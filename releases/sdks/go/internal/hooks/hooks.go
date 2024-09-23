// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package hooks

import (
	"context"
	"errors"
	"net/http"
)

type FailEarly struct {
	Cause error
}

var _ error = (*FailEarly)(nil)

func (f *FailEarly) Error() string {
	return f.Cause.Error()
}

// HTTPClient provides an interface for supplying the SDK with a custom HTTP client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HookContext struct {
	Context        context.Context
	OperationID    string
	OAuth2Scopes   []string
	SecuritySource func(context.Context) (interface{}, error)
}

type BeforeRequestContext struct {
	HookContext
}

type AfterSuccessContext struct {
	HookContext
}

type AfterErrorContext struct {
	HookContext
}

// sdkInitHook is called when the SDK is initializing. The hook can modify and return a new baseURL and HTTP client to be used by the SDK.
type sdkInitHook interface {
	SDKInit(baseURL string, client HTTPClient) (string, HTTPClient)
}

// beforeRequestHook is called before the SDK sends a request. The hook can modify the request before it is sent or return an error to stop the request from being sent.
type beforeRequestHook interface {
	BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error)
}

// afterSuccessHook is called after the SDK receives a response. The hook can modify the response before it is handled or return an error to stop the response from being handled.
type afterSuccessHook interface {
	AfterSuccess(hookCtx AfterSuccessContext, res *http.Response) (*http.Response, error)
}

// afterErrorHook is called after the SDK encounters an error, or a non-successful response. The hook can modify the response if available otherwise modify the error.
// All afterErrorHook hooks are called and returning an error won't stop the other hooks from being called. But if you want to stop the other hooks from being called, you can return a FailEarly error wrapping your error.
type afterErrorHook interface {
	AfterError(hookCtx AfterErrorContext, res *http.Response, err error) (*http.Response, error)
}

type Hooks struct {
	sdkInitHooks      []sdkInitHook
	beforeRequestHook []beforeRequestHook
	afterSuccessHook  []afterSuccessHook
	afterErrorHook    []afterErrorHook
}

func New() *Hooks {
	cc := NewClientCredentialsHook()

	h := &Hooks{
		sdkInitHooks: []sdkInitHook{
			cc,
		},
		beforeRequestHook: []beforeRequestHook{
			cc,
		},
		afterSuccessHook: []afterSuccessHook{},
		afterErrorHook: []afterErrorHook{
			cc,
		},
	}

	initHooks(h)

	return h
}

// registerSDKInitHook registers a hook to be used by the SDK for the initialization event.
func (h *Hooks) registerSDKInitHook(hook sdkInitHook) {
	h.sdkInitHooks = append(h.sdkInitHooks, hook)
}

// registerBeforeRequestHook registers a hook to be used by the SDK for the before request event.
func (h *Hooks) registerBeforeRequestHook(hook beforeRequestHook) {
	h.beforeRequestHook = append(h.beforeRequestHook, hook)
}

// registerAfterSuccessHook registers a hook to be used by the SDK for the after success event.
func (h *Hooks) registerAfterSuccessHook(hook afterSuccessHook) {
	h.afterSuccessHook = append(h.afterSuccessHook, hook)
}

// registerAfterErrorHook registers a hook to be used by the SDK for the after error event.
func (h *Hooks) registerAfterErrorHook(hook afterErrorHook) {
	h.afterErrorHook = append(h.afterErrorHook, hook)
}

func (h *Hooks) SDKInit(baseURL string, client HTTPClient) (string, HTTPClient) {
	for _, hook := range h.sdkInitHooks {
		baseURL, client = hook.SDKInit(baseURL, client)
	}
	return baseURL, client
}

func (h *Hooks) BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	for _, hook := range h.beforeRequestHook {
		var err error
		req, err = hook.BeforeRequest(hookCtx, req)
		if err != nil {
			return req, err
		}
	}
	return req, nil
}

func (h *Hooks) AfterSuccess(hookCtx AfterSuccessContext, res *http.Response) (*http.Response, error) {
	for _, hook := range h.afterSuccessHook {
		var err error
		res, err = hook.AfterSuccess(hookCtx, res)
		if err != nil {
			return res, err
		}
	}
	return res, nil
}

func (h *Hooks) AfterError(hookCtx AfterErrorContext, res *http.Response, err error) (*http.Response, error) {
	for _, hook := range h.afterErrorHook {
		res, err = hook.AfterError(hookCtx, res, err)
		var fe *FailEarly
		if errors.As(err, &fe) {
			return nil, fe.Cause
		}
	}
	return res, err
}
