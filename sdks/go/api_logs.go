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
	"time"
)


type LogsApi interface {

	/*
	ListLogs List the logs from a ledger

	List the logs from a ledger, sorted by ID in descending order.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param ledger Name of the ledger.
	@return ApiListLogsRequest
	*/
	ListLogs(ctx context.Context, ledger string) ApiListLogsRequest

	// ListLogsExecute executes the request
	//  @return LogsCursorResponse
	ListLogsExecute(r ApiListLogsRequest) (*LogsCursorResponse, *http.Response, error)
}

// LogsApiService LogsApi service
type LogsApiService service

type ApiListLogsRequest struct {
	ctx context.Context
	ApiService LogsApi
	ledger string
	pageSize *int64
	after *string
	startTime *time.Time
	endTime *time.Time
	cursor *string
}

// The maximum number of results to return per page. 
func (r ApiListLogsRequest) PageSize(pageSize int64) ApiListLogsRequest {
	r.pageSize = &pageSize
	return r
}

// Pagination cursor, will return the logs after a given ID. (in descending order).
func (r ApiListLogsRequest) After(after string) ApiListLogsRequest {
	r.after = &after
	return r
}

// Filter transactions that occurred after this timestamp. The format is RFC3339 and is inclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; includes the first second of 4th minute). 
func (r ApiListLogsRequest) StartTime(startTime time.Time) ApiListLogsRequest {
	r.startTime = &startTime
	return r
}

// Filter transactions that occurred before this timestamp. The format is RFC3339 and is exclusive (for example, \&quot;2023-01-02T15:04:01Z\&quot; excludes the first second of 4th minute). 
func (r ApiListLogsRequest) EndTime(endTime time.Time) ApiListLogsRequest {
	r.endTime = &endTime
	return r
}

// Parameter used in pagination requests. Maximum page size is set to 15. Set to the value of next for the next page of results. Set to the value of previous for the previous page of results. No other parameters can be set when this parameter is set. 
func (r ApiListLogsRequest) Cursor(cursor string) ApiListLogsRequest {
	r.cursor = &cursor
	return r
}

func (r ApiListLogsRequest) Execute() (*LogsCursorResponse, *http.Response, error) {
	return r.ApiService.ListLogsExecute(r)
}

/*
ListLogs List the logs from a ledger

List the logs from a ledger, sorted by ID in descending order.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param ledger Name of the ledger.
 @return ApiListLogsRequest
*/
func (a *LogsApiService) ListLogs(ctx context.Context, ledger string) ApiListLogsRequest {
	return ApiListLogsRequest{
		ApiService: a,
		ctx: ctx,
		ledger: ledger,
	}
}

// Execute executes the request
//  @return LogsCursorResponse
func (a *LogsApiService) ListLogsExecute(r ApiListLogsRequest) (*LogsCursorResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *LogsCursorResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "LogsApiService.ListLogs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/ledger/{ledger}/log"
	localVarPath = strings.Replace(localVarPath, "{"+"ledger"+"}", url.PathEscape(parameterToString(r.ledger, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.pageSize != nil {
		localVarQueryParams.Add("pageSize", parameterToString(*r.pageSize, ""))
	}
	if r.after != nil {
		localVarQueryParams.Add("after", parameterToString(*r.after, ""))
	}
	if r.startTime != nil {
		localVarQueryParams.Add("startTime", parameterToString(*r.startTime, ""))
	}
	if r.endTime != nil {
		localVarQueryParams.Add("endTime", parameterToString(*r.endTime, ""))
	}
	if r.cursor != nil {
		localVarQueryParams.Add("cursor", parameterToString(*r.cursor, ""))
	}
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
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
            		newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
            		newErr.model = v
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
