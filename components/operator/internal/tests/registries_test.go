package tests

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/operator/internal/tests/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
)

var _ = Describe("Registries", func() {
	var (
		stack            *v1beta1.Stack
		databaseSettings *v1beta1.Settings
		ledger           *v1beta1.Ledger
	)
	BeforeEach(func() {
		stack = &v1beta1.Stack{
			ObjectMeta: RandObjectMeta(),
		}
		databaseSettings = settings.New(uuid.NewString(), "postgres.*.uri", "postgresql://localhost", stack.Name)
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
		Expect(Create(databaseSettings)).To(Succeed())
		Expect(Create(ledger)).To(Succeed())

		DeferCleanup(func() {
			Expect(Delete(ledger)).To(Succeed())
			Expect(Delete(databaseSettings)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
	})
	Context("When overriding image in a settings", func() {
		var (
			settings     *v1beta1.Settings
			registry     = `"ghcr.io"`
			organization = "formancehq"
			ledgerpath   = fmt.Sprintf("%s/%s", organization, "ledger")
			imageRewrite = fmt.Sprintf("%s/%s", organization, "example")
		)

		BeforeEach(func() {
			settings = &v1beta1.Settings{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.SettingsSpec{
					Stacks: []string{"*"},
					Key:    "registries." + registry + ".images." + ledgerpath + ".rewrite",
					Value:  imageRewrite,
				},
			}
			Expect(Create(settings)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(settings)).To(Succeed())
		})
		It("Should have image re-written", func() {
			deployment := &appsv1.Deployment{}
			Eventually(func(g Gomega) error {
				g.Expect(LoadResource(stack.Name, "ledger", deployment)).To(Succeed())
				g.Expect(deployment.Spec.Template.Spec.Containers).To(HaveLen(1))
				g.Expect(deployment.Spec.Template.Spec.Containers[0].Image).To(ContainSubstring(imageRewrite))
				return nil
			}).Should(Succeed())
			Expect(deployment).To(BeControlledBy(ledger))
		})

	})
})
