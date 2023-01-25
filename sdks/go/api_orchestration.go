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


type OrchestrationApi interface {

	/*
	CreateWorkflow Create workflow

	Create a workflow

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiCreateWorkflowRequest
	*/
	CreateWorkflow(ctx context.Context) ApiCreateWorkflowRequest

	// CreateWorkflowExecute executes the request
	//  @return CreateWorkflowResponse
	CreateWorkflowExecute(r ApiCreateWorkflowRequest) (*CreateWorkflowResponse, *http.Response, error)

	/*
	GetInstance Get a workflow instance by id

	Get a workflow instance by id

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param instanceID The instance id
	@return ApiGetInstanceRequest
	*/
	GetInstance(ctx context.Context, instanceID string) ApiGetInstanceRequest

	// GetInstanceExecute executes the request
	//  @return GetWorkflowInstanceResponse
	GetInstanceExecute(r ApiGetInstanceRequest) (*GetWorkflowInstanceResponse, *http.Response, error)

	/*
	GetInstanceHistory Get a workflow instance history by id

	Get a workflow instance history by id

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param instanceID The instance id
	@return ApiGetInstanceHistoryRequest
	*/
	GetInstanceHistory(ctx context.Context, instanceID string) ApiGetInstanceHistoryRequest

	// GetInstanceHistoryExecute executes the request
	//  @return GetWorkflowInstanceHistoryResponse
	GetInstanceHistoryExecute(r ApiGetInstanceHistoryRequest) (*GetWorkflowInstanceHistoryResponse, *http.Response, error)

	/*
	GetInstanceStageHistory Get a workflow instance stage history

	Get a workflow instance stage history

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param instanceID The instance id
	@param number The stage number
	@return ApiGetInstanceStageHistoryRequest
	*/
	GetInstanceStageHistory(ctx context.Context, instanceID string, number int32) ApiGetInstanceStageHistoryRequest

	// GetInstanceStageHistoryExecute executes the request
	//  @return GetWorkflowInstanceHistoryStageResponse
	GetInstanceStageHistoryExecute(r ApiGetInstanceStageHistoryRequest) (*GetWorkflowInstanceHistoryStageResponse, *http.Response, error)

	/*
	GetWorkflow Get a flow by id

	Get a flow by id

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param flowId The flow id
	@return ApiGetWorkflowRequest
	*/
	GetWorkflow(ctx context.Context, flowId string) ApiGetWorkflowRequest

	// GetWorkflowExecute executes the request
	//  @return GetWorkflowResponse
	GetWorkflowExecute(r ApiGetWorkflowRequest) (*GetWorkflowResponse, *http.Response, error)

	/*
	ListInstances List instances of a workflow

	List instances of a workflow

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiListInstancesRequest
	*/
	ListInstances(ctx context.Context) ApiListInstancesRequest

	// ListInstancesExecute executes the request
	//  @return ListRunsResponse
	ListInstancesExecute(r ApiListInstancesRequest) (*ListRunsResponse, *http.Response, error)

	/*
	ListWorkflows List registered workflows

	List registered workflows

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiListWorkflowsRequest
	*/
	ListWorkflows(ctx context.Context) ApiListWorkflowsRequest

	// ListWorkflowsExecute executes the request
	//  @return ListWorkflowsResponse
	ListWorkflowsExecute(r ApiListWorkflowsRequest) (*ListWorkflowsResponse, *http.Response, error)

	/*
	OrchestrationgetServerInfo Get server info

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiOrchestrationgetServerInfoRequest
	*/
	OrchestrationgetServerInfo(ctx context.Context) ApiOrchestrationgetServerInfoRequest

	// OrchestrationgetServerInfoExecute executes the request
	//  @return ServerInfo
	OrchestrationgetServerInfoExecute(r ApiOrchestrationgetServerInfoRequest) (*ServerInfo, *http.Response, error)

	/*
	RunWorkflow Run workflow

	Run workflow

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param workflowID The flow id
	@return ApiRunWorkflowRequest
	*/
	RunWorkflow(ctx context.Context, workflowID string) ApiRunWorkflowRequest

	// RunWorkflowExecute executes the request
	//  @return RunWorkflowResponse
	RunWorkflowExecute(r ApiRunWorkflowRequest) (*RunWorkflowResponse, *http.Response, error)
}

// OrchestrationApiService OrchestrationApi service
type OrchestrationApiService service

type ApiCreateWorkflowRequest struct {
	ctx context.Context
	ApiService OrchestrationApi
	body *WorkflowConfig
}

func (r ApiCreateWorkflowRequest) Body(body WorkflowConfig) ApiCreateWorkflowRequest {
	r.body = &body
	return r
}

func (r ApiCreateWorkflowRequest) Execute() (*CreateWorkflowResponse, *http.Response, error) {
	return r.ApiService.CreateWorkflowExecute(r)
}

/*
CreateWorkflow Create workflow

Create a workflow

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiCreateWorkflowRequest
*/
func (a *OrchestrationApiService) CreateWorkflow(ctx context.Context) ApiCreateWorkflowRequest {
	return ApiCreateWorkflowRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return CreateWorkflowResponse
func (a *OrchestrationApiService) CreateWorkflowExecute(r ApiCreateWorkflowRequest) (*CreateWorkflowResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *CreateWorkflowResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrchestrationApiService.CreateWorkflow")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/orchestration/workflows"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

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
	// body params
	localVarPostBody = r.body
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
			var v Error
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

type ApiGetInstanceRequest struct {
	ctx context.Context
	ApiService OrchestrationApi
	instanceID string
}

func (r ApiGetInstanceRequest) Execute() (*GetWorkflowInstanceResponse, *http.Response, error) {
	return r.ApiService.GetInstanceExecute(r)
}

/*
GetInstance Get a workflow instance by id

Get a workflow instance by id

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param instanceID The instance id
 @return ApiGetInstanceRequest
*/
func (a *OrchestrationApiService) GetInstance(ctx context.Context, instanceID string) ApiGetInstanceRequest {
	return ApiGetInstanceRequest{
		ApiService: a,
		ctx: ctx,
		instanceID: instanceID,
	}
}

// Execute executes the request
//  @return GetWorkflowInstanceResponse
func (a *OrchestrationApiService) GetInstanceExecute(r ApiGetInstanceRequest) (*GetWorkflowInstanceResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *GetWorkflowInstanceResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrchestrationApiService.GetInstance")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/orchestration/instances/{instanceID}"
	localVarPath = strings.Replace(localVarPath, "{"+"instanceID"+"}", url.PathEscape(parameterValueToString(r.instanceID, "instanceID")), -1)

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
			var v Error
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

type ApiGetInstanceHistoryRequest struct {
	ctx context.Context
	ApiService OrchestrationApi
	instanceID string
}

func (r ApiGetInstanceHistoryRequest) Execute() (*GetWorkflowInstanceHistoryResponse, *http.Response, error) {
	return r.ApiService.GetInstanceHistoryExecute(r)
}

/*
GetInstanceHistory Get a workflow instance history by id

Get a workflow instance history by id

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param instanceID The instance id
 @return ApiGetInstanceHistoryRequest
*/
func (a *OrchestrationApiService) GetInstanceHistory(ctx context.Context, instanceID string) ApiGetInstanceHistoryRequest {
	return ApiGetInstanceHistoryRequest{
		ApiService: a,
		ctx: ctx,
		instanceID: instanceID,
	}
}

// Execute executes the request
//  @return GetWorkflowInstanceHistoryResponse
func (a *OrchestrationApiService) GetInstanceHistoryExecute(r ApiGetInstanceHistoryRequest) (*GetWorkflowInstanceHistoryResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *GetWorkflowInstanceHistoryResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrchestrationApiService.GetInstanceHistory")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/orchestration/instances/{instanceID}/history"
	localVarPath = strings.Replace(localVarPath, "{"+"instanceID"+"}", url.PathEscape(parameterValueToString(r.instanceID, "instanceID")), -1)

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
			var v Error
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

type ApiGetInstanceStageHistoryRequest struct {
	ctx context.Context
	ApiService OrchestrationApi
	instanceID string
	number int32
}

func (r ApiGetInstanceStageHistoryRequest) Execute() (*GetWorkflowInstanceHistoryStageResponse, *http.Response, error) {
	return r.ApiService.GetInstanceStageHistoryExecute(r)
}

/*
GetInstanceStageHistory Get a workflow instance stage history

Get a workflow instance stage history

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param instanceID The instance id
 @param number The stage number
 @return ApiGetInstanceStageHistoryRequest
*/
func (a *OrchestrationApiService) GetInstanceStageHistory(ctx context.Context, instanceID string, number int32) ApiGetInstanceStageHistoryRequest {
	return ApiGetInstanceStageHistoryRequest{
		ApiService: a,
		ctx: ctx,
		instanceID: instanceID,
		number: number,
	}
}

// Execute executes the request
//  @return GetWorkflowInstanceHistoryStageResponse
func (a *OrchestrationApiService) GetInstanceStageHistoryExecute(r ApiGetInstanceStageHistoryRequest) (*GetWorkflowInstanceHistoryStageResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *GetWorkflowInstanceHistoryStageResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrchestrationApiService.GetInstanceStageHistory")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/orchestration/instances/{instanceID}/stages/{number}/history"
	localVarPath = strings.Replace(localVarPath, "{"+"instanceID"+"}", url.PathEscape(parameterValueToString(r.instanceID, "instanceID")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"number"+"}", url.PathEscape(parameterValueToString(r.number, "number")), -1)

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
			var v Error
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

type ApiGetWorkflowRequest struct {
	ctx context.Context
	ApiService OrchestrationApi
	flowId string
}

func (r ApiGetWorkflowRequest) Execute() (*GetWorkflowResponse, *http.Response, error) {
	return r.ApiService.GetWorkflowExecute(r)
}

/*
GetWorkflow Get a flow by id

Get a flow by id

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param flowId The flow id
 @return ApiGetWorkflowRequest
*/
func (a *OrchestrationApiService) GetWorkflow(ctx context.Context, flowId string) ApiGetWorkflowRequest {
	return ApiGetWorkflowRequest{
		ApiService: a,
		ctx: ctx,
		flowId: flowId,
	}
}

// Execute executes the request
//  @return GetWorkflowResponse
func (a *OrchestrationApiService) GetWorkflowExecute(r ApiGetWorkflowRequest) (*GetWorkflowResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *GetWorkflowResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrchestrationApiService.GetWorkflow")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/orchestration/workflows/{flowId}"
	localVarPath = strings.Replace(localVarPath, "{"+"flowId"+"}", url.PathEscape(parameterValueToString(r.flowId, "flowId")), -1)

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
			var v Error
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

type ApiListInstancesRequest struct {
	ctx context.Context
	ApiService OrchestrationApi
	workflowID *string
}

// A workflow id
func (r ApiListInstancesRequest) WorkflowID(workflowID string) ApiListInstancesRequest {
	r.workflowID = &workflowID
	return r
}

func (r ApiListInstancesRequest) Execute() (*ListRunsResponse, *http.Response, error) {
	return r.ApiService.ListInstancesExecute(r)
}

/*
ListInstances List instances of a workflow

List instances of a workflow

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiListInstancesRequest
*/
func (a *OrchestrationApiService) ListInstances(ctx context.Context) ApiListInstancesRequest {
	return ApiListInstancesRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return ListRunsResponse
func (a *OrchestrationApiService) ListInstancesExecute(r ApiListInstancesRequest) (*ListRunsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *ListRunsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrchestrationApiService.ListInstances")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/orchestration/instances"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.workflowID == nil {
		return localVarReturnValue, nil, reportError("workflowID is required and must be specified")
	}

	parameterAddToQuery(localVarQueryParams, "workflowID", r.workflowID, "")
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
			var v Error
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

type ApiListWorkflowsRequest struct {
	ctx context.Context
	ApiService OrchestrationApi
}

func (r ApiListWorkflowsRequest) Execute() (*ListWorkflowsResponse, *http.Response, error) {
	return r.ApiService.ListWorkflowsExecute(r)
}

/*
ListWorkflows List registered workflows

List registered workflows

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiListWorkflowsRequest
*/
func (a *OrchestrationApiService) ListWorkflows(ctx context.Context) ApiListWorkflowsRequest {
	return ApiListWorkflowsRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return ListWorkflowsResponse
func (a *OrchestrationApiService) ListWorkflowsExecute(r ApiListWorkflowsRequest) (*ListWorkflowsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *ListWorkflowsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrchestrationApiService.ListWorkflows")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/orchestration/workflows"

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
			var v Error
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

type ApiOrchestrationgetServerInfoRequest struct {
	ctx context.Context
	ApiService OrchestrationApi
}

func (r ApiOrchestrationgetServerInfoRequest) Execute() (*ServerInfo, *http.Response, error) {
	return r.ApiService.OrchestrationgetServerInfoExecute(r)
}

/*
OrchestrationgetServerInfo Get server info

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiOrchestrationgetServerInfoRequest
*/
func (a *OrchestrationApiService) OrchestrationgetServerInfo(ctx context.Context) ApiOrchestrationgetServerInfoRequest {
	return ApiOrchestrationgetServerInfoRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return ServerInfo
func (a *OrchestrationApiService) OrchestrationgetServerInfoExecute(r ApiOrchestrationgetServerInfoRequest) (*ServerInfo, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *ServerInfo
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrchestrationApiService.OrchestrationgetServerInfo")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/orchestration/_info"

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
			var v Error
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

type ApiRunWorkflowRequest struct {
	ctx context.Context
	ApiService OrchestrationApi
	workflowID string
	wait *bool
	requestBody *map[string]string
}

// Wait end of the workflow before return
func (r ApiRunWorkflowRequest) Wait(wait bool) ApiRunWorkflowRequest {
	r.wait = &wait
	return r
}

func (r ApiRunWorkflowRequest) RequestBody(requestBody map[string]string) ApiRunWorkflowRequest {
	r.requestBody = &requestBody
	return r
}

func (r ApiRunWorkflowRequest) Execute() (*RunWorkflowResponse, *http.Response, error) {
	return r.ApiService.RunWorkflowExecute(r)
}

/*
RunWorkflow Run workflow

Run workflow

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param workflowID The flow id
 @return ApiRunWorkflowRequest
*/
func (a *OrchestrationApiService) RunWorkflow(ctx context.Context, workflowID string) ApiRunWorkflowRequest {
	return ApiRunWorkflowRequest{
		ApiService: a,
		ctx: ctx,
		workflowID: workflowID,
	}
}

// Execute executes the request
//  @return RunWorkflowResponse
func (a *OrchestrationApiService) RunWorkflowExecute(r ApiRunWorkflowRequest) (*RunWorkflowResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *RunWorkflowResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrchestrationApiService.RunWorkflow")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/orchestration/workflows/{workflowID}/instances"
	localVarPath = strings.Replace(localVarPath, "{"+"workflowID"+"}", url.PathEscape(parameterValueToString(r.workflowID, "workflowID")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.wait != nil {
		parameterAddToQuery(localVarQueryParams, "wait", r.wait, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

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
	// body params
	localVarPostBody = r.requestBody
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
			var v Error
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
