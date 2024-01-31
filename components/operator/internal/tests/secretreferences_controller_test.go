package tests

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/tests/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("SecretReferenceController", func() {
	var stack *v1beta1.Stack
	BeforeEach(func() {
		stack = &v1beta1.Stack{
			ObjectMeta: RandObjectMeta(),
		}
		Expect(Create(stack)).To(Succeed())
	})
	Context("With a secret created on default namespace", func() {
		var secret *corev1.Secret
		BeforeEach(func() {
			secret = &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      uuid.NewString(),
					Namespace: "default",
					Labels: map[string]string{
						v1beta1.StackLabel: stack.Name,
					},
				},
			}
			Expect(Create(secret)).To(Succeed())
		})
		When("Creating a secret reference on a stack", func() {
			var secretReference *v1beta1.SecretReference
			BeforeEach(func() {
				secretReference = &v1beta1.SecretReference{
					ObjectMeta: RandObjectMeta(),
					Spec: v1beta1.SecretReferenceSpec{
						StackDependency: v1beta1.StackDependency{
							Stack: stack.Name,
						},
						SecretName: secret.Name,
					},
				}
				Expect(Create(secretReference)).To(Succeed())
			})
			shouldHaveReplicatedSecret := func() {
				replicatedSecret := &corev1.Secret{}
				Eventually(func() error {
					return LoadResource(stack.Name, secretReference.Spec.SecretName, replicatedSecret)
				}).Should(Succeed())
			}
			It("Should replicate the secret to the stack namespace", shouldHaveReplicatedSecret)
			Context("then when updating the referenced secret to a new one", func() {
				var newSecret *corev1.Secret
				BeforeEach(func() {
					shouldHaveReplicatedSecret()
					newSecret = &corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name:      uuid.NewString(),
							Namespace: "default",
							Labels: map[string]string{
								v1beta1.StackLabel: stack.Name,
							},
						},
					}
					Expect(Create(newSecret)).To(Succeed())

					patch := client.MergeFrom(secretReference.DeepCopy())
					secretReference.Spec.SecretName = newSecret.Name
					Expect(Patch(secretReference, patch)).To(Succeed())
				})
				It("Should replicate the new secret to the stack namespace", shouldHaveReplicatedSecret)
				It("Should remove old replicated secret", func() {
					Eventually(func() error {
						return LoadResource(stack.Name, secret.Name, &corev1.Secret{})
					}).Should(BeNotFound())
				})
			})
		})
	})
})
