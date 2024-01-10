package formance_com_test

import (
	v1beta1 "github.com/formancehq/operator/v2/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controller/testing"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("PaymentsController", func() {
	Context("When creating a Payments object", func() {
		var (
			stack                 *v1beta1.Stack
			payments              *v1beta1.Payments
			databaseConfiguration *v1beta1.DatabaseConfiguration
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
			databaseConfiguration = &v1beta1.DatabaseConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name: uuid.NewString(),
					Labels: map[string]string{
						core.StackLabel:   stack.Name,
						core.ServiceLabel: "any",
					},
				},
				Spec: v1beta1.DatabaseConfigurationSpec{},
			}
		})
		JustBeforeEach(func() {
			Expect(Create(stack)).To(Succeed())
			Expect(Create(databaseConfiguration)).To(Succeed())
			Expect(Create(payments)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(payments)).To(Succeed())
			Expect(Delete(databaseConfiguration)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
		It("Should create a read deployment with a service", func() {
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
		It("Should create a connectors deployment with a service", func() {
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
		It("Should create a gateway", func() {
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return LoadResource(stack.Name, "payments", deployment)
			}).Should(Succeed())
			Expect(deployment).To(BeControlledBy(payments))
		})
		It("Should create a new HTTPAPI object", func() {
			httpService := &v1beta1.HTTPAPI{}
			Eventually(func() error {
				return LoadResource("", core.GetObjectName(stack.Name, "payments"), httpService)
			}).Should(Succeed())
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
