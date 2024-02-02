package tests_test

import (
	. "github.com/formancehq/operator/internal/tests/internal"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/api/stack.formance.com/v1beta3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("VersionsController (legacy)", func() {
	When("Creating a legacy Versions", func() {
		var (
			oldVersions *v1beta3.Versions
		)
		BeforeEach(func() {
			oldVersions = &v1beta3.Versions{
				ObjectMeta: RandObjectMeta(),
				Spec: v1beta3.VersionsSpec{
					Ledger: "1234",
				},
			}
			Expect(Create(oldVersions)).To(Succeed())
		})
		AfterEach(func() {
			Expect(Delete(oldVersions)).To(Succeed())
		})
		It("Should create a new Versions object with the correct versions", func() {
			version := &v1beta1.Versions{}
			Eventually(func(g Gomega) error {
				return LoadResource("", oldVersions.Name, version)
			}).Should(Succeed())
			Expect(version.Spec).To(HaveLen(1))
			Expect(version.Spec["ledger"]).To(Equal("1234"))
		})
	})
})
