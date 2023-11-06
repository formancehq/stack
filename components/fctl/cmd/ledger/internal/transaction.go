package internal

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	formance "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
)

func TransactionIDOrLastN(ctx context.Context, ledgerClient *formance.Formance, ledger, id string) (*big.Int, error) {
	if strings.HasPrefix(id, "last") {
		id = strings.TrimPrefix(id, "last")
		sub := int64(0)
		if id != "" {
			var err error
			sub, err = strconv.ParseInt(id, 10, 64)
			if err != nil {
				return nil, err
			}
		}
		pageSize := int64(1)
		request := operations.ListTransactionsRequest{
			Ledger:   ledger,
			PageSize: &pageSize,
		}
		response, err := ledgerClient.Ledger.ListTransactions(ctx, request)
		if err != nil {
			return nil, err
		}

		if response.ErrorResponse != nil {
			return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
		}

		if response.StatusCode >= 300 {
			return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}

		if len(response.TransactionsCursorResponse.Cursor.Data) == 0 {
			return nil, errors.New("no transaction found")
		}
		return response.TransactionsCursorResponse.Cursor.Data[0].ID.Sub(
			response.TransactionsCursorResponse.Cursor.Data[0].ID,
			big.NewInt(sub),
		), nil
	}

	v, ok := big.NewInt(0).SetString(id, 10)
	if !ok {
		return nil, fmt.Errorf("invalid bigint: %s", id)
	}

	return v, nil
}
