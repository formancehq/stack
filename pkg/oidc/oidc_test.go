package oidc_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/auth/pkg/oidc"
	"github.com/formancehq/auth/pkg/storage/sqlstorage"
	"github.com/gorilla/mux"
	"github.com/oauth2-proxy/mockoidc"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/client/rs"
	"github.com/zitadel/oidc/pkg/op"
	"gorm.io/driver/sqlite"
)

func init() {
	os.Setenv(op.OidcDevMode, "true")
}

func withServer(t *testing.T, fn func(storage *sqlstorage.Storage, provider op.OpenIDProvider)) {
	// Create a mock OIDC server which will always return a default user
	mockOIDC, err := mockoidc.Run()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, mockOIDC.Shutdown())
	}()

	// Prepare a tcp connection, listening on :0 to select a random port
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	// Compute server url, it will be the "issuer" of our oidc provider
	serverUrl := fmt.Sprintf("http://%s", l.Addr().String())

	// As our oidc provider, is also a relying party (it delegates authentication), we need to construct a relying party
	// with information from the mock
	serverRelyingParty, err := rp.NewRelyingPartyOIDC(mockOIDC.Issuer(), mockOIDC.ClientID, mockOIDC.ClientSecret,
		fmt.Sprintf("%s/authorize/callback", serverUrl), []string{"openid", "email"})
	require.NoError(t, err)

	// Construct our storage
	db, err := sqlstorage.LoadGorm(sqlite.Open(":memory:"))
	require.NoError(t, err)
	require.NoError(t, sqlstorage.MigrateTables(context.Background(), db))

	storage := sqlstorage.New(db)

	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	storageFacade := oidc.NewStorageFacade(storage, serverRelyingParty, key)

	// Construct our oidc provider
	provider, err := oidc.NewOpenIDProvider(context.TODO(), storageFacade, serverUrl)
	require.NoError(t, err)

	u, err := url.Parse(serverUrl)
	require.NoError(t, err)

	// Create the router
	router := mux.NewRouter()
	oidc.AddRoutes(router, provider, storage, serverRelyingParty, u)

	// Create our http server for our oidc provider
	providerHttpServer := &http.Server{
		Handler: router,
	}
	go func() {
		err := providerHttpServer.Serve(l)
		if err != http.ErrServerClosed {
			require.Fail(t, err.Error())
		}
	}()
	defer providerHttpServer.Close()

	fn(storage, provider)
}

func Test3LeggedFlow(t *testing.T) {

	withServer(t, func(storage *sqlstorage.Storage, provider op.OpenIDProvider) {
		// Create ou http server for our client (a web application for example)
		code := make(chan string, 1) // Just store codes coming from our provider inside a chan
		clientHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code <- r.URL.Query().Get("code")
		})
		clientHttpServer := httptest.NewServer(clientHandler)
		defer clientHttpServer.Close()

		// Create a OAuth2 client which represent our client application
		client := auth.NewClient(auth.ClientOptions{})
		client.RedirectURIs.Append(clientHttpServer.URL)          // Need to configure the redirect uri
		_, clear := client.GenerateNewSecret(auth.SecretCreate{}) // Need to generate a secret
		require.NoError(t, storage.SaveClient(context.TODO(), *client))

		// As our client is a relying party, we can use the library to get some helpers
		clientRelyingParty, err := rp.NewRelyingPartyOIDC(provider.Issuer(), client.Id, clear, client.RedirectURIs[0], []string{"openid", "email"})
		require.NoError(t, err)

		// Trigger an authentication request
		authUrl := rp.AuthURL("", clientRelyingParty)
		rsp, err := (&http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				fmt.Println(req.URL.String())
				return nil
			},
		}).Get(authUrl)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rsp.StatusCode)

		select {
		// As the mock automatically accept login response, we should have received a code
		case code := <-code:
			// And this code is used to get a token
			tokens, err := rp.CodeExchange(context.TODO(), code, clientRelyingParty)
			require.NoError(t, err)

			// Create a OAuth2 client which represent our client application
			secondaryClient := auth.NewClient(auth.ClientOptions{
				Trusted: true,
			})
			_, clear = secondaryClient.GenerateNewSecret(auth.SecretCreate{}) // Need to generate a secret
			require.NoError(t, storage.SaveClient(context.TODO(), *secondaryClient))

			resourceServer, err := rs.NewResourceServerClientCredentials(provider.Issuer(), secondaryClient.Id, clear)
			require.NoError(t, err)

			introspection, err := rs.Introspect(context.TODO(), resourceServer, tokens.AccessToken)
			require.NoError(t, err)
			require.True(t, introspection.IsActive())
		default:
			require.Fail(t, "code was expected")
		}
	})
}
