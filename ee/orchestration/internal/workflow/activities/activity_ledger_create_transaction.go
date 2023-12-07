package activities

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"unsafe"

	formance "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/pkg/utils"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"golang.org/x/mod/semver"
)

// CreateTransactionResponse - OK
type CreateTransactionResponse struct {
	Data []shared.Transaction `json:"data"`
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

// CreateTransaction - Create a new transaction to a ledger
func createTransactionV1(ctx context.Context, client *formance.Formance, baseURL string, request CreateTransactionRequest) (*CreateTransactionWrapper, error) {

	// Dirty hack to get the http client from the sdk client struct
	field := reflect.ValueOf(client).Elem().FieldByName("_securityClient")
	httpClient := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface().(formance.HTTPClient)

	url, err := utils.GenerateURL(ctx, baseURL, "/api/ledger/{ledger}/transactions", request, nil)
	if err != nil {
		return nil, fmt.Errorf("error generating URL: %w", err)
	}

	bodyReader, reqContentType, err := utils.SerializeRequestBody(ctx, request, false, false, "PostTransaction", "json", `request:"mediaType=application/json"`)
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
	rawBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	httpRes.Body.Close()
	httpRes.Body = io.NopCloser(bytes.NewBuffer(rawBody))

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
			if err := utils.UnmarshalJsonFromResponseBody(bytes.NewBuffer(rawBody), &out, ""); err != nil {
				return nil, err
			}

			res.CreateTransactionResponse = out
		}
	default:
		switch {
		case utils.MatchContentType(contentType, `application/json`):
			var out *shared.ErrorResponse
			if err := utils.UnmarshalJsonFromResponseBody(bytes.NewBuffer(rawBody), &out, ""); err != nil {
				return nil, err
			}

			res.ErrorResponse = out
		}
	}

	return res, nil
}

type CreateTransactionRequest struct {
	Ledger string                 `pathParam:"style=simple,explode=false,name=ledger"`
	Data   shared.PostTransaction `request:"mediaType=application/json"`
}

func (a Activities) CreateTransaction(ctx context.Context, request CreateTransactionRequest) (*shared.CreateTransactionResponse, error) {
	versionsResponse, err := a.client.GetVersions(ctx)
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

		v, err := createTransactionV1(ctx, a.client, baseURL, request)
		if err != nil {
			return nil, err
		}

		switch v.StatusCode {
		case http.StatusOK:
			return &shared.CreateTransactionResponse{
				Data: v.CreateTransactionResponse.Data[0],
			}, nil
		default:
			if v.ErrorResponse != nil {
				return nil, temporal.NewApplicationError(
					v.ErrorResponse.ErrorMessage,
					string(v.ErrorResponse.ErrorCode),
					v.ErrorResponse.Details)
			}

			return nil, fmt.Errorf("unexpected status code: %d", v.StatusCode)
		}
	} else {
		response, err := a.client.Ledger.V2.
			CreateTransaction(
				ctx,
				operations.CreateTransactionRequest{
					PostTransaction: request.Data,
					Ledger:          request.Ledger,
				},
			)
		if err != nil {
			return nil, err
		}

		switch response.StatusCode {
		case http.StatusOK:
			return response.CreateTransactionResponse, nil
		default:
			if response.ErrorResponse != nil {
				return nil, temporal.NewApplicationError(
					response.ErrorResponse.ErrorMessage,
					string(response.ErrorResponse.ErrorCode),
					response.ErrorResponse.Details)
			}

			return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}
	}
}

var CreateTransactionActivity = Activities{}.CreateTransaction

func CreateTransaction(ctx workflow.Context, ledger string, request shared.PostTransaction) (*shared.Transaction, error) {
	tx := &shared.CreateTransactionResponse{}
	if err := executeActivity(ctx, CreateTransactionActivity, tx, CreateTransactionRequest{
		Ledger: ledger,
		Data:   request,
	}); err != nil {
		return nil, err
	}
	return &tx.Data, nil
}
