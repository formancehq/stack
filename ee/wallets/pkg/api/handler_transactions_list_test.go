package api

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	wallet "github.com/formancehq/wallets/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTransactionsList(t *testing.T) {
	t.Parallel()

	w := wallet.NewWallet(uuid.NewString(), "default", metadata.Metadata{})

	var transactions []wallet.ExpandedTransaction
	for i := 0; i < 10; i++ {
		transactions = append(transactions, wallet.ExpandedTransaction{
			Postings: []wallet.Posting{{
				Amount:      big.NewInt(100),
				Asset:       "USD/2",
				Destination: "bank",
				Source:      "world",
			}},
			Metadata: metadata.Metadata{},
		})
	}
	const pageSize = 2
	numberOfPages := int64(len(transactions) / pageSize)

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithListTransactions(func(ctx context.Context, ledger string, query wallet.ListTransactionsQuery) (*wallet.TransactionsCursorResponseCursor, error) {
			if query.Cursor != "" {
				page, err := strconv.ParseInt(query.Cursor, 10, 64)
				if err != nil {
					panic(err)
				}

				if page >= numberOfPages-1 {
					return &wallet.TransactionsCursorResponseCursor{}, nil
				}
				hasMore := page < numberOfPages-1
				previous := fmt.Sprint(page - 1)
				next := fmt.Sprint(page + 1)

				return &wallet.TransactionsCursorResponseCursor{
					PageSize: pageSize,
					HasMore:  hasMore,
					Previous: previous,
					Next:     next,
					Data:     transactions[page*pageSize : (page+1)*pageSize],
				}, nil
			}

			require.Equal(t, pageSize, query.Limit)
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetMainBalanceAccount(w.ID), query.Account)

			hasMore := true
			next := "1"

			return &wallet.TransactionsCursorResponseCursor{
				PageSize: pageSize,
				HasMore:  hasMore,
				Next:     next,
				Data:     transactions[:pageSize],
			}, nil
		}),
	)

	req := newRequest(t, http.MethodGet, fmt.Sprintf("/transactions?pageSize=%d&walletID=%s", pageSize, w.ID), nil)
	rec := httptest.NewRecorder()
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	cursor := &sharedapi.Cursor[wallet.ExpandedTransaction]{}
	readCursor(t, rec, cursor)
	require.Len(t, cursor.Data, pageSize)
	require.EqualValues(t, transactions[:pageSize], cursor.Data)

	req = newRequest(t, http.MethodGet, fmt.Sprintf("/transactions?cursor=%s", cursor.Next), nil)
	rec = httptest.NewRecorder()
	testEnv.Router().ServeHTTP(rec, req)
	cursor = &sharedapi.Cursor[wallet.ExpandedTransaction]{}
	readCursor(t, rec, cursor)
	require.Len(t, cursor.Data, pageSize)
	require.EqualValues(t, transactions[pageSize:pageSize*2], cursor.Data)
}
