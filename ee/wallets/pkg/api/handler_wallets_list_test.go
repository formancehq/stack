package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/pointer"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	wallet "github.com/formancehq/wallets/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestWalletsList(t *testing.T) {
	t.Parallel()

	var wallets []wallet.Wallet
	for i := 0; i < 10; i++ {
		wallets = append(wallets, wallet.NewWallet(uuid.NewString(), "default", metadata.Metadata{}))
	}
	const pageSize = 2
	numberOfPages := int64(len(wallets) / pageSize)

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithListAccounts(func(ctx context.Context, ledger string, query wallet.ListAccountsQuery) (*wallet.AccountsCursorResponseCursor, error) {
			if query.Cursor != "" {
				page, err := strconv.ParseInt(query.Cursor, 10, 64)
				if err != nil {
					panic(err)
				}

				if page >= numberOfPages-1 {
					return &wallet.AccountsCursorResponseCursor{}, nil
				}
				hasMore := page < numberOfPages-1
				previous := fmt.Sprint(page - 1)
				next := fmt.Sprint(page + 1)
				accounts := make([]wallet.Account, 0)
				for _, w := range wallets[page*pageSize : (page+1)*pageSize] {
					accounts = append(accounts, wallet.Account{
						Address:  testEnv.Chart().GetMainBalanceAccount(w.ID),
						Metadata: metadataWithExpectingTypesAfterUnmarshalling(w.LedgerMetadata()),
					})
				}

				return &wallet.AccountsCursorResponseCursor{
					PageSize: pageSize,
					HasMore:  hasMore,
					Previous: pointer.For(previous),
					Next:     pointer.For(next),
					Data:     accounts,
				}, nil
			}

			require.Equal(t, pageSize, query.Limit)
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, metadata.Metadata{
				wallet.MetadataKeyWalletSpecType: wallet.PrimaryWallet,
			}, query.Metadata)

			next := "1"
			accounts := make([]wallet.Account, 0)
			for _, w := range wallets[:pageSize] {
				accounts = append(accounts, wallet.Account{
					Address:  testEnv.Chart().GetMainBalanceAccount(w.ID),
					Metadata: metadataWithExpectingTypesAfterUnmarshalling(w.LedgerMetadata()),
				})
			}

			return &wallet.AccountsCursorResponseCursor{
				PageSize: pageSize,
				HasMore:  true,
				Next:     pointer.For(next),
				Data:     accounts,
			}, nil
		}),
	)

	req := newRequest(t, http.MethodGet, fmt.Sprintf("/wallets?pageSize=%d", pageSize), nil)
	rec := httptest.NewRecorder()
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	cursor := &sharedapi.Cursor[wallet.Wallet]{}
	readCursor(t, rec, cursor)
	require.Len(t, cursor.Data, pageSize)
	require.EqualValues(t, wallets[:pageSize], cursor.Data)

	req = newRequest(t, http.MethodGet, fmt.Sprintf("/wallets?cursor=%s", cursor.Next), nil)
	rec = httptest.NewRecorder()
	testEnv.Router().ServeHTTP(rec, req)
	cursor = &sharedapi.Cursor[wallet.Wallet]{}
	readCursor(t, rec, cursor)
	require.Len(t, cursor.Data, pageSize)
	require.EqualValues(t, cursor.Data, wallets[pageSize:pageSize*2])
}

func TestWalletsListByName(t *testing.T) {
	t.Parallel()

	var wallets []wallet.Wallet
	for i := 0; i < 10; i++ {
		wallets = append(wallets, wallet.NewWallet(uuid.NewString(), "default", metadata.Metadata{}))
	}

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithListAccounts(func(ctx context.Context, ledger string, query wallet.ListAccountsQuery) (*wallet.AccountsCursorResponseCursor, error) {
			require.Equal(t, defaultLimit, query.Limit)
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, metadata.Metadata{
				wallet.MetadataKeyWalletSpecType: wallet.PrimaryWallet,
				wallet.MetadataKeyWalletName:     wallets[1].Name,
			}, query.Metadata)

			return &wallet.AccountsCursorResponseCursor{
				PageSize: defaultLimit,
				HasMore:  false,
				Data: []wallet.Account{{
					Address:  testEnv.Chart().GetMainBalanceAccount(wallets[1].ID),
					Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallets[1].LedgerMetadata()),
				}},
			}, nil
		}),
	)

	req := newRequest(t, http.MethodGet, fmt.Sprintf("/wallets?name=%s", wallets[1].Name), nil)
	rec := httptest.NewRecorder()
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	cursor := &sharedapi.Cursor[wallet.Wallet]{}
	readCursor(t, rec, cursor)
	require.Len(t, cursor.Data, 1)
	require.EqualValues(t, wallets[1], cursor.Data[0])
}

func TestWalletsListFilterMetadata(t *testing.T) {
	t.Parallel()

	var wallets []wallet.Wallet
	for i := 0; i < 10; i++ {
		wallets = append(wallets, wallet.NewWallet(uuid.NewString(), "default", metadata.Metadata{
			"wallet": fmt.Sprint(i),
		}))
	}

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithListAccounts(func(ctx context.Context, ledger string, query wallet.ListAccountsQuery) (*wallet.AccountsCursorResponseCursor, error) {
			require.Equal(t, defaultLimit, query.Limit)
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, metadata.Metadata{
				wallet.MetadataKeyWalletSpecType: wallet.PrimaryWallet,
			}.Merge(wallet.EncodeCustomMetadata(map[string]string{
				"wallet": "2",
			})), query.Metadata)

			return &wallet.AccountsCursorResponseCursor{
				PageSize: defaultLimit,
				Data: []wallet.Account{{
					Address:  testEnv.Chart().GetMainBalanceAccount(wallets[2].ID),
					Metadata: metadataWithExpectingTypesAfterUnmarshalling(wallets[2].LedgerMetadata()),
				}},
			}, nil
		}),
	)

	req := newRequest(t, http.MethodGet, "/wallets?metadata[wallet]=2", nil)
	rec := httptest.NewRecorder()
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	cursor := &sharedapi.Cursor[wallet.Wallet]{}
	readCursor(t, rec, cursor)
	require.Len(t, cursor.Data, 1)
	require.EqualValues(t, wallets[2], cursor.Data[0])
}
