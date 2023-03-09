package suite

import (
	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("An empty environment", func() {
	When("listing workflows", func() {
		var (
			ret *formance.ListWorkflowsResponse
			err error
		)
		BeforeEach(func() {
			ret, _, err = Client().OrchestrationApi.
				ListWorkflows(TestContext()).
				Execute()
			Expect(err).To(BeNil())
		})
		It("should respond with an empty list", func() {
			Expect(ret.Data).To(BeEmpty())
		})
	})
})
