package suite

import (
	"github.com/formancehq/stack/libs/go-libs/pointer"
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

var _ = WithModules([]*Module{modules.Ledger}, func() {
	now := time.Now().UTC().Round(time.Second)
	When("creating two transactions on a ledger with custom metadata", func() {
		var firstTransactionsInsertedAt time.Time
		BeforeEach(func() {
			response, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
				Ledger: "default",
			})
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusNoContent))

			_, err = Client().Ledger.V2.CreateBulk(TestContext(), operations.V2CreateBulkRequest{
				RequestBody: []shared.V2BulkElement{
					shared.CreateV2BulkElementCreateTransaction(shared.V2BulkElementCreateTransaction{
						Data: &shared.V2PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD/2",
								Destination: "bank1",
								Source:      "world",
							}},
							Timestamp: pointer.For(now.Add(-time.Minute)),
						},
					}),
					shared.CreateV2BulkElementCreateTransaction(shared.V2BulkElementCreateTransaction{
						Data: &shared.V2PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD/2",
								Destination: "bank2",
								Source:      "world",
							}},
							Timestamp: pointer.For(now.Add(-2 * time.Minute)),
						},
					}),
					shared.CreateV2BulkElementCreateTransaction(shared.V2BulkElementCreateTransaction{
						Data: &shared.V2PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD/2",
								Destination: "bank1",
								Source:      "world",
							}},
							Timestamp: pointer.For(now),
						},
					}),
					shared.CreateV2BulkElementAddMetadata(shared.V2BulkElementAddMetadata{
						Data: &shared.V2BulkElementAddMetadataData{
							Metadata: map[string]string{
								"category": "premium",
							},
							TargetID:   shared.CreateV2TargetIDStr("bank2"),
							TargetType: shared.V2TargetTypeAccount,
						},
					}),
					shared.CreateV2BulkElementAddMetadata(shared.V2BulkElementAddMetadata{
						Data: &shared.V2BulkElementAddMetadataData{
							Metadata: map[string]string{
								"category": "premium",
							},
							TargetID:   shared.CreateV2TargetIDStr("bank1"),
							TargetType: shared.V2TargetTypeAccount,
						},
					}),
				},
				Ledger: "default",
			})
			Expect(err).To(Succeed())

			firstTransactionsInsertedAt = time.Now()
			<-time.After(time.Second)

			_, err = Client().Ledger.V2.CreateBulk(TestContext(), operations.V2CreateBulkRequest{
				RequestBody: []shared.V2BulkElement{
					shared.CreateV2BulkElementCreateTransaction(shared.V2BulkElementCreateTransaction{
						Data: &shared.V2PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD/2",
								Destination: "bank1",
								Source:      "world",
							}},
							Timestamp: pointer.For(now),
						},
					}),
				},
				Ledger: "default",
			})
			Expect(err).To(Succeed())
		})
		It("should be ok when aggregating using the metadata", func() {
			response, err := Client().Ledger.V2.GetBalancesAggregated(
				TestContext(),
				operations.V2GetBalancesAggregatedRequest{
					RequestBody: map[string]any{
						"$match": map[string]any{
							"metadata[category]": "premium",
						},
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			Expect(response.V2AggregateBalancesResponse.Data).To(HaveLen(1))
			Expect(response.V2AggregateBalancesResponse.Data["USD/2"]).To(Equal(big.NewInt(400)))
		})
		It("should be ok when aggregating using pit on effective date", func() {
			response, err := Client().Ledger.V2.GetBalancesAggregated(
				TestContext(),
				operations.V2GetBalancesAggregatedRequest{
					Ledger: "default",
					Pit:    pointer.For(now.Add(-time.Minute)),
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"address": "bank1",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			Expect(response.V2AggregateBalancesResponse.Data).To(HaveLen(1))
			Expect(response.V2AggregateBalancesResponse.Data["USD/2"]).To(Equal(big.NewInt(100)))
		})
		It("should be ok when aggregating using pit on insertion date", func() {
			response, err := Client().Ledger.V2.GetBalancesAggregated(
				TestContext(),
				operations.V2GetBalancesAggregatedRequest{
					Ledger:           "default",
					Pit:              pointer.For(firstTransactionsInsertedAt),
					UseInsertionDate: pointer.For(true),
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"address": "bank1",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			Expect(response.V2AggregateBalancesResponse.Data).To(HaveLen(1))
			Expect(response.V2AggregateBalancesResponse.Data["USD/2"]).To(Equal(big.NewInt(200)))
		})
	})
})
