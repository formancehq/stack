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
			createWalletResponse *operations.CreateWalletResponse
			err                  error
		)
		BeforeEach(func() {
			createWalletResponse, err = Client().Wallets.V1.CreateWallet(
				TestContext(),
				operations.CreateWalletRequest{
					CreateWalletRequest: &shared.CreateWalletRequest{
						Name:     uuid.NewString(),
						Metadata: map[string]string{},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(createWalletResponse.StatusCode).To(Equal(201))
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
					ID: createWalletResponse.CreateWalletResponse.Data.ID,
				})
				Expect(err).To(Succeed())
			})
			Then("debiting it", func() {
				BeforeEach(func() {
					_, err := Client().Wallets.V1.DebitWallet(TestContext(), operations.DebitWalletRequest{
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
			Then("debiting it using timestamp", func() {
				now := time.Now().Round(time.Microsecond).UTC()
				BeforeEach(func() {
					_, err := Client().Wallets.V1.DebitWallet(TestContext(), operations.DebitWalletRequest{
						DebitWalletRequest: &shared.DebitWalletRequest{
							Amount: shared.Monetary{
								Amount: big.NewInt(100),
								Asset:  "USD/2",
							},
							Metadata:  map[string]string{},
							Timestamp: &now,
						},
						ID: createWalletResponse.CreateWalletResponse.Data.ID,
					})
					Expect(err).To(Succeed())
				})
				It("should create the transaction at the specified timestamp", func() {
					tx, err := Client().Ledger.V2.GetTransaction(TestContext(), operations.V2GetTransactionRequest{
						ID:     big.NewInt(1),
						Ledger: "wallets-002",
					})
					Expect(err).To(Succeed())
					Expect(tx.V2GetTransactionResponse.Data.Timestamp).To(Equal(now))
				})
			})
			Then("debiting it using a hold", func() {
				var (
					debitWalletResponse *operations.DebitWalletResponse
					ts                  *time.Time
				)
				JustBeforeEach(func() {
					debitWalletResponse, err = Client().Wallets.V1.DebitWallet(TestContext(), operations.DebitWalletRequest{
						DebitWalletRequest: &shared.DebitWalletRequest{
							Amount: shared.Monetary{
								Amount: big.NewInt(100),
								Asset:  "USD/2",
							},
							Pending:   pointer.For(true),
							Metadata:  map[string]string{},
							Timestamp: ts,
						},
						ID: createWalletResponse.CreateWalletResponse.Data.ID,
					})
					Expect(err).To(Succeed())
				})
				It("should be ok", func() {})
				Then("confirm the hold", func() {
					JustBeforeEach(func() {
						_, err := Client().Wallets.V1.ConfirmHold(TestContext(), operations.ConfirmHoldRequest{
							HoldID: debitWalletResponse.DebitWalletResponse.Data.ID,
						})
						Expect(err).To(Succeed())
					})
					It("should be ok", func() {})
				})
				Then("void the hold", func() {
					JustBeforeEach(func() {
						_, err := Client().Wallets.V1.VoidHold(TestContext(), operations.VoidHoldRequest{
							HoldID: debitWalletResponse.DebitWalletResponse.Data.ID,
						})
						Expect(err).To(Succeed())
					})
					It("should be ok", func() {})
				})
			})
			Then("debiting it using invalid destination", func() {
				It("should fail", func() {
					_, err := Client().Wallets.V1.DebitWallet(TestContext(), operations.DebitWalletRequest{
						DebitWalletRequest: &shared.DebitWalletRequest{
							Amount: shared.Monetary{
								Amount: big.NewInt(100),
								Asset:  "USD/2",
							},
							Metadata: map[string]string{},
							Destination: pointer.For(shared.CreateSubjectAccount(shared.LedgerAccountSubject{
								Identifier: "@xxx",
							})),
						},
						ID: createWalletResponse.CreateWalletResponse.Data.ID,
					})
					Expect(err).NotTo(Succeed())
					sdkError := &sdkerrors.WalletsErrorResponse{}
					Expect(errors.As(err, &sdkError)).To(BeTrue())
					Expect(sdkError.ErrorCode).To(Equal(sdkerrors.SchemasWalletsErrorResponseErrorCodeValidation))
				})
			})
			Then("debiting it using negative amount", func() {
				It("should fail", func() {
					_, err := Client().Wallets.V1.DebitWallet(TestContext(), operations.DebitWalletRequest{
						DebitWalletRequest: &shared.DebitWalletRequest{
							Amount: shared.Monetary{
								Amount: big.NewInt(-100),
								Asset:  "USD/2",
							},
							Metadata: map[string]string{},
						},
						ID: createWalletResponse.CreateWalletResponse.Data.ID,
					})
					Expect(err).NotTo(Succeed())
					sdkError := &sdkerrors.WalletsErrorResponse{}
					Expect(errors.As(err, &sdkError)).To(BeTrue())
					Expect(sdkError.ErrorCode).To(Equal(sdkerrors.SchemasWalletsErrorResponseErrorCodeValidation))
				})
			})
			Then("debiting it using invalid asset", func() {
				It("should fail", func() {
					_, err := Client().Wallets.V1.DebitWallet(TestContext(), operations.DebitWalletRequest{
						DebitWalletRequest: &shared.DebitWalletRequest{
							Amount: shared.Monetary{
								Amount: big.NewInt(100),
								Asset:  "test",
							},
							Metadata: map[string]string{},
						},
						ID: createWalletResponse.CreateWalletResponse.Data.ID,
					})
					Expect(err).NotTo(Succeed())
					sdkError := &sdkerrors.WalletsErrorResponse{}
					Expect(errors.As(err, &sdkError)).To(BeTrue())
					Expect(sdkError.ErrorCode).To(Equal(sdkerrors.SchemasWalletsErrorResponseErrorCodeValidation))
				})
			})
		})
	})
})
