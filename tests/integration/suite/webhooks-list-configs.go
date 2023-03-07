package suite

import (
	"context"

	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("some empty environment", func() {
	When("listing configs on wehbooks", func() {
		var (
			response *formance.ConfigsResponse
			err      error
		)
		BeforeEach(func() {
			response, _, err = Client().WebhooksApi.
				GetManyConfigs(context.Background()).
				Execute()
			Expect(err).To(BeNil())
		})

		It("should return no config", func() {
			Expect(response.Cursor.Data).To(BeEmpty())
		})
	})
})
