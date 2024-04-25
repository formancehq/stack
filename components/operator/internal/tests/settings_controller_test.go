package tests

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	helperSettings "github.com/formancehq/operator/internal/resources/settings"
	. "github.com/formancehq/operator/internal/tests/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SettingsController", func() {
	When("Creating a settings on specific stack", func() {
		var (
			setting *v1beta1.Settings
			stack   *v1beta1.Stack
		)
		BeforeEach(func() {
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
			}
			setting = &v1beta1.Settings{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.SettingsSpec{
					Stacks: []string{stack.Name},
					Key:    uuid.NewString() + `.` + uuid.NewString(),
					Value:  uuid.NewString(),
				},
			}

		})
		JustBeforeEach(func() {
			Expect(Create(stack)).To(Succeed())
			Expect(Create(setting)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(stack)).To(Succeed())
			Expect(Delete(setting)).To(Succeed())
		})
		It("Should create resources", func() {
			By("Should have created the setting", func() {
				Expect(LoadResource("", setting.Name, &v1beta1.Settings{})).To(Succeed())
			})
			By("Sould be able to retrieve the settings", func() {
				str, err := helperSettings.Get(TestContext(), stack.Name, helperSettings.SplitKeywordWithDot(setting.Spec.Key)...)
				Expect(err).To(BeNil())
				Expect(str).ToNot(BeNil())
				Expect(*str).To(Equal(setting.Spec.Value))
			})
		})
	})
	When("Creating a settings", func() {
		var (
			setting *v1beta1.Settings
			stack   *v1beta1.Stack
		)
		BeforeEach(func() {
			setting = &v1beta1.Settings{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.SettingsSpec{
					Stacks: []string{"*"},
					Key:    uuid.NewString() + `.` + uuid.NewString(),
					Value:  uuid.NewString(),
				},
			}

			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
			}

		})
		JustBeforeEach(func() {
			Expect(Create(setting)).To(Succeed())
			Expect(Create(stack)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(setting)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
		It("Should create resources", func() {
			By("Should have created the setting", func() {
				Expect(LoadResource("", setting.Name, &v1beta1.Settings{})).To(Succeed())
			})
			By("Sould be able to retrieve the settings", func() {
				str, err := helperSettings.Get(TestContext(), stack.Name, helperSettings.SplitKeywordWithDot(setting.Spec.Key)...)
				Expect(err).To(BeNil())
				Expect(str).ToNot(BeNil())
				Expect(*str).To(Equal(setting.Spec.Value))
			})
		})
	})
	When("Creating a settings with escaped key", func() {
		var (
			setting *v1beta1.Settings
			stack   *v1beta1.Stack
		)
		BeforeEach(func() {
			setting = &v1beta1.Settings{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta1.SettingsSpec{
					Stacks: []string{"*"},
					Key:    uuid.NewString() + `."example.net"`,
					Value:  uuid.NewString(),
				},
			}
			stack = &v1beta1.Stack{
				ObjectMeta: RandObjectMeta(),
			}
		})
		JustBeforeEach(func() {
			Expect(Create(setting)).To(Succeed())
			Expect(Create(stack)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(setting)).To(Succeed())
			Expect(Delete(stack)).To(Succeed())
		})
		It("Should create resources", func() {
			By("Should create a settings", func() {
				Expect(LoadResource("", setting.Name, &v1beta1.Settings{})).To(Succeed())
			})
			By("Sould be able to retrieve the settings", func() {
				str, err := helperSettings.Get(TestContext(), stack.Name, helperSettings.SplitKeywordWithDot(setting.Spec.Key)...)
				Expect(err).To(BeNil())
				Expect(str).ToNot(BeNil())
				Expect(*str).To(Equal(setting.Spec.Value))
			})
		})
	})
})
