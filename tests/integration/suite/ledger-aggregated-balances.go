package suite

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/big"
)

var _ = WithModules([]*Module{modules.Ledger, modules.Search}, func() {
	When("creating two transactions on a ledger with custom metadata", func() {
		BeforeEach(func() {
			_, err := Client().Ledger.V2.CreateBulk(TestContext(), operations.CreateBulkRequest{
				RequestBody: []shared.BulkElement{
					shared.CreateBulkElementCreateTransaction(shared.BulkElementBulkElementCreateTransaction{
						Data: &shared.PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD/2",
								Destination: "bank1",
								Source:      "world",
							}},
						},
					}),
					shared.CreateBulkElementCreateTransaction(shared.BulkElementBulkElementCreateTransaction{
						Data: &shared.PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.Posting{{
								Amount:      big.NewInt(100),
								Asset:       "USD/2",
								Destination: "bank2",
								Source:      "world",
							}},
						},
					}),
					shared.CreateBulkElementAddMetadata(shared.BulkElementBulkElementAddMetadata{
						Data: &shared.BulkElementBulkElementAddMetadataData{
							Metadata: map[string]string{
								"category": "premium",
							},
							TargetID:   shared.CreateTargetIDStr("bank2"),
							TargetType: shared.TargetTypeAccount,
						},
					}),
					shared.CreateBulkElementAddMetadata(shared.BulkElementBulkElementAddMetadata{
						Data: &shared.BulkElementBulkElementAddMetadataData{
							Metadata: map[string]string{
								"category": "premium",
							},
							TargetID:   shared.CreateTargetIDStr("bank1"),
							TargetType: shared.TargetTypeAccount,
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
				operations.GetBalancesAggregatedRequest{
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

			Expect(response.AggregateBalancesResponse.Data).To(HaveLen(1))
			Expect(response.AggregateBalancesResponse.Data["USD/2"]).To(Equal(big.NewInt(200)))
		})
	})
})
