package suite

import (
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/pointer"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"math/big"
	"time"
)

var _ = WithModules([]*Module{modules.Auth, modules.Ledger, modules.Wallets}, func() {

	When("creating a wallet", func() {
		var (
			response *operations.CreateWalletResponse
			err      error
		)
		BeforeEach(func() {
			response, err = Client().Wallets.V1.CreateWallet(
				TestContext(),
				operations.CreateWalletRequest{
					CreateWalletRequest: &shared.CreateWalletRequest{
						Name:     uuid.NewString(),
						Metadata: map[string]string{},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(201))
		})
		Then("crediting it", func() {
			BeforeEach(func() {
				_, err := Client().Wallets.V1.CreditWallet(TestContext(), operations.CreditWalletRequest{
					CreditWalletRequest: &shared.CreditWalletRequest{
						Amount: shared.Monetary{
							Amount: big.NewInt(1000),
							Asset:  "USD/2",
						},
						Sources:  []shared.Subject{},
						Metadata: map[string]string{},
					},
					ID:             response.CreateWalletResponse.Data.ID,
					IdempotencyKey: pointer.For("foo"),
				})
				Expect(err).To(Succeed())
			})
			It("should be ok", func() {})
			Then("crediting again with the same ik", func() {
				BeforeEach(func() {
					_, err := Client().Wallets.V1.CreditWallet(TestContext(), operations.CreditWalletRequest{
						CreditWalletRequest: &shared.CreditWalletRequest{
							Amount: shared.Monetary{
								Amount: big.NewInt(1000),
								Asset:  "USD/2",
							},
							Sources:  []shared.Subject{},
							Metadata: map[string]string{},
						},
						ID:             response.CreateWalletResponse.Data.ID,
						IdempotencyKey: pointer.For("foo"),
					})
					Expect(err).To(Succeed())
				})
				It("Should not trigger any movements", func() {
					balance, err := Client().Wallets.V1.GetBalance(TestContext(), operations.GetBalanceRequest{
						BalanceName: "main",
						ID:          response.CreateWalletResponse.Data.ID,
					})
					Expect(err).To(Succeed())
					Expect(balance.GetBalanceResponse.Data.Assets["USD/2"]).To(Equal(big.NewInt(1000)))
				})
			})
		})
		Then("crediting it with specified timestamp", func() {
			now := time.Now().Round(time.Microsecond).UTC()
			BeforeEach(func() {
				_, err := Client().Wallets.V1.CreditWallet(TestContext(), operations.CreditWalletRequest{
					CreditWalletRequest: &shared.CreditWalletRequest{
						Amount: shared.Monetary{
							Amount: big.NewInt(1000),
							Asset:  "USD/2",
						},
						Sources:   []shared.Subject{},
						Metadata:  map[string]string{},
						Timestamp: &now,
					},
					ID: response.CreateWalletResponse.Data.ID,
				})
				Expect(err).To(Succeed())
			})
			It("should create the transaction at the specified date", func() {
				tx, err := Client().Ledger.V2.GetTransaction(TestContext(), operations.V2GetTransactionRequest{
					ID:     big.NewInt(0),
					Ledger: "wallets-002",
				})
				Expect(err).To(Succeed())
				Expect(tx.V2GetTransactionResponse.Data.Timestamp).To(Equal(now))
			})
		})
		Then("crediting it with invalid source", func() {
			It("should fail", func() {
				_, err := Client().Wallets.V1.CreditWallet(TestContext(), operations.CreditWalletRequest{
					CreditWalletRequest: &shared.CreditWalletRequest{
						Amount: shared.Monetary{
							Amount: big.NewInt(1000),
							Asset:  "USD/2",
						},
						Sources: []shared.Subject{shared.CreateSubjectAccount(shared.LedgerAccountSubject{
							Identifier: "@xxx",
						})},
						Metadata: map[string]string{},
					},
					ID: response.CreateWalletResponse.Data.ID,
				})
				Expect(err).NotTo(Succeed())
				sdkError := &sdkerrors.WalletsErrorResponse{}
				Expect(errors.As(err, &sdkError)).To(BeTrue())
				Expect(sdkError.ErrorCode).To(Equal(sdkerrors.SchemasWalletsErrorResponseErrorCodeValidation))
			})
		})
		Then("crediting it with negative amount", func() {
			It("should fail", func() {
				_, err := Client().Wallets.V1.CreditWallet(TestContext(), operations.CreditWalletRequest{
					CreditWalletRequest: &shared.CreditWalletRequest{
						Amount: shared.Monetary{
							Amount: big.NewInt(-1000),
							Asset:  "USD/2",
						},
						Sources:  []shared.Subject{},
						Metadata: map[string]string{},
					},
					ID: response.CreateWalletResponse.Data.ID,
				})
				Expect(err).NotTo(Succeed())
				sdkError := &sdkerrors.WalletsErrorResponse{}
				Expect(errors.As(err, &sdkError)).To(BeTrue())
				Expect(sdkError.ErrorCode).To(Equal(sdkerrors.SchemasWalletsErrorResponseErrorCodeValidation))
			})
		})
		Then("crediting it with invalid asset name", func() {
			It("should fail", func() {
				_, err := Client().Wallets.V1.CreditWallet(TestContext(), operations.CreditWalletRequest{
					CreditWalletRequest: &shared.CreditWalletRequest{
						Amount: shared.Monetary{
							Amount: big.NewInt(1000),
							Asset:  "test",
						},
						Sources:  []shared.Subject{},
						Metadata: map[string]string{},
					},
					ID: response.CreateWalletResponse.Data.ID,
				})
				Expect(err).NotTo(Succeed())
				sdkError := &sdkerrors.WalletsErrorResponse{}
				Expect(errors.As(err, &sdkError)).To(BeTrue())
				Expect(sdkError.ErrorCode).To(Equal(sdkerrors.SchemasWalletsErrorResponseErrorCodeValidation))
			})
		})
	})
})
