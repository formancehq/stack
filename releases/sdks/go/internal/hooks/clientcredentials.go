// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package hooks

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type session struct {
	Credentials *credentials
	Token       string
	ExpiresAt   *int64
	Scopes      []string
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   *int64 `json:"expires_in"`
}

type credentials struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
}

type clientCredentialsHook struct {
	baseURL  string
	client   HTTPClient
	sessions map[string]*session
}

var (
	_ sdkInitHook       = (*clientCredentialsHook)(nil)
	_ beforeRequestHook = (*clientCredentialsHook)(nil)
	_ afterErrorHook    = (*clientCredentialsHook)(nil)
)

func NewClientCredentialsHook() *clientCredentialsHook {
	return &clientCredentialsHook{
		sessions: make(map[string]*session),
	}
}

func (c *clientCredentialsHook) SDKInit(baseURL string, client HTTPClient) (string, HTTPClient) {
	c.baseURL = baseURL
	c.client = client
	return baseURL, client
}

func (c *clientCredentialsHook) BeforeRequest(ctx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	if ctx.OAuth2Scopes == nil {
		// OAuth2 not in use
		return req, nil
	}

	credentials, err := c.getCredentials(ctx.Context, ctx.SecuritySource)
	if err != nil {
		return nil, &FailEarly{Cause: err}
	}
	if credentials == nil {
		return req, err
	}

	sessionKey := getSessionKey(credentials.ClientID, credentials.ClientSecret)
	sess, ok := c.sessions[sessionKey]
	if !ok || !hasRequiredScopes(sess.Scopes, ctx.OAuth2Scopes) || hasTokenExpired(sess.ExpiresAt) {
		s, err := c.doTokenRequest(ctx.Context, credentials, getScopes(ctx.OAuth2Scopes, sess))
		if err != nil {
			return nil, fmt.Errorf("failed to get token: %w", err)
		}

		c.sessions[sessionKey] = s
		sess = s
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sess.Token))

	return req, nil
}

func (c *clientCredentialsHook) AfterError(ctx AfterErrorContext, res *http.Response, err error) (*http.Response, error) {
	if ctx.OAuth2Scopes == nil {
		// OAuth2 not in use
		return res, err
	}

	// We don't want to refresh the token if the error is not related to the token
	if err != nil {
		return res, err
	}

	credentials, err := c.getCredentials(ctx.Context, ctx.SecuritySource)
	if err != nil {
		return nil, &FailEarly{Cause: err}
	}
	if credentials == nil {
		return res, err
	}

	if res != nil && res.StatusCode == http.StatusUnauthorized {
		sessionKey := getSessionKey(credentials.ClientID, credentials.ClientSecret)
		delete(c.sessions, sessionKey)
	}

	return res, err
}

func (c *clientCredentialsHook) doTokenRequest(ctx context.Context, credentials *credentials, scopes []string) (*session, error) {
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", credentials.ClientID)
	values.Set("client_secret", credentials.ClientSecret)

	if len(scopes) > 0 {
		values.Set("scope", strings.Join(scopes, " "))
	}

	tokenURL := credentials.TokenURL
	u, err := url.Parse(tokenURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token URL: %w", err)
	}
	if !u.IsAbs() {
		tokenURL, err = url.JoinPath(c.baseURL, tokenURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse token URL: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewBufferString(values.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send token request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("unexpected status code: %d: %s", res.StatusCode, body)
	}

	var tokenRes tokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenRes); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenRes.TokenType != "Bearer" {
		return nil, fmt.Errorf("unexpected token type: %s", tokenRes.TokenType)
	}

	var expiresAt *int64
	if tokenRes.ExpiresIn != nil {
		expiresAt = new(int64)
		*expiresAt = time.Now().Unix() + *tokenRes.ExpiresIn
	}

	return &session{
		Credentials: credentials,
		Token:       tokenRes.AccessToken,
		ExpiresAt:   expiresAt,
		Scopes:      scopes,
	}, nil
}

func (c *clientCredentialsHook) getCredentials(ctx context.Context, source func(ctx context.Context) (interface{}, error)) (*credentials, error) {
	if source == nil {
		return nil, nil
	}

	sec, err := source(ctx)
	if err != nil {
		return nil, err
	}

	security, ok := sec.(shared.Security)

	if !ok {
		return nil, fmt.Errorf("unexpected security type: %T", sec)
	}

	if security.ClientID == nil || security.ClientSecret == nil {
		return nil, nil
	}

	return &credentials{
		ClientID:     *security.ClientID,
		ClientSecret: *security.ClientSecret,
		TokenURL:     *security.GetTokenURL(),
	}, nil
}

func getSessionKey(clientID, clientSecret string) string {
	key := fmt.Sprintf("%s:%s", clientID, clientSecret)
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

func hasRequiredScopes(scopes []string, requiredScopes []string) bool {
	for _, requiredScope := range requiredScopes {
		found := false
		for _, scope := range scopes {
			if scope == requiredScope {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func getScopes(requiredScopes []string, sess *session) []string {
	scopes := requiredScopes
	if sess != nil {
		for _, scope := range sess.Scopes {
			found := false
			for _, requiredScope := range requiredScopes {
				if scope == requiredScope {
					found = true
					break
				}
			}
			if !found {
				scopes = append(scopes, scope)
			}
		}
	}

	return scopes
}

func hasTokenExpired(expiresAt *int64) bool {
	return expiresAt == nil || time.Now().Unix()+60 >= *expiresAt
}
