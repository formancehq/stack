package controller_test

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/controller/internal"
	. "github.com/formancehq/operator/v2/internal/controller/testing"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Labels: map[string]string{
						"formance.com/stack":   stack.Name,
						"formance.com/service": "any",
					},
				},
				Spec: v1beta1.DatabaseConfigurationSpec{},
			}
			ledger = &v1beta1.Ledger{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.LedgerSpec{
					Stack: stack.Name,
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
			Expect(deployment).To(BeOwnedBy(ledger))
		})
		It("Should create a new HTTPService object", func() {
			httpService := &v1beta1.HTTPAPI{}
			Eventually(func() error {
				return LoadResource("", internal.GetObjectName(stack.Name, "ledger"), httpService)
			}).Should(Succeed())
		})
		It("Should create a new Database object", func() {
			database := &v1beta1.Database{}
			Eventually(func() error {
				return LoadResource("", internal.GetObjectName(stack.Name, "ledger"), database)
			}).Should(Succeed())
		})
		Context("with monitoring enabled", func() {
			var (
				openTelemetryConfiguration *v1beta1.OpenTelemetryConfiguration
			)
			BeforeEach(func() {
				openTelemetryConfiguration = &v1beta1.OpenTelemetryConfiguration{
					ObjectMeta: RandObjectMeta(),
					Spec: v1beta1.OpenTelemetryConfigurationSpec{
						Stack: stack.Name,
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
			It("Should add correct env vars to the deployment", func() {
				Eventually(func(g Gomega) []corev1.EnvVar {
					deployment := &appsv1.Deployment{}
					g.Expect(Get(internal.GetNamespacedResourceName(stack.Name, "ledger"), deployment)).To(Succeed())
					return deployment.Spec.Template.Spec.Containers[0].Env
				}).Should(ContainElements(
					collectionutils.Map(
						internal.MonitoringEnvVars(openTelemetryConfiguration, "ledger"),
						collectionutils.ToAny[corev1.EnvVar],
					)...,
				))
			})
		})
		Context("with a Topic object existing on the ledger service", func() {
			deploymentShouldBeConfigured := func() {
				deployment := &appsv1.Deployment{}
				Eventually(func(g Gomega) bool {
					g.Expect(Get(internal.GetNamespacedResourceName(stack.Name, "ledger"), deployment)).To(Succeed())
					g.Expect(deployment.Spec.Template.Spec.Containers[0].Env).
						Should(ContainElement(internal.Env("BROKER", "nats")))
					return true
				}).Should(BeTrue())
			}
			var (
				brokerConfiguration *v1beta1.BrokerConfiguration
				topic               *v1beta1.Topic
			)
			JustBeforeEach(func() {
				brokerConfiguration = &v1beta1.BrokerConfiguration{
					ObjectMeta: metav1.ObjectMeta{
						Name: uuid.NewString(),
						Labels: map[string]string{
							"formance.com/stack":   stack.Name,
							"formance.com/service": "any",
						},
					},
					Spec: v1beta1.BrokerConfigSpec{
						Nats: &v1beta1.NatsConfig{
							URL: "nats://localhost:1234",
						},
					},
				}
				Expect(Create(brokerConfiguration)).To(Succeed())
				topic = &v1beta1.Topic{
					ObjectMeta: RandObjectMeta(),
					Spec: v1beta1.TopicSpec{
						Queries: []string{"orchestration"},
						Stack:   stack.Name,
						Service: "ledger",
					},
				}
				Expect(Create(topic)).To(Succeed())
			})
			It("Should start the deployment with env var defined for publishing in the event bus", deploymentShouldBeConfigured)
			Context("Then removing the Topic", func() {
				JustBeforeEach(func() {
					deploymentShouldBeConfigured()
					Expect(Delete(topic)).To(Succeed())
				})
				It("Should restart the deployment without broker env vars", func() {
					deployment := &appsv1.Deployment{}
					Eventually(func(g Gomega) bool {
						g.Expect(Get(internal.GetNamespacedResourceName(stack.Name, "ledger"), deployment)).To(Succeed())
						g.Expect(deployment.Spec.Template.Spec.Containers[0].Env).
							ShouldNot(ContainElements(internal.Env("BROKER", "nats")))
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
				Expect(reader).To(BeOwnedBy(ledger))

				readerService := &corev1.Service{}
				Eventually(func() error {
					return LoadResource(stack.Name, "ledger-read", readerService)
				}).Should(Succeed())
				Expect(readerService).To(BeOwnedBy(ledger))
				Expect(readerService).To(TargetDeployment(reader))

				writer := &appsv1.Deployment{}
				Eventually(func() error {
					return LoadResource(stack.Name, "ledger-write", writer)
				}).Should(Succeed())
				Expect(writer).To(BeOwnedBy(ledger))

				writerService := &corev1.Service{}
				Eventually(func() error {
					return LoadResource(stack.Name, "ledger-write", writerService)
				}).Should(Succeed())
				Expect(writerService).To(BeOwnedBy(ledger))
				Expect(writerService).To(TargetDeployment(writer))

				gateway := &appsv1.Deployment{}
				Eventually(func() error {
					return LoadResource(stack.Name, "ledger-gateway", gateway)
				}).Should(Succeed())
				Expect(gateway).To(BeOwnedBy(ledger))
			})
		})
	})
})