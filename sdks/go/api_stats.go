/*
Formance Stack API

Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions />

API version: develop
Contact: support@formance.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package formance

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type StatsApi interface {

	/*
			ReadStats Get Stats

			Get ledger stats (aggregate metrics on accounts and transactions)
		The stats for account


			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param ledger name of the ledger
			@return ApiReadStatsRequest
	*/
	ReadStats(ctx context.Context, ledger string) ApiReadStatsRequest

	// ReadStatsExecute executes the request
	//  @return StatsResponse
	ReadStatsExecute(r ApiReadStatsRequest) (*StatsResponse, *http.Response, error)
}

// StatsApiService StatsApi service
type StatsApiService service

type ApiReadStatsRequest struct {
	ctx        context.Context
	ApiService StatsApi
	ledger     string
}

func (r ApiReadStatsRequest) Execute() (*StatsResponse, *http.Response, error) {
	return r.ApiService.ReadStatsExecute(r)
}

/*
ReadStats Get Stats

Get ledger stats (aggregate metrics on accounts and transactions)
The stats for account

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param ledger name of the ledger
	@return ApiReadStatsRequest
*/
func (a *StatsApiService) ReadStats(ctx context.Context, ledger string) ApiReadStatsRequest {
	return ApiReadStatsRequest{
		ApiService: a,
		ctx:        ctx,
		ledger:     ledger,
	}
}

// Execute executes the request
//
//	@return StatsResponse
func (a *StatsApiService) ReadStatsExecute(r ApiReadStatsRequest) (*StatsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *StatsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StatsApiService.ReadStats")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/ledger/{ledger}/stats"
	localVarPath = strings.Replace(localVarPath, "{"+"ledger"+"}", url.PathEscape(parameterToString(r.ledger, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
