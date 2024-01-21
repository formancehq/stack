package tests_test

import (
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/opentelemetryconfigurations"
	. "github.com/formancehq/operator/internal/tests/internal"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var _ = Describe("LedgerController", func() {
	Context("When creating a Ledger", func() {
		var (
			stack                 *v1beta1.Stack
			ledger                *v1beta1.Ledger
			databaseConfiguration *v1beta1.DatabaseConfiguration
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			databaseConfiguration = &v1beta1.DatabaseConfiguration{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.DatabaseConfigurationSpec{
					ConfigurationProperties: v1beta1.ConfigurationProperties{
						Stacks: []string{stack.Name},
					},
					Service: "any",
				},
			}
			ledger = &v1beta1.Ledger{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.LedgerSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
				},
			}
		})
		JustBeforeEach(func() {
			Expect(Create(stack)).To(Succeed())
			Expect(Create(databaseConfiguration)).To(Succeed())
			Expect(Create(ledger)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(ledger)).To(Succeed())
			Expect(Delete(databaseConfiguration)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
		It("Should create a deployment", func() {
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return LoadResource(stack.Name, "ledger", deployment)
			}).Should(Succeed())
			Expect(deployment).To(BeControlledBy(ledger))
		})
		It("Should create a new HTTPAPI object", func() {
			httpService := &v1beta1.HTTPAPI{}
			Eventually(func() error {
				return LoadResource("", core.GetObjectName(stack.Name, "ledger"), httpService)
			}).Should(Succeed())
		})
		It("Should create a new Database object", func() {
			database := &v1beta1.Database{}
			Eventually(func() error {
				return LoadResource("", core.GetObjectName(stack.Name, "ledger"), database)
			}).Should(Succeed())
		})
		Context("with monitoring enabled", func() {
			var openTelemetryConfiguration *v1beta1.OpenTelemetryConfiguration
			BeforeEach(func() {
				openTelemetryConfiguration = &v1beta1.OpenTelemetryConfiguration{
					ObjectMeta: RandObjectMeta(),
					Spec: v1beta1.OpenTelemetryConfigurationSpec{
						ConfigurationProperties: v1beta1.ConfigurationProperties{
							Stacks: []string{stack.Name},
						},
						Traces: &v1beta1.TracesSpec{
							Otlp: &v1beta1.OtlpSpec{
								Endpoint: "otel-collector",
								Port:     4317,
								Insecure: true,
							},
						},
					},
				}
				Expect(Create(openTelemetryConfiguration)).To(Succeed())
			})
			AfterEach(func() {
				Expect(Delete(openTelemetryConfiguration)).To(Succeed())
			})
			It("Should add correct env vars to the deployment", func() {
				Eventually(func(g Gomega) []corev1.EnvVar {
					deployment := &appsv1.Deployment{}
					g.Expect(Get(core.GetNamespacedResourceName(stack.Name, "ledger"), deployment)).To(Succeed())
					return deployment.Spec.Template.Spec.Containers[0].Env
				}).Should(ContainElements(
					collectionutils.Map(
						opentelemetryconfigurations.GetEnvVars(openTelemetryConfiguration, "ledger"),
						collectionutils.ToAny[corev1.EnvVar],
					)...,
				))
			})
		})
		Context("with a BrokerTopic object existing on the ledger service", func() {
			deploymentShouldBeConfigured := func() {
				deployment := &appsv1.Deployment{}
				Eventually(func(g Gomega) bool {
					g.Expect(Get(core.GetNamespacedResourceName(stack.Name, "ledger"), deployment)).To(Succeed())
					g.Expect(deployment.Spec.Template.Spec.Containers[0].Env).
						Should(ContainElement(core.Env("BROKER", "nats")))
					return true
				}).Should(BeTrue())
			}
			var (
				brokerConfiguration *v1beta1.BrokerConfiguration
				brokerTopic         *v1beta1.BrokerTopic
			)
			JustBeforeEach(func() {
				brokerConfiguration = &v1beta1.BrokerConfiguration{
					ObjectMeta: RandObjectMeta(),
					Spec: v1beta1.BrokerConfigurationSpec{
						ConfigurationProperties: v1beta1.ConfigurationProperties{
							Stacks: []string{stack.Name},
						},
						Nats: &v1beta1.BrokerNatsConfig{
							URL: "nats://localhost:1234",
						},
					},
				}
				Expect(Create(brokerConfiguration)).To(Succeed())
				brokerTopic = &v1beta1.BrokerTopic{
					ObjectMeta: RandObjectMeta(),
					Spec: v1beta1.BrokerTopicSpec{
						StackDependency: v1beta1.StackDependency{
							Stack: stack.Name,
						},
						Service: "ledger",
					},
				}
				Expect(controllerutil.SetOwnerReference(stack, brokerTopic, GetScheme()))
				Expect(Create(brokerTopic)).To(Succeed())
			})
			AfterEach(func() {
				Expect(Delete(brokerConfiguration)).To(Succeed())
				Expect(client.IgnoreNotFound(Delete(brokerTopic))).To(Succeed())
			})
			It("Should start the deployment with env var defined for publishing in the event bus", deploymentShouldBeConfigured)
			Context("Then removing the BrokerTopic", func() {
				JustBeforeEach(func() {
					deploymentShouldBeConfigured()
					Expect(Delete(brokerTopic)).To(Succeed())
				})
				It("Should restart the deployment without broker env vars", func() {
					deployment := &appsv1.Deployment{}
					Eventually(func(g Gomega) bool {
						g.Expect(Get(core.GetNamespacedResourceName(stack.Name, "ledger"), deployment)).To(Succeed())
						g.Expect(deployment.Spec.Template.Spec.Containers[0].Env).
							ShouldNot(ContainElements(core.Env("BROKER", "nats")))
						return true
					}).Should(BeTrue())
				})
			})
		})
		Context("with multi ready deployment strategy", func() {
			JustBeforeEach(func() {
				Eventually(func() error {
					return LoadResource(stack.Name, "ledger", &appsv1.Deployment{})
				}).Should(Succeed())
				patch := client.MergeFrom(ledger.DeepCopy())
				ledger.Spec.DeploymentStrategy = v1beta1.DeploymentStrategyMonoWriterMultipleReader
				Expect(Patch(ledger, patch)).To(Succeed())
			})
			It("Should remove the original deployment", func() {
				Eventually(func() error {
					return LoadResource(stack.Name, "ledger", &appsv1.Deployment{})
				}).Should(BeNotFound())
			})
			It("Should create two deployments, two services and a gateway", func() {
				reader := &appsv1.Deployment{}
				Eventually(func() error {
					return LoadResource(stack.Name, "ledger-read", reader)
				}).Should(Succeed())
				Expect(reader).To(BeControlledBy(ledger))

				readerService := &corev1.Service{}
				Eventually(func() error {
					return LoadResource(stack.Name, "ledger-read", readerService)
				}).Should(Succeed())
				Expect(readerService).To(BeControlledBy(ledger))
				Expect(readerService).To(TargetDeployment(reader))

				writer := &appsv1.Deployment{}
				Eventually(func() error {
					return LoadResource(stack.Name, "ledger-write", writer)
				}).Should(Succeed())
				Expect(writer).To(BeControlledBy(ledger))

				writerService := &corev1.Service{}
				Eventually(func() error {
					return LoadResource(stack.Name, "ledger-write", writerService)
				}).Should(Succeed())
				Expect(writerService).To(BeControlledBy(ledger))
				Expect(writerService).To(TargetDeployment(writer))

				gateway := &appsv1.Deployment{}
				Eventually(func() error {
					return LoadResource(stack.Name, "ledger-gateway", gateway)
				}).Should(Succeed())
				Expect(gateway).To(BeControlledBy(ledger))
			})
		})
		Context("With Search enabled", func() {
			var search *v1beta1.Search
			BeforeEach(func() {
				search = &v1beta1.Search{
					ObjectMeta: RandObjectMeta(),
					Spec: v1beta1.SearchSpec{
						StackDependency: v1beta1.StackDependency{
							Stack: stack.Name,
						},
					},
				}
			})
			JustBeforeEach(func() {
				Expect(Create(search)).To(Succeed())
			})
			JustAfterEach(func() {
				Expect(client.IgnoreNotFound(Delete(search))).To(Succeed())
			})
			checkStreamsExists := func() {
				l := &v1beta1.StreamList{}
				Eventually(func(g Gomega) int {
					g.Expect(List(l)).To(Succeed())
					return len(collectionutils.Filter(l.Items, func(stream v1beta1.Stream) bool {
						return stream.Spec.Stack == stack.Name
					}))
				}).Should(BeNumerically(">", 0))
			}
			It("Should create streams", checkStreamsExists)
			Context("Then when removing search", func() {
				JustBeforeEach(func() {
					checkStreamsExists()
					Expect(Delete(search)).To(Succeed())
				})
				It("Should remove streams", func() {
					l := &v1beta1.StreamList{}
					Eventually(func(g Gomega) int {
						g.Expect(List(l)).To(Succeed())
						return len(collectionutils.Filter(l.Items, func(stream v1beta1.Stream) bool {
							return stream.Spec.Stack == stack.Name
						}))
					}).Should(BeZero())
				})
			})
		})
	})
})
