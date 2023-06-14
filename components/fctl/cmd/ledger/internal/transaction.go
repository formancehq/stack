package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"golang.org/x/mod/semver"
)

func TransactionIDOrLastN(ctx context.Context, ledgerClient *formance.Formance, ledger, id string) (int64, error) {
	if strings.HasPrefix(id, "last") {
		id = strings.TrimPrefix(id, "last")
		sub := int64(0)
		if id != "" {
			var err error
			sub, err = strconv.ParseInt(id, 10, 64)
			if err != nil {
				return 0, err
			}
		}
		pageSize := int64(1)
		request := operations.ListTransactionsRequest{
			Ledger:   ledger,
			PageSize: &pageSize,
		}
		response, err := ledgerClient.Ledger.ListTransactions(ctx, request)
		if err != nil {
			return 0, err
		}

		if response.ErrorResponse != nil {
			return 0, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
		}

		if response.StatusCode >= 300 {
			return 0, fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}

		if len(response.TransactionsCursorResponse.Cursor.Data) == 0 {
			return 0, errors.New("no transaction found")
		}
		return response.TransactionsCursorResponse.Cursor.Data[0].Txid + sub, nil
	}

	return strconv.ParseInt(id, 10, 64)
}

func CreateTransaction(client *formance.Formance, ctx context.Context, ledger string, request operations.CreateTransactionRequest) (*shared.Transaction, error) {

	versionsResponse, err := client.GetVersions(ctx)
	if err != nil {
		return nil, err
	}
	if versionsResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when getting versions", versionsResponse.StatusCode)
	}

	version := collectionutils.Filter(versionsResponse.GetVersionsResponse.Versions, func(version shared.Version) bool {
		return version.Name == "ledger"
	})[0].Version

	response, err := client.Ledger.CreateTransaction(ctx, request)

	if semver.IsValid(version) && semver.Compare(version, "v2.0.0") < 0 {
		baseResponse := api.BaseResponse[[]shared.Transaction]{}
		if err := json.NewDecoder(response.RawResponse.Body).Decode(&baseResponse); err != nil {
			return nil, err
		}
		return &(*baseResponse.Data)[0], nil
	} else {
		if err != nil {
			return nil, err
		}
		if response.StatusCode > 300 {
			return nil, fmt.Errorf("unexpected status code %d when creating transaction", response.StatusCode)
		}
		return &response.CreateTransactionResponse.Data, nil
	}
}
