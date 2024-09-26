package suite

import (
	"database/sql"
	"math/big"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/formancehq/stack/tests/integration/internal/modules"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	webhooks "github.com/formancehq/webhooks/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var _ = WithModules([]*Module{modules.Ledger, modules.Webhooks}, func() {
	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
	})
	Context("the endpoint only returning errors", func() {
		It("with an exponential backoff starting at 1s with a 3s timeout, 3 attempts have to be made and all should have a failed status", func() {
			httpServer := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, _ *http.Request) {
					http.Error(w, "error", http.StatusNotFound)
				}))
			defer func() {
				httpServer.Close()
			}()
			sqldb := sql.OpenDB(
				pgdriver.NewConnector(
					pgdriver.WithDSN(CurrentTest().GetDatabaseSourceName("webhooks"))))
			db := bun.NewDB(sqldb, pgdialect.New())
			defer func() {
				_ = db.Close()
			}()

			response, err := Client().Webhooks.V1.InsertConfig(
				TestContext(),
				shared.ConfigUser{
					Endpoint: httpServer.URL,
					EventTypes: []string{
						"ledger.committed_transactions",
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			createTransactionResponse, err := Client().Ledger.V2.CreateTransaction(
				TestContext(),
				operations.V2CreateTransactionRequest{
					V2PostTransaction: shared.V2PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.V2Posting{
							{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Source:      "world",
								Destination: "alice",
							},
						},
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(createTransactionResponse.StatusCode).To(Equal(http.StatusOK))

			Eventually(db.Ping()).
				WithTimeout(5 * time.Second).Should(Succeed())

			Eventually(getNumAttemptsToRetry).WithArguments(db).
				WithTimeout(5 * time.Second).
				Should(BeNumerically(">", 0))

			Eventually(getNumFailedAttempts).WithArguments(db).
				WithTimeout(5 * time.Second).
				Should(BeNumerically(">=", 3))

			toRetry, err := getNumAttemptsToRetry(db)
			Expect(err).ToNot(HaveOccurred())
			Expect(toRetry).To(Equal(0))
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

func getNumAttemptsToRetry(db *bun.DB) (int, error) {
	var results []webhooks.Attempt
	err := db.NewSelect().Model(&results).
		Where("status = ?", "to retry").
		Scan(TestContext())
	if err != nil {
		return 0, err
	}
	return len(results), nil
}

func getNumFailedAttempts(db *bun.DB) (int, error) {
	var results []webhooks.Attempt
	err := db.NewSelect().Model(&results).
		Where("status = ?", "failed").
		Scan(TestContext())
	if err != nil {
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

func getAttemptsStatus(db *bun.DB) (string, error) {
	atts, err := getAttempts(db)
	if err != nil {
		return "", err
	}
	if len(atts) == 0 {
		return "", nil
	}
	return atts[0].Status, nil
}
