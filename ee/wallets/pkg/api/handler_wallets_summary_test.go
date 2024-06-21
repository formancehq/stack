package api

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/stack/libs/go-libs/metadata"
	wallet "github.com/formancehq/wallets/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestWalletSummary(t *testing.T) {
	t.Parallel()

	w := wallet.NewWallet(uuid.NewString(), "default", metadata.Metadata{})

	req := newRequest(t, http.MethodGet, "/wallets/"+w.ID+"/summary", nil)
	rec := httptest.NewRecorder()

	coupon1Balance := wallet.NewBalance("coupon1", ptr(time.Now().Add(-time.Minute).Round(time.Second).UTC()))
	coupon2Balance := wallet.NewBalance("coupon2", ptr(time.Now().Add(time.Minute).Round(time.Second).UTC()))
	hold1 := wallet.NewDebitHold(w.ID, wallet.NewLedgerAccountSubject("bank"), "USD", "", metadata.Metadata{})
	hold2 := wallet.NewDebitHold(w.ID, wallet.NewLedgerAccountSubject("bank"), "USD", "", metadata.Metadata{})

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithGetAccount(func(ctx context.Context, ledger, account string) (*wallet.AccountWithVolumesAndBalances, error) {
			require.Equal(t, testEnv.LedgerName(), ledger)
			switch account {
			case testEnv.Chart().GetMainBalanceAccount(w.ID):
				return &wallet.AccountWithVolumesAndBalances{
					Account: wallet.Account{
						Address:  account,
						Metadata: metadataWithExpectingTypesAfterUnmarshalling(w.LedgerMetadata()),
					},
					Balances: map[string]*big.Int{
						"USD": big.NewInt(100),
					},
				}, nil
			case testEnv.Chart().GetBalanceAccount(w.ID, coupon1Balance.Name):
				return &wallet.AccountWithVolumesAndBalances{
					Account: wallet.Account{
						Address:  account,
						Metadata: metadataWithExpectingTypesAfterUnmarshalling(coupon1Balance.LedgerMetadata(w.ID)),
					},
					Balances: map[string]*big.Int{
						"USD": big.NewInt(10),
					},
				}, nil
			case testEnv.Chart().GetBalanceAccount(w.ID, coupon2Balance.Name):
				return &wallet.AccountWithVolumesAndBalances{
					Account: wallet.Account{
						Address:  account,
						Metadata: metadataWithExpectingTypesAfterUnmarshalling(coupon2Balance.LedgerMetadata(w.ID)),
					},
					Balances: map[string]*big.Int{
						"USD": big.NewInt(20),
					},
				}, nil
			case testEnv.Chart().GetHoldAccount(hold1.ID):
				return &wallet.AccountWithVolumesAndBalances{
					Account: wallet.Account{
						Address:  account,
						Metadata: metadataWithExpectingTypesAfterUnmarshalling(hold1.LedgerMetadata(testEnv.Chart())),
					},
					Balances: map[string]*big.Int{
						"USD": big.NewInt(10),
					},
				}, nil
			case testEnv.Chart().GetHoldAccount(hold2.ID):
				return &wallet.AccountWithVolumesAndBalances{
					Account: wallet.Account{
						Address:  account,
						Metadata: metadataWithExpectingTypesAfterUnmarshalling(hold2.LedgerMetadata(testEnv.Chart())),
					},
					Balances: map[string]*big.Int{
						"USD": big.NewInt(20),
					},
				}, nil
			default:
				require.Fail(t, "unexpected account query")
			}
			panic("should not happen")
		}),
		WithListAccounts(func(ctx context.Context, ledger string, query wallet.ListAccountsQuery) (*wallet.AccountsCursorResponseCursor, error) {
			switch {
			case query.Metadata[wallet.MetadataKeyWalletID] == w.ID:
				return &wallet.AccountsCursorResponseCursor{
					Data: []wallet.AccountWithVolumesAndBalances{
						{
							Account: wallet.Account{
								Address:  testEnv.Chart().GetMainBalanceAccount(w.ID),
								Metadata: metadataWithExpectingTypesAfterUnmarshalling(w.LedgerMetadata()),
							},
						},
						{
							Account: wallet.Account{
								Address:  testEnv.Chart().GetBalanceAccount(w.ID, "coupon1"),
								Metadata: metadataWithExpectingTypesAfterUnmarshalling(coupon1Balance.LedgerMetadata(w.ID)),
							},
						},
						{
							Account: wallet.Account{
								Address:  testEnv.Chart().GetBalanceAccount(w.ID, "coupon2"),
								Metadata: metadataWithExpectingTypesAfterUnmarshalling(coupon2Balance.LedgerMetadata(w.ID)),
							},
						},
					},
				}, nil
			case query.Metadata[wallet.MetadataKeyHoldWalletID] == w.ID:
				return &wallet.AccountsCursorResponseCursor{
					Data: []wallet.AccountWithVolumesAndBalances{
						{
							Account: wallet.Account{
								Address:  testEnv.Chart().GetHoldAccount(hold1.ID),
								Metadata: metadataWithExpectingTypesAfterUnmarshalling(hold1.LedgerMetadata(testEnv.Chart())),
							},
						},
						{
							Account: wallet.Account{
								Address:  testEnv.Chart().GetHoldAccount(hold2.ID),
								Metadata: metadataWithExpectingTypesAfterUnmarshalling(hold2.LedgerMetadata(testEnv.Chart())),
							},
						},
					},
				}, nil
			default:
				require.Fail(t, "unexpected list accounts query")
			}
			panic("should not happen")
		}),
	)
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	summary := wallet.Summary{}
	readResponse(t, rec, &summary)

	require.Equal(t, wallet.Summary{
		Balances: []wallet.ExpandedBalance{
			{
				Balance: wallet.Balance{
					Name: "main",
				},
				Assets: map[string]*big.Int{
					"USD": big.NewInt(100),
				},
			},
			{
				Balance: coupon1Balance,
				Assets: map[string]*big.Int{
					"USD": big.NewInt(10),
				},
			},
			{
				Balance: coupon2Balance,
				Assets: map[string]*big.Int{
					"USD": big.NewInt(20),
				},
			},
		},
		AvailableFunds: map[string]*big.Int{
			"USD": big.NewInt(120),
		},
		ExpiredFunds: map[string]*big.Int{
			"USD": big.NewInt(10),
		},
		ExpirableFunds: map[string]*big.Int{
			"USD": big.NewInt(20),
		},
		HoldFunds: map[string]*big.Int{
			"USD": big.NewInt(30),
		},
	}, summary)
}
