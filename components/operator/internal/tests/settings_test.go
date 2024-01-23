package tests

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/operator/internal/tests/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SettingsUtils", func() {
	Context("When creating some settings for a database", func() {
		var (
			stack                   *v1beta1.Stack
			wildcardPostgresSetting *v1beta1.Settings
			ledgerPostgresSetting   *v1beta1.Settings
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
			}
			Expect(Create(stack)).To(Succeed())
			wildcardPostgresSetting = settings.New(uuid.NewString(), "postgres.*.uri", "postgresql://postgres1:5432", stack.Name)
			ledgerPostgresSetting = settings.New(uuid.NewString(), "postgres.ledger.uri", "postgresql://postgres2:5432", stack.Name)
			Expect(Create(wildcardPostgresSetting, ledgerPostgresSetting)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(wildcardPostgresSetting, ledgerPostgresSetting)).To(Succeed())
		})
		FIt("Should resolve properly", func() {
			value, err := settings.GetByPriority(TestContext(), stack.Name, "postgres", "ledger", "uri")
			Expect(err).To(BeNil())
			Expect(value).NotTo(BeNil())
			Expect(*value).To(Equal(ledgerPostgresSetting.Spec.Value))

			value, err = settings.GetByPriority(TestContext(), stack.Name, "postgres", "payments", "uri")
			Expect(err).To(BeNil())
			Expect(value).NotTo(BeNil())
			Expect(*value).To(Equal(wildcardPostgresSetting.Spec.Value))
		})
	})
})
