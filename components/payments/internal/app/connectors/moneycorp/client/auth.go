package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Cannot use "golang.org/x/oauth2/clientcredentials" lib because moneycorp
// is only accepting request with "application/json" content type, and the lib
// sets it as application/x-www-form-urlencoded, giving us a 415 error.
type apiTransport struct {
	logger logging.Logger

	clientID string
	apiKey   string
	endpoint string

	accessToken          string
	accessTokenExpiresAt time.Time

	underlying *otelhttp.Transport
}

func (t *apiTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.ensureAccessTokenIsValid(req.Context()); err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+t.accessToken)

	return t.underlying.RoundTrip(req)
}

func (t *apiTransport) ensureAccessTokenIsValid(ctx context.Context) error {
	if t.accessTokenExpiresAt.After(time.Now().Add(5 * time.Second)) {
		return nil
	}

	return t.login(ctx)
}

type loginRequest struct {
	ClientID string `json:"loginId"`
	APIKey   string `json:"apiKey"`
}

type loginResponse struct {
	Data struct {
		AccessToken string `json:"accessToken"`
		ExpiresIn   int    `json:"expiresIn"`
	} `json:"data"`
}

func (t *apiTransport) login(ctx context.Context) error {
	lreq := loginRequest{
		ClientID: t.clientID,
		APIKey:   t.apiKey,
	}

	requestBody, err := json.Marshal(lreq)
	if err != nil {
		return fmt.Errorf("failed to marshal login request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		t.endpoint+"/login", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			t.logger.Error(err)
		}
	}()

	var res loginResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("failed to decode login response: %w", err)
	}

	t.accessToken = res.Data.AccessToken
	t.accessTokenExpiresAt = time.Now().Add(time.Duration(res.Data.ExpiresIn) * time.Second)

	return nil
}
