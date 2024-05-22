package tests_test

import (
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/operator/internal/tests/internal"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var _ = Describe("LedgerController", func() {
	Context("When creating a Ledger", func() {
		var (
			stack            *v1beta1.Stack
			ledger           *v1beta1.Ledger
			databaseSettings *v1beta1.Settings
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			databaseSettings = settings.New(uuid.NewString(), "postgres.*.uri", "postgresql://localhost", stack.Name)
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
			Expect(Create(databaseSettings)).To(Succeed())
			Expect(Create(ledger)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(ledger)).To(Succeed())
			Expect(Delete(databaseSettings)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
		It("Should add an owner reference on the stack", func() {
			Eventually(func(g Gomega) bool {
				g.Expect(LoadResource("", ledger.Name, ledger)).To(Succeed())
				reference, err := core.HasOwnerReference(TestContext(), stack, ledger)
				g.Expect(err).To(BeNil())
				return reference
			}).Should(BeTrue())
		})
		It("Should create a deployment", func() {
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return LoadResource(stack.Name, "ledger", deployment)
			}).Should(Succeed())
			Expect(deployment).To(BeControlledBy(ledger))
		})
		It("Should create a new GatewayHTTPAPI object", func() {
			httpService := &v1beta1.GatewayHTTPAPI{}
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
			var (
				otelTracesDSNSetting *v1beta1.Settings
			)
			BeforeEach(func() {
				otelTracesDSNSetting = settings.New(uuid.NewString(), "opentelemetry.traces.dsn", "grpc://collector", stack.Name)
				Expect(Create(otelTracesDSNSetting)).To(Succeed())
			})
			AfterEach(func() {
				Expect(Delete(otelTracesDSNSetting)).To(Succeed())
			})
			It("Should add correct env vars to the deployment", func() {
				Eventually(func(g Gomega) []corev1.EnvVar {
					deployment := &appsv1.Deployment{}
					g.Expect(Get(core.GetNamespacedResourceName(stack.Name, "ledger"), deployment)).To(Succeed())
					return deployment.Spec.Template.Spec.Containers[0].Env
				}).Should(ContainElements(corev1.EnvVar{
					Name:  "OTEL_SERVICE_NAME",
					Value: "ledger",
				}))
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
				brokerDSNSettings *v1beta1.Settings
				brokerTopic       *v1beta1.BrokerTopic
			)
			JustBeforeEach(func() {
				brokerDSNSettings = settings.New(uuid.NewString(), "broker.dsn", "nats://localhost:1234", stack.Name)
				Expect(Create(brokerDSNSettings)).To(BeNil())
				brokerTopic = &v1beta1.BrokerTopic{
					ObjectMeta: RandObjectMeta(),
					Spec: v1beta1.BrokerTopicSpec{
						StackDependency: v1beta1.StackDependency{
							Stack: stack.Name,
						},
						Service: "ledger",
					},
				}
				// notes(gfyrag): add "fake" owner reference to prevent auto deletion of topics
				Expect(controllerutil.SetOwnerReference(databaseSettings, brokerTopic, GetScheme()))
				Expect(controllerutil.SetOwnerReference(stack, brokerTopic, GetScheme()))
				Expect(Create(brokerTopic)).To(Succeed())
			})
			AfterEach(func() {
				Expect(Delete(brokerDSNSettings)).To(Succeed())
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
			checkResourcesExists := func() {
				l := &v1beta1.BenthosStreamList{}
				Eventually(func(g Gomega) int {
					g.Expect(List(l)).To(Succeed())
					return len(collectionutils.Filter(l.Items, func(stream v1beta1.BenthosStream) bool {
						return stream.Spec.Stack == stack.Name
					}))
				}).Should(BeNumerically(">", 0))

				cronJob := &v1.CronJob{}
				Eventually(func() error {
					return Get(types.NamespacedName{
						Namespace: stack.Name,
						Name:      "reindex-ledger",
					}, cronJob)
				}).Should(BeNil())
			}
			It("Should create appropriate resources", checkResourcesExists)
			Context("Then when removing search", func() {
				JustBeforeEach(func() {
					checkResourcesExists()
					Expect(Delete(search)).To(Succeed())
				})
				It("Should remove resources", func() {
					l := &v1beta1.BenthosStreamList{}
					Eventually(func(g Gomega) int {
						g.Expect(List(l)).To(Succeed())
						return len(collectionutils.Filter(l.Items, func(stream v1beta1.BenthosStream) bool {
							return stream.Spec.Stack == stack.Name
						}))
					}).Should(BeZero())
					Eventually(func() error {
						return Get(types.NamespacedName{
							Namespace: stack.Name,
							Name:      "reindex-ledger",
						}, &v1.CronJob{})
					}).Should(BeNotFound())
				})
			})
		})
	})
})
