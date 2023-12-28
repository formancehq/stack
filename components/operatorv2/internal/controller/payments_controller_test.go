package controller_test

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/controller/internal"
	. "github.com/formancehq/operator/v2/internal/controller/testing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
						"formance.com/stack":   stack.Name,
						"formance.com/service": "any",
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
			Expect(deployment).To(BeOwnedBy(payments))

			service := &corev1.Service{}
			Eventually(func() error {
				return LoadResource(stack.Name, "payments-read", service)
			}).Should(Succeed())
			Expect(service).To(BeOwnedBy(payments))
		})
		It("Should create a connectors deployment with a service", func() {
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return LoadResource(stack.Name, "payments-connectors", deployment)
			}).Should(Succeed())
			Expect(deployment).To(BeOwnedBy(payments))

			service := &corev1.Service{}
			Eventually(func() error {
				return LoadResource(stack.Name, "payments-connectors", service)
			}).Should(Succeed())
			Expect(service).To(BeOwnedBy(payments))
		})
		It("Should create a gateway", func() {
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return LoadResource(stack.Name, "payments", deployment)
			}).Should(Succeed())
			Expect(deployment).To(BeOwnedBy(payments))
		})
		It("Should create a new HTTPAPI object", func() {
			httpService := &v1beta1.HTTPAPI{}
			Eventually(func() error {
				return LoadResource("", internal.GetObjectName(stack.Name, "payments"), httpService)
			}).Should(Succeed())
		})
	})
})
