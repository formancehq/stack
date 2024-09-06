package httpclient

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/oauth2-proxy/mockoidc"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestClient(t *testing.T) {
	t.Parallel()

	// Create a mock OIDC server which will always return a default user
	mockOIDC, err := mockoidc.NewServer(nil)
	require.NoError(t, err)

	// mockoidc does not support client credentials, add a middleware to mock the feature
	mockCalled := make(chan struct{})
	err = mockOIDC.AddMiddleware(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/oidc/.well-known/openid-configuration" {
				handler.ServeHTTP(w, r)
				return
			}

			w.Header().Set("Content-Type", "application/json")

			tokenEndpoint, err := url.Parse(mockOIDC.TokenEndpoint())
			require.NoError(t, err)

			require.NoError(t, r.ParseForm())
			require.Equal(t, tokenEndpoint.Path, r.URL.Path)
			require.Equal(t, "client_credentials", r.Form.Get("grant_type"))
			username, password, ok := r.BasicAuth()
			require.True(t, ok)
			require.Equal(t, mockOIDC.ClientID, username)
			require.Equal(t, mockOIDC.ClientSecret, password)

			require.NoError(t, json.NewEncoder(w).Encode(oauth2.Token{
				AccessToken: "xxx",
				TokenType:   "bearer",
			}))

			close(mockCalled)
		})
	})
	require.NoError(t, err)

	// Start a tcp socket
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	// Start the oidc mock server on or newly created tcp socket
	require.NoError(t, mockOIDC.Start(ln, nil))
	t.Cleanup(func() {
		_ = mockOIDC.Shutdown()
	})

	// Client should be created properly
	client, err := NewStackAuthenticatedClient(
		OAuth2Config{
			Issuer:       mockOIDC.Issuer(),
			ClientID:     mockOIDC.ClientID,
			ClientSecret: mockOIDC.ClientSecret,
		},
		http.DefaultClient,
		"module1",
		"module1",
	)
	require.NoError(t, err)

	// Create a testing server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	t.Cleanup(srv.Close)

	// Do a simple request to check if the oauth2 client is working properly
	rsp, err := client.Get(srv.URL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rsp.StatusCode)

	// Check that the mock server has been called
	select {
	case <-mockCalled:
	default:
		require.Fail(t, "oidc mock should have been called")
	}
}
