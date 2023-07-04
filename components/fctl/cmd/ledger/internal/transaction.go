package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/pkg/utils"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"golang.org/x/mod/semver"
)

// Copy from SDK, + fix int64 to big.Int
type Transaction struct {
	Metadata  map[string]interface{} `json:"metadata"`
	Postings  []*shared.Posting      `json:"postings"`
	Reference *string                `json:"reference,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Txid      int64                  `json:"txid"`
}

// Copy from SDK, + change metadata type to map[string]interface{}
// The is 2 types of metadata, one is map[string]interface{}, another is map[string]string
type ExpandedTransaction struct {
	Metadata          map[string]interface{}               `json:"metadata"`
	PostCommitVolumes map[string]map[string]*shared.Volume `json:"postCommitVolumes"`
	Postings          []*shared.Posting                    `json:"postings"`
	PreCommitVolumes  map[string]map[string]*shared.Volume `json:"preCommitVolumes"`
	Reference         *string                              `json:"reference,omitempty"`
	Timestamp         time.Time                            `json:"timestamp"`
	Txid              int64                                `json:"txid"`
}

// CreateTransactionResponse - OK
type CreateTransactionResponse struct {
	Data []Transaction `json:"data"`
}

type CreateTransactionWrapper struct {
	ContentType string
	// OK
	CreateTransactionResponse *CreateTransactionResponse
	// Error
	ErrorResponse *shared.ErrorResponse
	StatusCode    int
	RawResponse   *http.Response
}

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

// CreateTransaction - Create a new transaction to a ledger
func createTransactionV1(ctx context.Context, client *formance.Formance, baseURL string, request operations.CreateTransactionRequest) (*CreateTransactionWrapper, error) {

	// Dirty hack to get the http client from the sdk client struct
	field := reflect.ValueOf(client).Elem().FieldByName("_securityClient")
	httpClient := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface().(formance.HTTPClient)

	url, err := utils.GenerateURL(ctx, baseURL, "/api/ledger/{ledger}/transactions", request, nil)
	if err != nil {
		return nil, fmt.Errorf("error generating URL: %w", err)
	}

	bodyReader, reqContentType, err := utils.SerializeRequestBody(ctx, request, "PostTransaction", "json")
	if err != nil {
		return nil, fmt.Errorf("error serializing request body: %w", err)
	}
	if bodyReader == nil {
		return nil, fmt.Errorf("request body is required")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json;q=1, application/json;q=0")
	req.Header.Set("Content-Type", reqContentType)

	utils.PopulateHeaders(ctx, req, request)

	if err := utils.PopulateQueryParams(ctx, req, request, nil); err != nil {
		return nil, fmt.Errorf("error populating query params: %w", err)
	}

	httpRes, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	if httpRes == nil {
		return nil, fmt.Errorf("error sending request: no response")
	}
	defer httpRes.Body.Close()

	contentType := httpRes.Header.Get("Content-Type")

	res := &CreateTransactionWrapper{
		StatusCode:  httpRes.StatusCode,
		ContentType: contentType,
		RawResponse: httpRes,
	}
	switch {
	case httpRes.StatusCode == 200:
		switch {
		case utils.MatchContentType(contentType, `application/json`):
			var out *CreateTransactionResponse
			if err := utils.UnmarshalJsonFromResponseBody(httpRes.Body, &out); err != nil {
				return nil, err
			}

			res.CreateTransactionResponse = out
		}
	default:
		switch {
		case utils.MatchContentType(contentType, `application/json`):
			var out *shared.ErrorResponse
			if err := utils.UnmarshalJsonFromResponseBody(httpRes.Body, &out); err != nil {
				return nil, err
			}

			res.ErrorResponse = out
		}
	}

	return res, nil
}

func CreateTransaction(client *formance.Formance, ctx context.Context, request operations.CreateTransactionRequest) (*Transaction, error) {

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

	if semver.IsValid(version) && semver.Compare(version, "v2.0.0") < 0 {
		baseURL := strings.TrimSuffix(versionsResponse.RawResponse.Request.URL.String(), "/versions")

		v, err := createTransactionV1(ctx, client, baseURL, request)
		if err != nil {
			return nil, err
		}
		return &v.CreateTransactionResponse.Data[0], nil
	} else {
		response, err := client.Ledger.CreateTransaction(ctx, request)
		if err != nil {
			return nil, err
		}
		if response.StatusCode > 300 {
			return nil, fmt.Errorf("unexpected status code %d when creating transaction", response.StatusCode)
		}

		st := &response.CreateTransactionResponse.Data

		metadata := make(map[string]interface{})
		if st.Metadata != nil {
			for k, v := range st.Metadata {
				metadata[k] = v
			}
		}

		t := &Transaction{
			Txid:      st.Txid,
			Postings:  make([]*shared.Posting, len(st.Postings)),
			Reference: st.Reference,
			Timestamp: st.Timestamp,
			Metadata:  metadata,
		}

		return t, nil
	}
}
