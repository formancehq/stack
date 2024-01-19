package tests_test

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	. "github.com/formancehq/operator/internal/tests/internal"
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
		})
		JustBeforeEach(func() {
			Expect(Create(stack)).To(Succeed())
		})
		JustAfterEach(func() {
			Expect(Delete(stack)).To(Succeed())
		})
		It("Should create a new namespace", func() {
			Eventually(func() error {
				return Get(core.GetResourceName(stack.Name), &corev1.Namespace{})
			}).Should(Succeed())
		})
		It("Should resolve to 'latest' version", func() {
			version, err := core.GetModuleVersion(TestContext(), stack, &v1beta1.Ledger{})
			Expect(err).To(Succeed())
			Expect(version).To(Equal("latest"))
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
							core.StackLabel: stack.Name,
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
		Context("with version specified", func() {
			BeforeEach(func() {
				stack.Spec.Version = "1234"
			})
			It("should resolve a module to the specified version", func() {
				version, err := core.GetModuleVersion(TestContext(), stack, &v1beta1.Ledger{})
				Expect(err).To(Succeed())
				Expect(version).To(Equal("1234"))
			})
		})
		Context("with version file specified", func() {
			var versions *v1beta1.Versions
			BeforeEach(func() {
				versions = &v1beta1.Versions{
					ObjectMeta: RandObjectMeta(),
					Spec:       map[string]string{},
				}
				stack.Spec.VersionsFromFile = versions.Name
			})
			JustBeforeEach(func() {
				Expect(Create(versions)).To(Succeed())
			})
			JustAfterEach(func() {
				Expect(Delete(versions)).To(Succeed())
			})
			Context("with no specific version", func() {
				It("should resolve a module to 'latest'", func() {
					version, err := core.GetModuleVersion(TestContext(), stack, &v1beta1.Ledger{})
					Expect(err).To(Succeed())
					Expect(version).To(Equal("latest"))
				})
			})
			Context("with specific version for a module", func() {
				BeforeEach(func() {
					versions.Spec["ledger"] = "5678"
				})
				It("should resolve to the correct version", func() {
					version, err := core.GetModuleVersion(TestContext(), stack, &v1beta1.Ledger{})
					Expect(err).To(Succeed())
					Expect(version).To(Equal("5678"))
				})
			})
		})
	})
})
