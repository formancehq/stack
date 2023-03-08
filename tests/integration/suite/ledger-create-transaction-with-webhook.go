package suite

import (
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/webhooks/cmd/flag"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/security"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var _ = Given("An environment configured with a webhook sent on created transaction", func() {
	var (
		httpServer *httptest.Server
		now        = time.Now().Round(time.Second).UTC()
		called     chan struct{}
		secret     = webhooks.NewSecret()
	)

	BeforeEach(func() {
		called = make(chan struct{})
		httpServer = httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					close(called)
				}()
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
				"ledger.committed_transactions",
			},
		}).
			Execute()
		Expect(err).To(BeNil())
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
			Expect(err).To(BeNil())
		})

		It("should trigger a call to the webhook endpoint", func() {
			Eventually(ChanClosed(called)).Should(BeTrue())
		})

		It("should insert a successful attempt in DB", func() {
			sqldb := sql.OpenDB(
				pgdriver.NewConnector(
					pgdriver.WithDSN(viper.GetString(flag.StoragePostgresConnString))))
			db := bun.NewDB(sqldb, pgdialect.New())
			defer db.Close()

			var attempts []webhooks.Attempt
			Expect(db.NewSelect().Model(&attempts).Scan(TestContext())).To(Succeed())
			Expect(attempts).To(HaveLen(1))
			Expect(attempts[0].Status).To(Equal("success"))
		})
	})
})
