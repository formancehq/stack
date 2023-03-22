package suite

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("An empty environment", func() {
	When("listing workflows", func() {
		var (
			ret *shared.ListWorkflowsResponse
		)
		BeforeEach(func() {
			response, err := Client().Orchestration.ListWorkflows(
				TestContext(),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			ret = response.ListWorkflowsResponse
		})
		It("should respond with an empty list", func() {
			Expect(ret.Data).To(BeEmpty())
		})
	})
})
