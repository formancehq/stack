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
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

var _ = Describe("ResourceReferenceController", func() {
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
		When("Creating a resource reference on a stack", func() {
			var resourceReference *v1beta1.ResourceReference
			BeforeEach(func() {
				gvk, err := apiutil.GVKForObject(secret, GetScheme())
				Expect(err).To(BeNil())

				resourceReference = &v1beta1.ResourceReference{
					ObjectMeta: RandObjectMeta(),
					Spec: v1beta1.ResourceReferenceSpec{
						StackDependency: v1beta1.StackDependency{
							Stack: stack.Name,
						},
						Name: secret.Name,
						GroupVersionKind: &metav1.GroupVersionKind{
							Group:   gvk.Group,
							Version: gvk.Version,
							Kind:    gvk.Kind,
						},
					},
				}
				Expect(Create(resourceReference)).To(Succeed())
			})
			shouldHaveReplicatedSecret := func() {
				replicatedSecret := &corev1.Secret{}
				Eventually(func() error {
					return LoadResource(stack.Name, resourceReference.Spec.Name, replicatedSecret)
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

					patch := client.MergeFrom(resourceReference.DeepCopy())
					resourceReference.Spec.Name = newSecret.Name
					Expect(Patch(resourceReference, patch)).To(Succeed())
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
