package suite

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/webhooks/cmd/flag"
	webhooks "github.com/formancehq/webhooks/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var _ = Given("an environment configured with a webhook sent on created transaction", func() {
	Context("the endpoint only returning errors", func() {
		It("with a retries schedule of [1s,2s], 3 attempts have to be made and all should have a failed status", func() {
			httpServer := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, _ *http.Request) {
					http.Error(w, "error", http.StatusNotFound)
				}))
			sqldb := sql.OpenDB(
				pgdriver.NewConnector(
					pgdriver.WithDSN(viper.GetString(flag.StoragePostgresConnString))))
			db := bun.NewDB(sqldb, pgdialect.New())
			DeferCleanup(func() {
				httpServer.Close()
				Expect(db.Close()).To(Succeed())
			})

			_, _, err := Client().WebhooksApi.
				InsertConfig(TestContext()).ConfigUser(formance.ConfigUser{
				Endpoint: httpServer.URL,
				EventTypes: []string{
					"ledger.committed_transactions",
				},
			}).Execute()
			Expect(err).ToNot(HaveOccurred())

			_, _, err = Client().TransactionsApi.
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

			Eventually(db.Ping()).
				WithTimeout(5 * time.Second).Should(Succeed())

			Eventually(getNumAttempts).WithArguments(db).
				WithTimeout(5 * time.Second).
				Should(BeNumerically("==", 3))

			att, err := getAttempts(db)
			Expect(err).ToNot(HaveOccurred())

			for _, a := range att {
				Expect(a.Status).To(Equal("failed"))
			}
		})
	})
})

func getNumAttempts(db *bun.DB) (int, error) {
	var results []webhooks.Attempt
	if err := db.NewSelect().Model(&results).Scan(TestContext()); err != nil {
		return 0, err
	}
	return len(results), nil
}

func getAttempts(db *bun.DB) ([]webhooks.Attempt, error) {
	var results []webhooks.Attempt
	if err := db.NewSelect().Model(&results).Scan(TestContext()); err != nil {
		return []webhooks.Attempt{}, err
	}
	return results, nil
}
