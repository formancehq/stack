package tests_test

import (
	"fmt"

	. "github.com/formancehq/operator/internal/tests/internal"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("AuthClientController", func() {
	Context("When creating a AuthClient object", func() {
		var (
			stack      *v1beta1.Stack
			authClient *v1beta1.AuthClient
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			authClient = &v1beta1.AuthClient{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.AuthClientSpec{
					ID: uuid.NewString(),
					StackDependency: v1beta1.StackDependency{
						Stack: stack.Name,
					},
					Secret: uuid.NewString(),
				},
			}
		})
		JustBeforeEach(func() {
			Expect(Create(stack)).To(Succeed())
			Expect(Create(authClient)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(authClient)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
		It("Should create a secret", func() {
			secret := &corev1.Secret{}
			Eventually(func() error {
				return LoadResource(stack.Name, fmt.Sprintf("auth-client-%s", authClient.Name), secret)
			}).Should(Succeed())
		})
	})
})
