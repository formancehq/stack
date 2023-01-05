package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	sdk "github.com/formancehq/formance-sdk-go"
	sharedapi "github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/metadata"
	"github.com/formancehq/wallets/pkg/core"
	"github.com/formancehq/wallets/pkg/wallet"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTransactionsList(t *testing.T) {
	t.Parallel()

	w := core.NewWallet(uuid.NewString(), metadata.Metadata{})

	var transactions []sdk.Transaction
	for i := 0; i < 10; i++ {
		transactions = append(transactions, sdk.Transaction{
			Postings: []sdk.Posting{{
				Amount:      100,
				Asset:       "USD/2",
				Destination: "bank",
				Source:      "world",
			}},
		})
	}
	const pageSize = 2
	numberOfPages := int64(len(transactions) / pageSize)

	var testEnv *testEnv
	testEnv = newTestEnv(
		WithListTransactions(func(ctx context.Context, ledger string, query wallet.ListTransactionsQuery) (*sdk.ListTransactions200ResponseCursor, error) {
			if query.PaginationToken != "" {
				page, err := strconv.ParseInt(query.PaginationToken, 10, 64)
				if err != nil {
					panic(err)
				}

				if page >= numberOfPages-1 {
					return &sdk.ListTransactions200ResponseCursor{}, nil
				}
				hasMore := page < numberOfPages-1
				previous := fmt.Sprint(page - 1)
				next := fmt.Sprint(page + 1)

				return &sdk.ListTransactions200ResponseCursor{
					PageSize: pageSize,
					HasMore:  &hasMore,
					Previous: &previous,
					Next:     &next,
					Data:     transactions[page*pageSize : (page+1)*pageSize],
				}, nil
			}

			require.Equal(t, pageSize, query.Limit)
			require.Equal(t, testEnv.LedgerName(), ledger)
			require.Equal(t, testEnv.Chart().GetMainAccount(w.ID), query.Account)

			hasMore := true
			next := "1"

			return &sdk.ListTransactions200ResponseCursor{
				PageSize: pageSize,
				HasMore:  &hasMore,
				Next:     &next,
				Data:     transactions[:pageSize],
			}, nil
		}),
	)

	req := newRequest(t, http.MethodGet, fmt.Sprintf("/transactions?limit=%d&walletID=%s", pageSize, w.ID), nil)
	rec := httptest.NewRecorder()
	testEnv.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	cursor := &sharedapi.Cursor[sdk.Transaction]{}
	readCursor(t, rec, cursor)
	require.Len(t, cursor.Data, pageSize)
	require.EqualValues(t, cursor.Data, transactions[:pageSize])

	req = newRequest(t, http.MethodGet, fmt.Sprintf("/transactions?cursor=%s", cursor.Next), nil)
	rec = httptest.NewRecorder()
	testEnv.Router().ServeHTTP(rec, req)
	cursor = &sharedapi.Cursor[sdk.Transaction]{}
	readCursor(t, rec, cursor)
	require.Len(t, cursor.Data, pageSize)
	require.EqualValues(t, cursor.Data, transactions[pageSize:pageSize*2])
}
