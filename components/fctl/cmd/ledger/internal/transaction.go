package internal

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"golang.org/x/mod/semver"
)

func TransactionIDOrLastN(ctx context.Context, ledgerClient *formance.APIClient, ledger, id string) (int64, error) {
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
		response, _, err := ledgerClient.TransactionsApi.
			ListTransactions(ctx, ledger).
			PageSize(1).
			Execute()
		if err != nil {
			return 0, err
		}
		if len(response.Cursor.Data) == 0 {
			return 0, errors.New("no transaction found")
		}
		return response.Cursor.Data[0].Txid + sub, nil
	}

	return strconv.ParseInt(id, 10, 64)
}

func CreateTransaction(client *formance.APIClient, ctx context.Context, ledger string, transaction formance.PostTransaction) (*formance.Transaction, error) {

	versions, _, err := client.DefaultApi.GetVersions(ctx).Execute()
	if err != nil {
		return nil, err
	}

	version := collectionutils.Filter(versions.Versions, func(version formance.Version) bool {
		return version.Name == "ledger"
	})[0].Version

	response, httpResponse, err := client.TransactionsApi.CreateTransaction(ctx, ledger).
		PostTransaction(transaction).Execute()

	if semver.IsValid(version) && semver.Compare(version, "v2.0.0") < 0 {
		baseResponse := api.BaseResponse[[]formance.Transaction]{}
		if err := json.NewDecoder(httpResponse.Body).Decode(&baseResponse); err != nil {
			return nil, err
		}
		return &(*baseResponse.Data)[0], nil
	} else {
		if err != nil {
			return nil, err
		}
		return &response.Data, nil
	}
}
