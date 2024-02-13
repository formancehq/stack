package tests_test

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/operator/internal/tests/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("StackController", FlakeAttempts(3), func() {
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
			Expect(client.IgnoreNotFound(Delete(stack))).To(Succeed())
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
		It("Should be ready", func() {
			Eventually(func() bool {
				Expect(LoadResource("", stack.Name, stack)).To(Succeed())
				return stack.Status.Ready
			}).Should(BeTrue())
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
		Context("with a module", func() {
			var (
				ledger *v1beta1.Ledger
			)
			JustBeforeEach(func() {
				ledger = &v1beta1.Ledger{
					ObjectMeta: RandObjectMeta(),
					Spec: v1beta1.LedgerSpec{
						StackDependency: v1beta1.StackDependency{
							Stack: stack.Name,
						},
					},
				}

				Expect(Create(ledger)).To(Succeed())
				Eventually(func(g Gomega) *v1beta1.Ledger {
					g.Expect(LoadResource("", ledger.Name, ledger)).To(Succeed())
					return ledger
				}).Should(BeOwnedBy(stack))
			})
			JustAfterEach(func() {
				Expect(client.IgnoreNotFound(Delete(ledger))).To(Succeed())
			})
			It("the stack should not be ready anymore", func() {
				Eventually(func() bool {
					Expect(LoadResource("", stack.Name, stack)).To(Succeed())
					return stack.Status.Ready
				}).Should(BeFalse())
			})
			It("should not be ready", func() {
				Eventually(func() bool {
					Expect(LoadResource("", ledger.Name, ledger)).To(Succeed())
					return ledger.Status.Ready
				}).Should(BeFalse())
			})
			When("each module are ready", func() {
				var databaseSettings *v1beta1.Settings
				BeforeEach(func() {
					databaseSettings = settings.New(uuid.NewString(), "postgres.*.uri", "postgresql://localhost", stack.Name)
					Expect(Create(databaseSettings)).Should(Succeed())
				})
				JustAfterEach(func() {
					Expect(Delete(databaseSettings)).To(Succeed())
				})
				JustBeforeEach(func() {
					Eventually(func() bool {
						Expect(LoadResource("", ledger.Name, ledger)).To(Succeed())
						return ledger.Status.Ready
					}).Should(BeTrue())
				})
				It("the stack should be ready", func() {
					Eventually(func() bool {
						err := LoadResource("", stack.Name, stack)
						Expect(err).ToNot(HaveOccurred())
						return stack.Status.Ready
					}).Should(BeTrue())
				})
			})
			When("deleting the stack", func() {
				JustBeforeEach(func() {
					Expect(Delete(stack)).To(Succeed())
				})
				It("Should also delete the module", func() {
					Eventually(func(g Gomega) error {
						return LoadResource("", ledger.Name, ledger)
					}).Should(BeNotFound())
					Eventually(func(g Gomega) error {
						return LoadResource("", stack.Name, stack)
					}).Should(BeNotFound())
				})
			})
		})
	})
})
