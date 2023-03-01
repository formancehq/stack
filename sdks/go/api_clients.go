/*
Formance Stack API

Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 

API version: v1.0.20230301
Contact: support@formance.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package formance

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)


type ClientsApi interface {

	/*
	AddScopeToClient Add scope to client

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clientId Client ID
	@param scopeId Scope ID
	@return ApiAddScopeToClientRequest
	*/
	AddScopeToClient(ctx context.Context, clientId string, scopeId string) ApiAddScopeToClientRequest

	// AddScopeToClientExecute executes the request
	AddScopeToClientExecute(r ApiAddScopeToClientRequest) (*http.Response, error)

	/*
	CreateClient Create client

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiCreateClientRequest
	*/
	CreateClient(ctx context.Context) ApiCreateClientRequest

	// CreateClientExecute executes the request
	//  @return CreateClientResponse
	CreateClientExecute(r ApiCreateClientRequest) (*CreateClientResponse, *http.Response, error)

	/*
	CreateSecret Add a secret to a client

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clientId Client ID
	@return ApiCreateSecretRequest
	*/
	CreateSecret(ctx context.Context, clientId string) ApiCreateSecretRequest

	// CreateSecretExecute executes the request
	//  @return CreateSecretResponse
	CreateSecretExecute(r ApiCreateSecretRequest) (*CreateSecretResponse, *http.Response, error)

	/*
	DeleteClient Delete client

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clientId Client ID
	@return ApiDeleteClientRequest
	*/
	DeleteClient(ctx context.Context, clientId string) ApiDeleteClientRequest

	// DeleteClientExecute executes the request
	DeleteClientExecute(r ApiDeleteClientRequest) (*http.Response, error)

	/*
	DeleteScopeFromClient Delete scope from client

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clientId Client ID
	@param scopeId Scope ID
	@return ApiDeleteScopeFromClientRequest
	*/
	DeleteScopeFromClient(ctx context.Context, clientId string, scopeId string) ApiDeleteScopeFromClientRequest

	// DeleteScopeFromClientExecute executes the request
	DeleteScopeFromClientExecute(r ApiDeleteScopeFromClientRequest) (*http.Response, error)

	/*
	DeleteSecret Delete a secret from a client

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clientId Client ID
	@param secretId Secret ID
	@return ApiDeleteSecretRequest
	*/
	DeleteSecret(ctx context.Context, clientId string, secretId string) ApiDeleteSecretRequest

	// DeleteSecretExecute executes the request
	DeleteSecretExecute(r ApiDeleteSecretRequest) (*http.Response, error)

	/*
	ListClients List clients

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiListClientsRequest
	*/
	ListClients(ctx context.Context) ApiListClientsRequest

	// ListClientsExecute executes the request
	//  @return ListClientsResponse
	ListClientsExecute(r ApiListClientsRequest) (*ListClientsResponse, *http.Response, error)

	/*
	ReadClient Read client

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clientId Client ID
	@return ApiReadClientRequest
	*/
	ReadClient(ctx context.Context, clientId string) ApiReadClientRequest

	// ReadClientExecute executes the request
	//  @return ReadClientResponse
	ReadClientExecute(r ApiReadClientRequest) (*ReadClientResponse, *http.Response, error)

	/*
	UpdateClient Update client

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clientId Client ID
	@return ApiUpdateClientRequest
	*/
	UpdateClient(ctx context.Context, clientId string) ApiUpdateClientRequest

	// UpdateClientExecute executes the request
	//  @return CreateClientResponse
	UpdateClientExecute(r ApiUpdateClientRequest) (*CreateClientResponse, *http.Response, error)
}

// ClientsApiService ClientsApi service
type ClientsApiService service

type ApiAddScopeToClientRequest struct {
	ctx context.Context
	ApiService ClientsApi
	clientId string
	scopeId string
}

func (r ApiAddScopeToClientRequest) Execute() (*http.Response, error) {
	return r.ApiService.AddScopeToClientExecute(r)
}

/*
AddScopeToClient Add scope to client

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param clientId Client ID
 @param scopeId Scope ID
 @return ApiAddScopeToClientRequest
*/
func (a *ClientsApiService) AddScopeToClient(ctx context.Context, clientId string, scopeId string) ApiAddScopeToClientRequest {
	return ApiAddScopeToClientRequest{
		ApiService: a,
		ctx: ctx,
		clientId: clientId,
		scopeId: scopeId,
	}
}

// Execute executes the request
func (a *ClientsApiService) AddScopeToClientExecute(r ApiAddScopeToClientRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClientsApiService.AddScopeToClient")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/clients/{clientId}/scopes/{scopeId}"
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(parameterValueToString(r.clientId, "clientId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"scopeId"+"}", url.PathEscape(parameterValueToString(r.scopeId, "scopeId")), -1)

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
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiCreateClientRequest struct {
	ctx context.Context
	ApiService ClientsApi
	body *ClientOptions
}

func (r ApiCreateClientRequest) Body(body ClientOptions) ApiCreateClientRequest {
	r.body = &body
	return r
}

func (r ApiCreateClientRequest) Execute() (*CreateClientResponse, *http.Response, error) {
	return r.ApiService.CreateClientExecute(r)
}

/*
CreateClient Create client

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiCreateClientRequest
*/
func (a *ClientsApiService) CreateClient(ctx context.Context) ApiCreateClientRequest {
	return ApiCreateClientRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return CreateClientResponse
func (a *ClientsApiService) CreateClientExecute(r ApiCreateClientRequest) (*CreateClientResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *CreateClientResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClientsApiService.CreateClient")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/clients"

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

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
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

type ApiCreateSecretRequest struct {
	ctx context.Context
	ApiService ClientsApi
	clientId string
	body *SecretOptions
}

func (r ApiCreateSecretRequest) Body(body SecretOptions) ApiCreateSecretRequest {
	r.body = &body
	return r
}

func (r ApiCreateSecretRequest) Execute() (*CreateSecretResponse, *http.Response, error) {
	return r.ApiService.CreateSecretExecute(r)
}

/*
CreateSecret Add a secret to a client

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param clientId Client ID
 @return ApiCreateSecretRequest
*/
func (a *ClientsApiService) CreateSecret(ctx context.Context, clientId string) ApiCreateSecretRequest {
	return ApiCreateSecretRequest{
		ApiService: a,
		ctx: ctx,
		clientId: clientId,
	}
}

// Execute executes the request
//  @return CreateSecretResponse
func (a *ClientsApiService) CreateSecretExecute(r ApiCreateSecretRequest) (*CreateSecretResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *CreateSecretResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClientsApiService.CreateSecret")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/clients/{clientId}/secrets"
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(parameterValueToString(r.clientId, "clientId")), -1)

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

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
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

type ApiDeleteClientRequest struct {
	ctx context.Context
	ApiService ClientsApi
	clientId string
}

func (r ApiDeleteClientRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteClientExecute(r)
}

/*
DeleteClient Delete client

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param clientId Client ID
 @return ApiDeleteClientRequest
*/
func (a *ClientsApiService) DeleteClient(ctx context.Context, clientId string) ApiDeleteClientRequest {
	return ApiDeleteClientRequest{
		ApiService: a,
		ctx: ctx,
		clientId: clientId,
	}
}

// Execute executes the request
func (a *ClientsApiService) DeleteClientExecute(r ApiDeleteClientRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClientsApiService.DeleteClient")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/clients/{clientId}"
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(parameterValueToString(r.clientId, "clientId")), -1)

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
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiDeleteScopeFromClientRequest struct {
	ctx context.Context
	ApiService ClientsApi
	clientId string
	scopeId string
}

func (r ApiDeleteScopeFromClientRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteScopeFromClientExecute(r)
}

/*
DeleteScopeFromClient Delete scope from client

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param clientId Client ID
 @param scopeId Scope ID
 @return ApiDeleteScopeFromClientRequest
*/
func (a *ClientsApiService) DeleteScopeFromClient(ctx context.Context, clientId string, scopeId string) ApiDeleteScopeFromClientRequest {
	return ApiDeleteScopeFromClientRequest{
		ApiService: a,
		ctx: ctx,
		clientId: clientId,
		scopeId: scopeId,
	}
}

// Execute executes the request
func (a *ClientsApiService) DeleteScopeFromClientExecute(r ApiDeleteScopeFromClientRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClientsApiService.DeleteScopeFromClient")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/clients/{clientId}/scopes/{scopeId}"
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(parameterValueToString(r.clientId, "clientId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"scopeId"+"}", url.PathEscape(parameterValueToString(r.scopeId, "scopeId")), -1)

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
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiDeleteSecretRequest struct {
	ctx context.Context
	ApiService ClientsApi
	clientId string
	secretId string
}

func (r ApiDeleteSecretRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteSecretExecute(r)
}

/*
DeleteSecret Delete a secret from a client

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param clientId Client ID
 @param secretId Secret ID
 @return ApiDeleteSecretRequest
*/
func (a *ClientsApiService) DeleteSecret(ctx context.Context, clientId string, secretId string) ApiDeleteSecretRequest {
	return ApiDeleteSecretRequest{
		ApiService: a,
		ctx: ctx,
		clientId: clientId,
		secretId: secretId,
	}
}

// Execute executes the request
func (a *ClientsApiService) DeleteSecretExecute(r ApiDeleteSecretRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClientsApiService.DeleteSecret")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/clients/{clientId}/secrets/{secretId}"
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(parameterValueToString(r.clientId, "clientId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"secretId"+"}", url.PathEscape(parameterValueToString(r.secretId, "secretId")), -1)

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
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiListClientsRequest struct {
	ctx context.Context
	ApiService ClientsApi
}

func (r ApiListClientsRequest) Execute() (*ListClientsResponse, *http.Response, error) {
	return r.ApiService.ListClientsExecute(r)
}

/*
ListClients List clients

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiListClientsRequest
*/
func (a *ClientsApiService) ListClients(ctx context.Context) ApiListClientsRequest {
	return ApiListClientsRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return ListClientsResponse
func (a *ClientsApiService) ListClientsExecute(r ApiListClientsRequest) (*ListClientsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *ListClientsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClientsApiService.ListClients")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/clients"

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

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
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

type ApiReadClientRequest struct {
	ctx context.Context
	ApiService ClientsApi
	clientId string
}

func (r ApiReadClientRequest) Execute() (*ReadClientResponse, *http.Response, error) {
	return r.ApiService.ReadClientExecute(r)
}

/*
ReadClient Read client

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param clientId Client ID
 @return ApiReadClientRequest
*/
func (a *ClientsApiService) ReadClient(ctx context.Context, clientId string) ApiReadClientRequest {
	return ApiReadClientRequest{
		ApiService: a,
		ctx: ctx,
		clientId: clientId,
	}
}

// Execute executes the request
//  @return ReadClientResponse
func (a *ClientsApiService) ReadClientExecute(r ApiReadClientRequest) (*ReadClientResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *ReadClientResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClientsApiService.ReadClient")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/clients/{clientId}"
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(parameterValueToString(r.clientId, "clientId")), -1)

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

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
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

type ApiUpdateClientRequest struct {
	ctx context.Context
	ApiService ClientsApi
	clientId string
	body *ClientOptions
}

func (r ApiUpdateClientRequest) Body(body ClientOptions) ApiUpdateClientRequest {
	r.body = &body
	return r
}

func (r ApiUpdateClientRequest) Execute() (*CreateClientResponse, *http.Response, error) {
	return r.ApiService.UpdateClientExecute(r)
}

/*
UpdateClient Update client

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param clientId Client ID
 @return ApiUpdateClientRequest
*/
func (a *ClientsApiService) UpdateClient(ctx context.Context, clientId string) ApiUpdateClientRequest {
	return ApiUpdateClientRequest{
		ApiService: a,
		ctx: ctx,
		clientId: clientId,
	}
}

// Execute executes the request
//  @return CreateClientResponse
func (a *ClientsApiService) UpdateClientExecute(r ApiUpdateClientRequest) (*CreateClientResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPut
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *CreateClientResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ClientsApiService.UpdateClient")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/auth/clients/{clientId}"
	localVarPath = strings.Replace(localVarPath, "{"+"clientId"+"}", url.PathEscape(parameterValueToString(r.clientId, "clientId")), -1)

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

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
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
