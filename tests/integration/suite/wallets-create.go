package suite

import (
	"fmt"
	"github.com/formancehq/go-libs/pointer"
	"math/big"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Auth, modules.Ledger, modules.Wallets}, func() {
	const countWallets = 3
	When(fmt.Sprintf("creating %d wallets", countWallets), func() {
		JustBeforeEach(func() {
			for i := 0; i < countWallets; i++ {
				name := uuid.NewString()
				response, err := Client().Wallets.V1.CreateWallet(
					TestContext(),
					operations.CreateWalletRequest{
						CreateWalletRequest: &shared.CreateWalletRequest{
							Metadata: map[string]string{
								"wallets_number": fmt.Sprint(i),
							},
							Name: name,
						},
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(201))

				_, err = Client().Wallets.V1.CreditWallet(TestContext(), operations.CreditWalletRequest{
					CreditWalletRequest: &shared.CreditWalletRequest{
						Amount: shared.Monetary{
							Amount: big.NewInt(100),
							Asset:  "USD",
						},
					},
					ID: response.CreateWalletResponse.Data.ID,
				})
				Expect(err).ToNot(HaveOccurred())
			}
		})
		Then("listing them", func() {
			var (
				request  operations.ListWalletsRequest
				response *operations.ListWalletsResponse
				err      error
			)
			BeforeEach(func() {
				// reset between each test
				request = operations.ListWalletsRequest{}
			})
			JustBeforeEach(func() {
				Eventually(func(g Gomega) bool {
					response, err = Client().Wallets.V1.ListWallets(TestContext(), request)
					g.Expect(err).ToNot(HaveOccurred())
					g.Expect(response.StatusCode).To(Equal(200))
					return true
				}).Should(BeTrue())
			})
			It(fmt.Sprintf("should return %d items", countWallets), func() {
				Expect(response.ListWalletsResponse.Cursor.Data).To(HaveLen(countWallets))
			})
			Context("using a metadata filter", func() {
				BeforeEach(func() {
					request.Metadata = map[string]string{
						"wallets_number": "0",
					}
				})
				It("should return only one item", func() {
					Expect(response.ListWalletsResponse.Cursor.Data).To(HaveLen(1))
				})
			})
			Context("expanding balances", func() {
				BeforeEach(func() {
					request.Expand = pointer.For("balances")
				})
				It("should return all items with volumes and balances", func() {
					Expect(response.ListWalletsResponse.Cursor.Data).To(HaveLen(3))
					for _, wallet := range response.ListWalletsResponse.Cursor.Data {
						Expect(wallet.Balances.Main.Assets["USD"]).To(Equal(big.NewInt(100)))
					}
				})
			})
		})
	})
})
