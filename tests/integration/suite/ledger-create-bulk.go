package suite

import (
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Ledger}, func() {
	BeforeEach(func() {
		response, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(response.StatusCode).To(Equal(http.StatusNoContent))
	})
	When("creating a bulk on a ledger", func() {
		var (
			now = time.Now().Round(time.Microsecond).UTC()
		)
		BeforeEach(func() {
			_, err := Client().Ledger.V2.CreateBulk(TestContext(), operations.V2CreateBulkRequest{
				RequestBody: []shared.V2BulkElement{
					shared.CreateV2BulkElementCreateTransaction(shared.V2BulkElementCreateTransaction{
						Data: &shared.V2PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD/2",
								Destination: "bank",
								Source:      "world",
							}},
							Timestamp: &now,
						},
					}),
					shared.CreateV2BulkElementAddMetadata(shared.V2BulkElementAddMetadata{
						Data: &shared.V2BulkElementAddMetadataData{
							Metadata: metadata.Metadata{
								"foo":  "bar",
								"role": "admin",
							},
							TargetID:   shared.CreateV2TargetIDBigint(big.NewInt(0)),
							TargetType: shared.V2TargetTypeTransaction,
						},
					}),
					shared.CreateV2BulkElementDeleteMetadata(shared.V2BulkElementDeleteMetadata{
						Data: &shared.V2BulkElementDeleteMetadataData{
							Key:        "foo",
							TargetID:   shared.CreateV2TargetIDBigint(big.NewInt(0)),
							TargetType: shared.V2TargetTypeTransaction,
						},
					}),
					shared.CreateV2BulkElementRevertTransaction(shared.V2BulkElementRevertTransaction{
						Data: &shared.V2BulkElementRevertTransactionData{
							ID: big.NewInt(0),
						},
					}),
				},
				Ledger: "default",
			})
			Expect(err).To(Succeed())
		})
		It("should be ok", func() {
			tx, err := Client().Ledger.V2.GetTransaction(TestContext(), operations.V2GetTransactionRequest{
				ID:     big.NewInt(0),
				Ledger: "default",
			})
			Expect(err).To(Succeed())
			Expect(tx.V2GetTransactionResponse.Data).To(Equal(shared.V2ExpandedTransaction{
				ID: big.NewInt(0),
				Metadata: metadata.Metadata{
					"role": "admin",
				},
				Postings: []shared.V2Posting{{
					Amount:      big.NewInt(100),
					Asset:       "USD/2",
					Destination: "bank",
					Source:      "world",
				}},
				Reverted:  true,
				Timestamp: now,
			}))
		})
	})
	When("creating a bulk with an error on a ledger", func() {
		var (
			now          = time.Now().Round(time.Microsecond).UTC()
			err          error
			bulkResponse *operations.V2CreateBulkResponse
		)
		BeforeEach(func() {
			bulkResponse, err = Client().Ledger.V2.CreateBulk(TestContext(), operations.V2CreateBulkRequest{
				RequestBody: []shared.V2BulkElement{
					shared.CreateV2BulkElementCreateTransaction(shared.V2BulkElementCreateTransaction{
						Data: &shared.V2PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD/2",
								Destination: "bank",
								Source:      "world",
							}},
							Timestamp: &now,
						},
					}),
					shared.CreateV2BulkElementCreateTransaction(shared.V2BulkElementCreateTransaction{
						Data: &shared.V2PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.V2Posting{{
								Amount:      big.NewInt(200), // Insufficient fund
								Asset:       "USD/2",
								Destination: "user",
								Source:      "bank",
							}},
							Timestamp: &now,
						},
					}),
				},
				Ledger: "default",
			})
			Expect(err).To(Succeed())
		})
		It("should respond with an error", func() {
			Expect(bulkResponse.V2BulkResponse.Data[1].Type).To(Equal(shared.V2BulkElementResultType("ERROR")))
			Expect(bulkResponse.V2BulkResponse.Data[1].V2BulkElementResultErrorSchemas.ErrorCode).To(Equal("INSUFFICIENT_FUND"))
		})
	})
})
