package suite

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/big"
)

var _ = WithModules([]*Module{modules.Auth, modules.Ledger, modules.Wallets}, func() {

	When("creating a wallet", func() {
		var (
			createWalletResponse *operations.CreateWalletResponse
			err                  error
		)
		BeforeEach(func() {
			createWalletResponse, err = Client().Wallets.CreateWallet(
				TestContext(),
				&shared.CreateWalletRequest{
					Name:     uuid.NewString(),
					Metadata: map[string]string{},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(createWalletResponse.StatusCode).To(Equal(201))
		})
		Then("crediting it", func() {
			BeforeEach(func() {
				_, err := Client().Wallets.CreditWallet(TestContext(), operations.CreditWalletRequest{
					CreditWalletRequest: &shared.CreditWalletRequest{
						Amount: shared.Monetary{
							Amount: big.NewInt(1000),
							Asset:  "USD/2",
						},
						Sources:  []shared.Subject{},
						Metadata: map[string]string{},
					},
					ID: createWalletResponse.CreateWalletResponse.Data.ID,
				})
				Expect(err).To(Succeed())
			})
			Then("debiting it", func() {
				BeforeEach(func() {
					_, err := Client().Wallets.DebitWallet(TestContext(), operations.DebitWalletRequest{
						DebitWalletRequest: &shared.DebitWalletRequest{
							Amount: shared.Monetary{
								Amount: big.NewInt(100),
								Asset:  "USD/2",
							},
							Metadata: map[string]string{},
						},
						ID: createWalletResponse.CreateWalletResponse.Data.ID,
					})
					Expect(err).To(Succeed())
				})
				It("should be ok", func() {})
			})
			Then("debiting it using a hold", func() {
				var (
					debitWalletResponse *operations.DebitWalletResponse
				)
				BeforeEach(func() {
					debitWalletResponse, err = Client().Wallets.DebitWallet(TestContext(), operations.DebitWalletRequest{
						DebitWalletRequest: &shared.DebitWalletRequest{
							Amount: shared.Monetary{
								Amount: big.NewInt(100),
								Asset:  "USD/2",
							},
							Pending:  pointer.For(true),
							Metadata: map[string]string{},
						},
						ID: createWalletResponse.CreateWalletResponse.Data.ID,
					})
					Expect(err).To(Succeed())
				})
				It("should be ok", func() {})
				Then("confirm the hold", func() {
					BeforeEach(func() {
						_, err := Client().Wallets.ConfirmHold(TestContext(), operations.ConfirmHoldRequest{
							HoldID: debitWalletResponse.DebitWalletResponse.Data.ID,
						})
						Expect(err).To(Succeed())
					})
					It("should be ok", func() {})
				})
				Then("void the hold", func() {
					BeforeEach(func() {
						_, err := Client().Wallets.VoidHold(TestContext(), operations.VoidHoldRequest{
							HoldID: debitWalletResponse.DebitWalletResponse.Data.ID,
						})
						Expect(err).To(Succeed())
					})
					It("should be ok", func() {})
				})
			})
		})
	})
})
