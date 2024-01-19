package api

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/pointer"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	wallet "github.com/formancehq/wallets/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTransactionsList(t *testing.T) {
	t.Parallel()

	w := wallet.NewWallet(uuid.NewString(), "default", metadata.Metadata{})

	var transactions []shared.Transaction
	for i := 0; i < 10; i++ {
		transactions = append(transactions, shared.Transaction{
			Postings: []shared.Posting{{
				Amount:      big.NewInt(100),
				Asset:       "USD/2",
				Destination: "bank",
				Source:      "world",
			}},
			Metadata: map[string]any{},
		})
	}
	const pageSize = 2
	numberOfPages := int64(len(transactions) / pageSize)

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithListTransactions(func(ctx context.Context, ledger string, query wallet.ListTransactionsQuery) (*shared.TransactionsCursorResponseCursor, error) {
			if query.Cursor != "" {
				page, err := strconv.ParseInt(query.Cursor, 10, 64)
				if err != nil {
					panic(err)
				}

				if page >= numberOfPages-1 {
					return &shared.TransactionsCursorResponseCursor{}, nil
				}
				hasMore := page < numberOfPages-1
				previous := fmt.Sprint(page - 1)
				next := fmt.Sprint(page + 1)

				return &shared.TransactionsCursorResponseCursor{
					PageSize: pageSize,
					HasMore:  hasMore,
					Previous: pointer.For(previous),
					Next:     pointer.For(next),
					Data:     transactions[page*pageSize : (page+1)*pageSize],
				}, nil
			}

			require.Equal(t, pageSize, query.Limit)
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetMainBalanceAccount(w.ID), query.Account)

			return &shared.TransactionsCursorResponseCursor{
				PageSize: pageSize,
				HasMore:  true,
				Next:     pointer.For("1"),
				Data:     transactions[:pageSize],
			}, nil
		}),
	)

	req := newRequest(t, http.MethodGet, fmt.Sprintf("/transactions?pageSize=%d&walletID=%s", pageSize, w.ID), nil)
	rec := httptest.NewRecorder()
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	cursor := &sharedapi.Cursor[shared.Transaction]{}
	readCursor(t, rec, cursor)
	require.Len(t, cursor.Data, pageSize)
	require.EqualValues(t, transactions[:pageSize], cursor.Data)

	req = newRequest(t, http.MethodGet, fmt.Sprintf("/transactions?cursor=%s", cursor.Next), nil)
	rec = httptest.NewRecorder()
	testEnv.Router().ServeHTTP(rec, req)
	cursor = &sharedapi.Cursor[shared.Transaction]{}
	readCursor(t, rec, cursor)
	require.Len(t, cursor.Data, pageSize)
	require.EqualValues(t, transactions[pageSize:pageSize*2], cursor.Data)
}
