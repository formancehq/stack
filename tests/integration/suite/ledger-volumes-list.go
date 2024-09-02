package suite

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Transaction struct {
	Amount        int64
	Asset         string
	Source        string
	Destination   string
	EffectiveDate time.Time
}

var now = time.Now()

var transactions = []Transaction{
	Transaction{Amount: 100, Asset: "USD", Source: "world", Destination: "account:user1", EffectiveDate: now.Add(-4 * time.Hour)},        //user1:100, world:-100
	Transaction{Amount: 125, Asset: "USD", Source: "world", Destination: "account:user2", EffectiveDate: now.Add(-3 * time.Hour)},        //user1:100, user2:125, world:-225
	Transaction{Amount: 75, Asset: "USD", Source: "account:user1", Destination: "account:user2", EffectiveDate: now.Add(-2 * time.Hour)}, //user1:25, user2:200, world:-200
	Transaction{Amount: 175, Asset: "USD", Source: "world", Destination: "account:user1", EffectiveDate: now.Add(-1 * time.Hour)},        //user1:200, user2:200, world:-400
	Transaction{Amount: 50, Asset: "USD", Source: "account:user2", Destination: "bank", EffectiveDate: now},                              //user1:200, user2:150, world:-400, bank:50
	Transaction{Amount: 100, Asset: "USD", Source: "account:user2", Destination: "account:user1", EffectiveDate: now.Add(1 * time.Hour)}, //user1:300, user2:50, world:-400, bank:50
	Transaction{Amount: 150, Asset: "USD", Source: "account:user1", Destination: "bank", EffectiveDate: now.Add(2 * time.Hour)},          //user1:150, user2:50, world:-400, bank:200
}

var _ = WithModules([]*Module{modules.Ledger}, func() {

	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))

		for _, transaction := range transactions {
			response, err := Client().Ledger.V2.CreateTransaction(
				TestContext(),
				operations.V2CreateTransactionRequest{
					V2PostTransaction: shared.V2PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.V2Posting{
							{
								Amount:      big.NewInt(transaction.Amount),
								Asset:       transaction.Asset,
								Source:      transaction.Source,
								Destination: transaction.Destination,
							},
						},
						Timestamp: &transaction.EffectiveDate,
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
		}
	})

	When(fmt.Sprint("Get current Volumes and Balances From origin of time till now (insertion-date)"), func() {

		It("should be ok", func() {

			inserDate := true

			response, err := Client().Ledger.V2.GetVolumesWithBalances(
				TestContext(),
				operations.V2GetVolumesWithBalancesRequest{
					InsertionDate: &inserDate,
					Ledger:        "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			ret := response.V2VolumesWithBalanceCursorResponse
			Expect(len(ret.Cursor.Data)).To(Equal(4))
			for _, volume := range ret.Cursor.Data {
				if volume.Account == "account:user1" {
					Expect(volume.Balance).To(Equal(big.NewInt(150)))
				}
				if volume.Account == "account:user2" {
					Expect(volume.Balance).To(Equal(big.NewInt(50)))
				}

				if volume.Account == "bank" {
					Expect(volume.Balance).To(Equal(big.NewInt(200)))
				}

				if volume.Account == "world" {
					Expect(volume.Balance).To(Equal(big.NewInt(-400)))
				}

			}

		})

	})

	When(fmt.Sprint("Get Volumes and Balances From oot til oot+2 hours (effectiveDate) "), func() {

		It("should be ok", func() {

			startTime := now.Add(-4 * time.Hour)
			stopTime := now.Add(-2 * time.Hour)

			response, err := Client().Ledger.V2.GetVolumesWithBalances(
				TestContext(),
				operations.V2GetVolumesWithBalancesRequest{
					StartTime: &startTime,
					EndTime:   &stopTime,
					Ledger:    "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			ret := response.V2VolumesWithBalanceCursorResponse
			Expect(len(ret.Cursor.Data)).To(Equal(3))
			for _, volume := range ret.Cursor.Data {
				//fmt.Println(fmt.Sprintf("%s | %s | %d | %d | %d", volume.Account, volume.Asset, volume.Input, volume.Output, volume.Balance))

				if volume.Account == "account:user1" {
					Expect(volume.Balance).To(Equal(big.NewInt(25)))
				}
				if volume.Account == "account:user2" {
					Expect(volume.Balance).To(Equal(big.NewInt(200)))
				}

				if volume.Account == "world" {
					Expect(volume.Balance).To(Equal(big.NewInt(-225)))
				}

			}

		})

	})

	When(fmt.Sprint("Get Volumes and Balances Filter by address account"), func() {

		It("should be ok", func() {
			inserDate := true
			response, err := Client().Ledger.V2.GetVolumesWithBalances(
				TestContext(),
				operations.V2GetVolumesWithBalancesRequest{
					InsertionDate: &inserDate,
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"account": "account:",
						},
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			ret := response.V2VolumesWithBalanceCursorResponse
			Expect(len(ret.Cursor.Data)).To(Equal(2))
			for _, volume := range ret.Cursor.Data {
				//fmt.Println(fmt.Sprintf("%s | %s | %d | %d | %d", volume.Account, volume.Asset, volume.Input, volume.Output, volume.Balance))

				if volume.Account == "account:user1" {
					Expect(volume.Balance).To(Equal(big.NewInt(150)))
				}
				if volume.Account == "account:user2" {
					Expect(volume.Balance).To(Equal(big.NewInt(50)))
				}

			}

		})

	})

	When(fmt.Sprint("Get Volumes and Balances Filter by address account a,d and end-time now effective"), func() {

		It("should be ok", func() {

			response, err := Client().Ledger.V2.GetVolumesWithBalances(
				TestContext(),
				operations.V2GetVolumesWithBalancesRequest{

					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"account": "account:",
						},
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			ret := response.V2VolumesWithBalanceCursorResponse
			Expect(len(ret.Cursor.Data)).To(Equal(2))
			for _, volume := range ret.Cursor.Data {
				//fmt.Println(fmt.Sprintf("%s | %s | %d | %d | %d", volume.Account, volume.Asset, volume.Input, volume.Output, volume.Balance))
				if volume.Account == "account:user1" {
					Expect(volume.Balance).To(Equal(big.NewInt(200)))
				}
				if volume.Account == "account:user2" {
					Expect(volume.Balance).To(Equal(big.NewInt(150)))
				}

			}

		})

	})

	When(fmt.Sprint("Get Volumes and Balances Filter by address account which doesn't exist"), func() {

		It("should be ok", func() {

			response, err := Client().Ledger.V2.GetVolumesWithBalances(
				TestContext(),
				operations.V2GetVolumesWithBalancesRequest{

					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"account": "foo:",
						},
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			ret := response.V2VolumesWithBalanceCursorResponse
			Expect(len(ret.Cursor.Data)).To(Equal(0))

		})

	})

	When(fmt.Sprint("Get Volumes and Balances Filter With futures dates empty"), func() {

		It("should be ok", func() {

			startDate := time.Now().Add(8 * time.Hour)
			endDate := time.Now().Add(12 * time.Hour)

			response, err := Client().Ledger.V2.GetVolumesWithBalances(
				TestContext(),
				operations.V2GetVolumesWithBalancesRequest{
					StartTime: &startDate,
					EndTime:   &endDate,
					Ledger:    "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			ret := response.V2VolumesWithBalanceCursorResponse
			Expect(len(ret.Cursor.Data)).To(Equal(0))

		})

	})

	When(fmt.Sprint("Get Volumes and Balances Filter by address account aggregation by level 1"), func() {

		It("should be ok", func() {
			inserDate := true
			groupBylvl := int64(1)
			response, err := Client().Ledger.V2.GetVolumesWithBalances(
				TestContext(),
				operations.V2GetVolumesWithBalancesRequest{
					InsertionDate: &inserDate,
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"account": "account:",
						},
					},
					GroupBy: &groupBylvl,
					Ledger:  "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			ret := response.V2VolumesWithBalanceCursorResponse
			Expect(len(ret.Cursor.Data)).To(Equal(1))
			for _, volume := range ret.Cursor.Data {
				//fmt.Println(fmt.Sprintf("%s | %s | %d | %d | %d", volume.Account, volume.Asset, volume.Input, volume.Output, volume.Balance))

				if volume.Account == "account" {
					Expect(volume.Balance).To(Equal(big.NewInt(200)))
				}

			}

		})

	})

})
