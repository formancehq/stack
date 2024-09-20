package tests_test

import (
	"github.com/formancehq/go-libs/collectionutils"
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/operator/internal/tests/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("PaymentsController", func() {
	Context("When creating a Payments object", func() {
		var (
			stack            *v1beta1.Stack
			payments         *v1beta1.Payments
			databaseSettings *v1beta1.Settings
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			payments = &v1beta1.Payments{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.PaymentsSpec{
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
				},
			}
			databaseSettings = settings.New(uuid.NewString(), "postgres.*.uri", "postgresql://localhost", stack.Name)
		})
		JustBeforeEach(func() {
			Expect(Create(stack)).To(Succeed())
			Expect(Create(databaseSettings)).To(Succeed())
			Expect(Create(payments)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(payments)).To(Succeed())
			Expect(Delete(databaseSettings)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
		It("Should create resources", func() {
			By("Should add an owner reference on the stack", func() {
				Eventually(func(g Gomega) bool {
					g.Expect(LoadResource("", payments.Name, payments)).To(Succeed())
					reference, err := core.HasOwnerReference(TestContext(), stack, payments)
					g.Expect(err).To(BeNil())
					return reference
				}).Should(BeTrue())
			})
			By("Should create a read deployment with a service", func() {
				deployment := &appsv1.Deployment{}
				Eventually(func() error {
					return LoadResource(stack.Name, "payments-read", deployment)
				}).Should(Succeed())
				Expect(deployment).To(BeControlledBy(payments))

				service := &corev1.Service{}
				Eventually(func() error {
					return LoadResource(stack.Name, "payments-read", service)
				}).Should(Succeed())
				Expect(service).To(BeControlledBy(payments))
			})
			By("Should create a connectors deployment with a service", func() {
				deployment := &appsv1.Deployment{}
				Eventually(func() error {
					return LoadResource(stack.Name, "payments-connectors", deployment)
				}).Should(Succeed())
				Expect(deployment).To(BeControlledBy(payments))

				service := &corev1.Service{}
				Eventually(func() error {
					return LoadResource(stack.Name, "payments-connectors", service)
				}).Should(Succeed())
				Expect(service).To(BeControlledBy(payments))
			})
			By("Should create a gateway", func() {
				deployment := &appsv1.Deployment{}
				Eventually(func() error {
					return LoadResource(stack.Name, "payments", deployment)
				}).Should(Succeed())
				Expect(deployment).To(BeControlledBy(payments))
			})
			By("Should create a new GatewayHTTPAPI object", func() {
				httpService := &v1beta1.GatewayHTTPAPI{}
				Eventually(func() error {
					return LoadResource("", core.GetObjectName(stack.Name, "payments"), httpService)
				}).Should(Succeed())
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
				l := &v1beta1.BenthosStreamList{}
				Eventually(func(g Gomega) int {
					g.Expect(List(l)).To(Succeed())
					return len(collectionutils.Filter(l.Items, func(stream v1beta1.BenthosStream) bool {
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
					l := &v1beta1.BenthosStreamList{}
					Eventually(func(g Gomega) int {
						g.Expect(List(l)).To(Succeed())
						return len(collectionutils.Filter(l.Items, func(stream v1beta1.BenthosStream) bool {
							return stream.Spec.Stack == stack.Name
						}))
					}).Should(BeZero())
				})
			})
		})
	})
})
