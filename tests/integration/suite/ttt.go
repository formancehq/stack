package suite

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("some empty environment", func() {
	When("creating a transaction on a ledger", func() {
		It("should fail", func() {
			resp, httpResp, err := Client().TransactionsApi.
				CreateTransaction(TestContext(), "default").
				PostTransaction(formance.PostTransaction{
					Postings: []formance.Posting{{
						Amount:      100,
						Asset:       "USD",
						Source:      "bob",
						Destination: "alice",
					}},
				}).
				Execute()
			Expect(err).To(HaveOccurred())
			spew.Dump("TXRESP", resp, "HTTPRESP", httpResp, "TXERR", err)
		})
	})
})
