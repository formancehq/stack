package controller_test

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/controller"
	"github.com/formancehq/operator/v2/internal/controller/internal"
	. "github.com/formancehq/operator/v2/internal/controller/testing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("StackController", func() {
	Context("When creating a stack", func() {
		var (
			stack *v1beta1.Stack
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			Expect(Create(stack)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(stack)).To(Succeed())
		})
		It("Should create a new namespace", func() {
			Eventually(func() error {
				return Get(internal.GetResourceName(stack.Name), &corev1.Namespace{})
			}).Should(Succeed())
		})
		Context("With some secrets annotated with our annotations", func() {
			var (
				secret *corev1.Secret
			)
			BeforeEach(func() {
				secret = &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      uuid.NewString(),
						Labels: map[string]string{
							controller.ApplyToStacksLabel: stack.Name,
						},
					},
					StringData: map[string]string{
						"username": "formance",
						"password": "formance",
					},
				}
				Expect(Create(secret)).To(Succeed())
			})
			AfterEach(func() {
				Expect(client.IgnoreNotFound(Delete(secret))).To(Succeed())
			})
			It("Should replicate secret in the stack namespace", func() {
				Eventually(func() error {
					return LoadResource(stack.Name, secret.Name, &corev1.Secret{})
				}).Should(BeNil())
			})
		})
	})
})
