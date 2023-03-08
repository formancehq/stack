package suite

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("An environment configured with a webhook sent on created transaction", func() {
	var (
		testingHttpServer *httptest.Server
		now               = time.Now().Round(time.Second).UTC()
		called            chan struct{}
	)
	BeforeEach(func() {
		called = make(chan struct{})
		testingHttpServer = httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				close(called)
			}))
		_, _, err := Client().WebhooksApi.
			InsertConfig(TestContext()).ConfigUser(formance.ConfigUser{
			Endpoint: testingHttpServer.URL,
			EventTypes: []string{
				"ledger.committed_transactions",
			},
		}).
			Execute()
		Expect(err).To(Not(HaveOccurred()))
	})
	AfterEach(func() {
		testingHttpServer.Close()
	})
	When("creating a transaction", func() {
		BeforeEach(func() {
			_, _, err := Client().TransactionsApi.
				CreateTransaction(TestContext(), "default").
				PostTransaction(formance.PostTransaction{
					Timestamp: &now,
					Postings: []formance.Posting{{
						Amount:      100,
						Asset:       "USD",
						Source:      "world",
						Destination: "alice",
					}},
				}).
				Execute()
			Expect(err).To(Not(HaveOccurred()))
		})
		It("should trigger a call to the webhook endpoint", func() {
			Eventually(ChanClosed(called)).Should(BeTrue())
		})
	})
})
