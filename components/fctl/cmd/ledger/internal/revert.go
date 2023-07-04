package internal

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"unsafe"

	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/pkg/utils"
)

type RevertTransactionResponseData struct {
	Data Transaction `json:"data"`
}
type RevertTransactionResponse struct {
	ContentType string
	// Error
	ErrorResponse *shared.ErrorResponse
	// OK
	RevertTransactionResponse *RevertTransactionResponseData
	StatusCode                int
	RawResponse               *http.Response
}

func RevertTransaction(client *formance.Formance, ctx context.Context, baseURL string, request operations.RevertTransactionRequest) (*RevertTransactionResponse, error) {
	// Dirty hack to get the http client from the sdk client struct
	field := reflect.ValueOf(client).Elem().FieldByName("_securityClient")
	httpClient := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface().(formance.HTTPClient)

	url, err := utils.GenerateURL(ctx, baseURL, "/api/ledger/{ledger}/transactions/{txid}/revert", request, nil)
	if err != nil {
		return nil, fmt.Errorf("error generating URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json;q=1, application/json;q=0")

	httpRes, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	if httpRes == nil {
		return nil, fmt.Errorf("error sending request: no response")
	}
	defer httpRes.Body.Close()

	contentType := httpRes.Header.Get("Content-Type")

	res := &RevertTransactionResponse{
		StatusCode:  httpRes.StatusCode,
		ContentType: contentType,
		RawResponse: httpRes,
	}
	switch {
	case httpRes.StatusCode == 201:
		switch {
		case utils.MatchContentType(contentType, `application/json`):
			var out *RevertTransactionResponseData
			if err := utils.UnmarshalJsonFromResponseBody(httpRes.Body, &out); err != nil {
				return nil, err
			}

			res.RevertTransactionResponse = out
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
