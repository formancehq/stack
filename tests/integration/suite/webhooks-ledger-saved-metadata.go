package suite

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/security"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("An environment configured with a webhook sent on saved metadata", func() {
	var (
		httpServer *httptest.Server
		called     chan struct{}
		secret     = webhooks.NewSecret()
	)

	BeforeEach(func() {
		called = make(chan struct{})
		httpServer = httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer close(called)
				id := r.Header.Get("formance-webhook-id")
				ts := r.Header.Get("formance-webhook-timestamp")
				signatures := r.Header.Get("formance-webhook-signature")
				timeInt, err := strconv.ParseInt(ts, 10, 64)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				payload, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				ok, err := security.Verify(signatures, id, timeInt, secret, payload)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if !ok {
					http.Error(w, "WEBHOOKS SIGNATURE VERIFICATION NOK", http.StatusBadRequest)
					return
				}
			}))
		DeferCleanup(func() {
			httpServer.Close()
		})

		_, _, err := Client().WebhooksApi.
			InsertConfig(TestContext()).ConfigUser(formance.ConfigUser{
			Endpoint: httpServer.URL,
			Secret:   &secret,
			EventTypes: []string{
				"ledger.saved_metadata",
			},
		}).Execute()
		Expect(err).ToNot(HaveOccurred())
	})

	When("creating a transaction and adding metadata to it", func() {
		BeforeEach(func() {
			insertResp, _, err := Client().TransactionsApi.
				CreateTransaction(TestContext(), "default").
				PostTransaction(formance.PostTransaction{
					Postings: []formance.Posting{{
						Amount:      100,
						Asset:       "USD",
						Source:      "world",
						Destination: "alice",
					}},
				}).
				Execute()
			Expect(err).ToNot(HaveOccurred())

			_, err = Client().TransactionsApi.AddMetadataOnTransaction(TestContext(),
				"default", insertResp.Data.Txid).
				RequestBody(map[string]interface{}{"key": "value"}).Execute()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should trigger a call to the webhook endpoint", func() {
			Eventually(ChanClosed(called)).Should(BeTrue())
		})
	})
})
