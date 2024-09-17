package httpwrapper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"
)

var (
	ErrStatusCodeUnexpected = errors.New("unexpected status code")

	defaultHttpErrorCheckerFn = func(statusCode int) error {
		if statusCode >= http.StatusBadRequest {
			return ErrStatusCodeUnexpected
		}
		return nil
	}
)

// Client is a convenience wrapper that encapsulates common code related to interacting with HTTP endpoints
type Client interface {
	// Do performs an HTTP request while handling errors and unmarshaling success and error responses into the provided interfaces
	// expectedBody and errorBody should be pointers to structs
	Do(req *http.Request, expectedBody, errorBody any) (statusCode int, err error)
}

type client struct {
	httpClient *http.Client

	httpErrorCheckerFn func(statusCode int) error
}

func NewClient(config *Config) (Client, error) {
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}
	if config.Transport != nil {
		config.Transport = otelhttp.NewTransport(config.Transport)
	} else {
		config.Transport = http.DefaultTransport.(*http.Transport).Clone()
	}

	httpClient := &http.Client{
		Timeout:   config.Timeout,
		Transport: config.Transport,
	}
	if config.OAuthConfig != nil {
		// pass a pre-configured http client to oauth lib via the context
		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
		httpClient = config.OAuthConfig.Client(ctx)
	}

	if config.HttpErrorCheckerFn == nil {
		config.HttpErrorCheckerFn = defaultHttpErrorCheckerFn
	}

	return &client{
		httpErrorCheckerFn: config.HttpErrorCheckerFn,
		httpClient:         httpClient,
	}, nil
}

func (c *client) Do(req *http.Request, expectedBody, errorBody any) (int, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to make request: %w", err)
	}

	reqErr := c.httpErrorCheckerFn(resp.StatusCode)
	// the caller doesn't care about the response body so we return early
	if resp.Body == nil || (reqErr == nil && expectedBody == nil) || (reqErr != nil && errorBody == nil) {
		return resp.StatusCode, reqErr
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			_ = err
			// TODO(polo): log error
		}
	}()

	// TODO: reading everything into memory might not be optimal if we expect long responses
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	if reqErr != nil {
		if err = json.Unmarshal(rawBody, errorBody); err != nil {
			return resp.StatusCode, fmt.Errorf("failed to unmarshal error response with status %d: %w", resp.StatusCode, err)
		}
		return resp.StatusCode, reqErr
	}

	// TODO: assuming json bodies for now, but may need to handle other body types
	if err = json.Unmarshal(rawBody, expectedBody); err != nil {
		return resp.StatusCode, fmt.Errorf("failed to unmarshal response with status %d: %w", resp.StatusCode, err)
	}
	return resp.StatusCode, nil
}
