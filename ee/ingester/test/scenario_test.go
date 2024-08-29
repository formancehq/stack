//go:build it

package test_suite

import (
	"context"
	"github.com/formancehq/ingester/ingesterclient"
	. "github.com/formancehq/stack/ee/ingester/pkg/testserver"
	"github.com/formancehq/stack/libs/go-libs/logging"
	. "github.com/formancehq/stack/libs/go-libs/testing/api"
	. "github.com/formancehq/stack/libs/go-libs/testing/platform/natstesting"
	. "github.com/formancehq/stack/libs/go-libs/testing/platform/pgtesting"
	. "github.com/formancehq/stack/libs/go-libs/testing/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
)

func runTest(stackName *Deferred[string], natsServer *Deferred[*NatsServer], connectorName string, connectorAdapter *Deferred[Connector]) {
	var (
		db        = UsePostgresDatabase(pgServer)
		ctx       = logging.TestingContext()
		module    = "module1"
		pipeline  ingesterclient.Pipeline
		connector ingesterclient.Connector
		service   *FakeModule
	)
	BeforeEach(func() {
		By("creating a new fake module to allow the pull of logs")
		service = NewFakeModule(GinkgoT(), stackName.GetValue(), module, PublisherFn(func(_ require.TestingT, stack string, module string, data []byte) {
			natsServer.GetValue().Publish(GinkgoT(), stack, module, data)
		}))
	})

	testServer := UseNewTestServer(func() Configuration {
		return Configuration{
			Stack:                 stackName.GetValue(),
			PostgresConfiguration: db.GetValue().ConnectionOptions(),
			NatsURL:               natsServer.GetValue().ClientURL(),
			Output:                GinkgoWriter,
			ModuleURLTemplate:     service.URL(),
			Debug:                 debug,
		}
	})
	When("reading a not existing connector", func() {
		It("should return a NOT_FOUND error", func() {
			_, err := GetConnector(context.Background(), testServer.GetValue(), "xxx")
			Expect(err).To(HaveErrorCode("NOT_FOUND"))
		})
	})
	When("reading a not existing pipeline", func() {
		It("should return a NOT_FOUND error", func() {
			_, err := GetPipeline(context.Background(), testServer.GetValue(), "xxx")
			Expect(err).To(HaveErrorCode("NOT_FOUND"))
		})
	})
	When("deleting a not existing connector", func() {
		It("should return a NOT_FOUND error", func() {
			err := DeleteConnector(context.Background(), testServer.GetValue(), "xxx")
			Expect(err).To(HaveErrorCode("NOT_FOUND"))
		})
	})
	When("deleting a not existing pipeline", func() {
		It("should return a NOT_FOUND error", func() {
			err := DeletePipeline(context.Background(), testServer.GetValue(), "xxx")
			Expect(err).To(HaveErrorCode("NOT_FOUND"))
		})
	})
	When("creating a connector with invalid configuration", func() {
		It("should fail with VALIDATION error", func() {
			_, err := CreateConnector(context.Background(), testServer.GetValue(), ingesterclient.ConnectorConfiguration{
				Driver: "noop",
				Config: map[string]any{
					"batching": map[string]any{
						"maxItems": -10,
					},
				},
			})
			Expect(err).To(HaveErrorCode("VALIDATION"))
		})
	})
	When("Creating a connector with valid configuration", func() {
		BeforeEach(func() {
			config := connectorAdapter.GetValue().Config()
			config["batching"] = map[string]any{
				"flushInterval": "50ms",
			}
			connector = MustCreateConnector(testServer, ingesterclient.ConnectorConfiguration{
				Driver: connectorName,
				Config: config,
			})
			Expect(connector.CreatedAt).NotTo(BeZero())
			Expect(connector.Driver).To(Equal(connectorName))
			Expect(connector.Id).NotTo(BeZero())
			Expect(connector.Config).NotTo(BeEmpty())
		})
		When("listing connectors", func() {
			It("should be ok", func() {
				connectors := ListConnectors(context.Background(), testServer.GetValue()).Data
				Expect(connectors).To(HaveLen(1))
				Expect(connectors[0]).To(Equal(connector))
			})
		})
		When("reading connector", func() {
			It("should be ok", func() {
				connectorFromAPI := MustGetConnector(testServer, connector.Id)
				Expect(connectorFromAPI).To(Equal(connector))
			})
		})
		When("deleting it", func() {
			It("should be ok", func() {
				MustDeleteConnector(testServer, connector.Id)
			})
		})

		When("creating a new pipeline", func() {
			BeforeEach(func() {
				pipelineConfiguration := *ingesterclient.NewPipelineConfiguration("module1", connector.Id)
				pipeline = MustCreatePipeline(testServer, pipelineConfiguration)
				Expect(pipeline.CreatedAt).NotTo(BeZero())
				Expect(pipeline.ConnectorID).To(Equal(connector.Id))
				Expect(pipeline.Id).NotTo(BeZero())
				Expect(pipeline.Module).To(Equal(module))
			})
			When("reading pipeline", func() {
				It("should be ok", func() {
					pipelineFromAPI := MustGetPipeline(testServer, pipeline.Id)
					pipelineFromAPI.State = pipeline.State
					Expect(pipelineFromAPI).To(Equal(pipeline))
				})
			})
			When("listing pipelines", func() {
				It("should be ok", func() {
					pipelines := ListPipelines(context.Background(), testServer.GetValue()).Data
					Expect(pipelines).To(HaveLen(1))
					pipelineFromAPI := pipelines[0]
					pipelineFromAPI.State = pipeline.State
					Expect(pipelineFromAPI).To(Equal(pipeline))
				})
			})

			When("sending a bunch of messages", func() {
				const numberOfMessages = 50
				BeforeEach(func() {
					By("sending a bunch of messages")
					sendBatch(service, 0, numberOfMessages)
				})
				It("it should eventually be inserted", func() {
					shouldHaveInsert(connectorAdapter, numberOfMessages)
				})
				When("deleting its associated connector", func() {
					It("should fail", func() {
						err := DeleteConnector(ctx, testServer.GetValue(), connector.Id)
						Expect(err).To(HaveErrorCode("VALIDATION"))
					})
				})
				When("deleting it", func() {
					It("should be ok", func() {
						Expect(DeletePipeline(ctx, testServer.GetValue(), pipeline.Id)).To(Succeed())
					})
				})
				When("creating another pipeline and listing pipelines", func() {
					var newPipelineConfiguration ingesterclient.PipelineConfiguration
					BeforeEach(func() {
						newPipelineConfiguration = *ingesterclient.NewPipelineConfiguration("module2", connector.Id)
						MustCreatePipeline(testServer, newPipelineConfiguration)
					})
					It("should be ok", func() {
						Expect(ListPipelines(ctx, testServer.GetValue()).Data).To(HaveLen(2))
					})
				})
				When("resetting it", func() {
					BeforeEach(func() {
						shouldHaveInsert(connectorAdapter, numberOfMessages)
						Expect(connectorAdapter.GetValue().Clear(ctx)).To(BeNil())
						shouldHaveInsert(connectorAdapter, 0)

						MustResetPipeline(testServer, pipeline.Id)
					})
					It("it should have poll old logs", func() {
						shouldHaveInsert(connectorAdapter, numberOfMessages)
					})
				})
				When("pausing the pipeline", func() {
					BeforeEach(func() {
						// We still need to wait for the messages ingestion
						shouldHaveInsert(connectorAdapter, numberOfMessages)
						MustPausePipeline(testServer, pipeline.Id)
					})
					When("sending another bunch of messages", func() {
						BeforeEach(func() {
							sendBatch(service, numberOfMessages, numberOfMessages)
						})
						It("should consistently not sending messages to the connector", shouldHaveInsertConsistently(connectorAdapter, numberOfMessages))
						When("resuming the connector", func() {
							BeforeEach(func() {
								MustResumePipeline(testServer, pipeline.Id)
							})
							It("should ingest pending messages", func() {
								shouldHaveInsert(connectorAdapter, 2*numberOfMessages)
							})
							When("resuming again the connector", func() {
								It("should fail when resuming again", func() {
									Expect(ResumePipeline(ctx, testServer.GetValue(), pipeline.Id)).To(HaveErrorCode("VALIDATION"))
								})
							})
						})
					})
					It("should fail when pausing again", func() {
						Expect(PausePipeline(ctx, testServer.GetValue(), pipeline.Id)).To(HaveErrorCode("VALIDATION"))
					})
				})
				When("stopping the pipeline", func() {
					BeforeEach(func() {
						// We still need to wait for the messages ingestion
						shouldHaveInsert(connectorAdapter, numberOfMessages)
						MustStopPipeline(testServer, pipeline.Id)
					})
					It("should consistently not sending messages to the connector", shouldHaveInsertConsistently(connectorAdapter, numberOfMessages))
					When("stopping again the pipeline", func() {
						It("should fail", func() {
							Expect(StopPipeline(ctx, testServer.GetValue(), pipeline.Id)).To(HaveErrorCode("VALIDATION"))
						})
					})
					When("starting the pipeline", func() {
						BeforeEach(func() {
							MustStartPipeline(testServer, pipeline.Id)
						})
						It("should be ok", func() {})
						When("sending a message", func() {
							BeforeEach(func() {
								sendBatch(service, numberOfMessages, 1)
							})
							It("should be ingested", func() {
								shouldHaveInsert(connectorAdapter, numberOfMessages+1)
							})
						})
						When("deleting it", func() {
							It("should be ok", func() {
								Expect(DeletePipeline(ctx, testServer.GetValue(), pipeline.Id)).To(Succeed())
							})
						})
					})
				})
				It("should be ok to stop the server", func() {
					testServer.GetValue().Stop()
				})
				When("restarting the server", func() {
					BeforeEach(func() {
						testServer.GetValue().Restart()
					})
					When("sending a message", func() {
						BeforeEach(func() {
							sendBatch(service, numberOfMessages, 1)
						})
						It("should eventually be ingested", func() {
							shouldHaveInsert(connectorAdapter, numberOfMessages+1)
						})
					})
				})
			})
		})
	})
}
