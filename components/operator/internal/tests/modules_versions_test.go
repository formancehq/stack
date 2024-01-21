package tests_test

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	. "github.com/formancehq/operator/internal/tests/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Modules Versions", func() {
	Context("When creating a Ledger object", func() {
		var (
			stack                 *v1beta1.Stack
			ledger                *v1beta1.Ledger
			databaseConfiguration *v1beta1.DatabaseConfiguration
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
				Spec:       v1beta1.StackSpec{},
			}
			databaseConfiguration = &v1beta1.DatabaseConfiguration{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.DatabaseConfigurationSpec{
					ConfigurationProperties: v1beta1.ConfigurationProperties{
						Stacks: []string{stack.Name},
					},
					Service: "any",
				},
			}
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
			Expect(Create(databaseConfiguration)).To(Succeed())
			Expect(Create(ledger)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(ledger)).To(Succeed())
			Expect(Delete(databaseConfiguration)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
		Context("With version defined at Ledger level", func() {
			BeforeEach(func() {
				ledger.Spec.Version = "v1.2.3"
			})
			It("Should install v1.2.3 version", func() {
				versionsHistory := &v1beta1.VersionsHistory{}
				Eventually(func() error {
					return LoadResource("", core.InstalledVersionName(TestContext(), ledger, "v1.2.3"), versionsHistory)
				}).Should(Succeed())
				Expect(versionsHistory.Spec.Version).To(Equal("v1.2.3"))
				Expect(versionsHistory.Spec.Module).To(Equal("ledger"))
			})
		})
		Context("With version defined at Stack level", func() {
			BeforeEach(func() {
				stack.Spec.Version = "v2.3.4"
			})
			It("Should install v2.3.4 version", func() {
				versionsHistory := &v1beta1.VersionsHistory{}
				Eventually(func() error {
					return LoadResource("", core.InstalledVersionName(TestContext(), ledger, "v2.3.4"), versionsHistory)
				}).Should(Succeed())
				Expect(versionsHistory.Spec.Version).To(Equal("v2.3.4"))
				Expect(versionsHistory.Spec.Module).To(Equal("ledger"))
			})
		})
		Context("With version defined at Versions level", func() {
			var versions *v1beta1.Versions
			BeforeEach(func() {
				versions = &v1beta1.Versions{
					ObjectMeta: RandObjectMeta(),
					Spec: map[string]string{
						"ledger": "v3.4.5",
					},
				}
				Expect(Create(versions)).To(Succeed())
				stack.Spec.VersionsFromFile = versions.Name
			})
			AfterEach(func() {
				Expect(Delete(versions)).To(Succeed())
			})
			It("Should install v3.4.5 version", func() {
				versionsHistory := &v1beta1.VersionsHistory{}
				Eventually(func() error {
					return LoadResource("", core.InstalledVersionName(TestContext(), ledger, "v3.4.5"), versionsHistory)
				}).Should(Succeed())
				Expect(versionsHistory.Spec.Version).To(Equal("v3.4.5"))
				Expect(versionsHistory.Spec.Module).To(Equal("ledger"))
			})
		})
		Context("With no version defined", func() {
			It("should use 'latest'", func() {
				deployment := &appsv1.Deployment{}
				Eventually(func() error {
					return LoadResource(stack.Name, "ledger", deployment)
				}).Should(Succeed())
				Expect(deployment.Spec.Template.Spec.Containers[0].Image).To(
					Equal("ghcr.io/formancehq/ledger:latest"))
			})
			Context("then when updating the version at stack level", func() {
				JustBeforeEach(func() {
					Eventually(func(g Gomega) *v1beta1.Ledger {
						g.Expect(LoadResource("", ledger.Name, ledger)).To(Succeed())
						return ledger
					}).Should(BeReady())
					patch := client.MergeFrom(stack.DeepCopyObject().(client.Object))
					stack.Spec.Version = "v1.2.3"
					Expect(Patch(stack, patch)).To(Succeed())
				})
				It("Should install v1.2.3 version", func() {
					versionsHistory := &v1beta1.VersionsHistory{}
					Eventually(func() error {
						return LoadResource("", core.InstalledVersionName(TestContext(), ledger, "v1.2.3"), versionsHistory)
					}).Should(Succeed())
					Expect(versionsHistory.Spec.Version).To(Equal("v1.2.3"))
					Expect(versionsHistory.Spec.Module).To(Equal("ledger"))
				})
			})
			Context("then when updating the version at version file level", func() {
				var versions *v1beta1.Versions
				JustBeforeEach(func() {
					versions = &v1beta1.Versions{
						ObjectMeta: RandObjectMeta(),
						Spec:       map[string]string{},
					}
					Expect(Create(versions)).To(Succeed())
					patch := client.MergeFrom(stack.DeepCopyObject().(client.Object))
					stack.Spec.VersionsFromFile = versions.Name
					Expect(Patch(stack, patch)).To(Succeed())
				})
				JustAfterEach(func() {
					Expect(Delete(versions)).To(Succeed())
				})
				It("should use 'latest'", func() {
					deployment := &appsv1.Deployment{}
					Eventually(func(g Gomega) string {
						g.Expect(LoadResource(stack.Name, "ledger", deployment)).To(Succeed())
						return deployment.Spec.Template.Spec.Containers[0].Image
					}).Should(Equal("ghcr.io/formancehq/ledger:latest"))
				})
				Context("Then when updating the version file with another version", func() {
					JustBeforeEach(func() {
						patch := client.MergeFrom(versions.DeepCopyObject().(client.Object))
						versions.Spec = map[string]string{
							"ledger": "v2",
						}
						Expect(Patch(versions, patch)).To(Succeed())
					})
					It("Should install the new version", func() {
						versionsHistory := &v1beta1.VersionsHistory{}
						Eventually(func() error {
							return LoadResource("", core.InstalledVersionName(TestContext(), ledger, "v2"), versionsHistory)
						}).Should(Succeed())
						Expect(versionsHistory.Spec.Version).To(Equal("v2"))
						Expect(versionsHistory.Spec.Module).To(Equal("ledger"))
					})
				})
			})
		})
	})
})
