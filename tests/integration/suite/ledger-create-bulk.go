package suite

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/big"
	"net/http"
	"time"
)

var _ = WithModules([]*Module{modules.Ledger}, func() {
	BeforeEach(func() {
		response, err := Client().Ledger.CreateLedger(TestContext(), operations.CreateLedgerRequest{
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
			_, err := Client().Ledger.CreateBulk(TestContext(), operations.CreateBulkRequest{
				RequestBody: []shared.BulkElement{
					shared.CreateBulkElementCreateTransaction(shared.BulkElementBulkElementCreateTransaction{
						Data: &shared.PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD/2",
								Destination: "bank",
								Source:      "world",
							}},
							Timestamp: &now,
						},
					}),
					shared.CreateBulkElementAddMetadata(shared.BulkElementBulkElementAddMetadata{
						Data: &shared.BulkElementBulkElementAddMetadataData{
							Metadata: metadata.Metadata{
								"foo":  "bar",
								"role": "admin",
							},
							TargetID:   shared.CreateTargetIDBigint(big.NewInt(0)),
							TargetType: shared.TargetTypeTransaction,
						},
					}),
					shared.CreateBulkElementDeleteMetadata(shared.BulkElementBulkElementDeleteMetadata{
						Data: &shared.BulkElementBulkElementDeleteMetadataData{
							Key:        "foo",
							TargetID:   shared.CreateTargetIDBigint(big.NewInt(0)),
							TargetType: shared.TargetTypeTransaction,
						},
					}),
					shared.CreateBulkElementRevertTransaction(shared.BulkElementBulkElementRevertTransaction{
						Data: &shared.BulkElementBulkElementRevertTransactionData{
							ID: big.NewInt(0),
						},
					}),
				},
				Ledger: "default",
			})
			Expect(err).To(Succeed())
		})
		It("should be ok", func() {
			tx, err := Client().Ledger.GetTransaction(TestContext(), operations.GetTransactionRequest{
				ID:     big.NewInt(0),
				Ledger: "default",
			})
			Expect(err).To(Succeed())
			Expect(tx.GetTransactionResponse.Data).To(Equal(shared.ExpandedTransaction{
				ID: big.NewInt(0),
				Metadata: metadata.Metadata{
					"role": "admin",
				},
				Postings: []shared.Posting{{
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
			bulkResponse *operations.CreateBulkResponse
		)
		BeforeEach(func() {
			bulkResponse, err = Client().Ledger.CreateBulk(TestContext(), operations.CreateBulkRequest{
				RequestBody: []shared.BulkElement{
					shared.CreateBulkElementCreateTransaction(shared.BulkElementBulkElementCreateTransaction{
						Data: &shared.PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD/2",
								Destination: "bank",
								Source:      "world",
							}},
							Timestamp: &now,
						},
					}),
					shared.CreateBulkElementCreateTransaction(shared.BulkElementBulkElementCreateTransaction{
						Data: &shared.PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.Posting{{
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
			Expect(bulkResponse.BulkResponse.Data[1].Type).To(Equal(shared.BulkElementResultType("ERROR")))
			Expect(bulkResponse.BulkResponse.Data[1].BulkElementResultBulkElementResultError.ErrorCode).To(Equal("INSUFFICIENT_FUND"))
		})
	})
})
