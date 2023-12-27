package authorization

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"
	"net/http"
	"testing"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/auth/pkg/delegatedauth"
	authoidc "github.com/formancehq/auth/pkg/oidc"
	"github.com/formancehq/auth/pkg/storage/sqlstorage"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/oauth2-proxy/mockoidc"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zitadel/oidc/v2/pkg/op"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestVerifyAccessToken(t *testing.T) {
	t.Parallel()

	mockOIDC, err := mockoidc.Run()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, mockOIDC.Shutdown())
	}()

	// Prepare a tcp connection, listening on :0 to select a random port
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	// Compute server url, it will be the "issuer" of our oidc provider
	serverURL := fmt.Sprintf("http://%s", l.Addr().String())

	// Construct our storage
	postgresDB := pgtesting.NewPostgresDatabase(t)
	dialector := postgres.Open(postgresDB.ConnString())

	db, err := sqlstorage.LoadGorm(dialector, &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	require.NoError(t, sqlstorage.MigrateTables(context.Background(), db))

	storage := sqlstorage.New(db)

	serverRelyingParty, err := rp.NewRelyingPartyOIDC(mockOIDC.Issuer(), mockOIDC.ClientID, mockOIDC.ClientSecret,
		fmt.Sprintf("%s/authorize/callback", serverURL), []string{"openid", "email"})
	require.NoError(t, err)

	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	staticClients := []auth.StaticClient{{
		ClientOptions: auth.ClientOptions{
			Id:                     "test",
			Public:                 true,
			RedirectURIs:           []string{"http://localhost:3000/auth-callback"},
			Name:                   "test",
			PostLogoutRedirectUris: []string{"http://localhost:3000/"},
			Trusted:                true,
		},
	}}
	storageFacade := authoidc.NewStorageFacade(storage, serverRelyingParty, key, staticClients...)

	keySet, err := authoidc.ReadKeySet(http.DefaultClient, context.Background(), delegatedauth.Config{
		Issuer:       mockOIDC.Issuer(),
		ClientID:     mockOIDC.ClientID,
		ClientSecret: mockOIDC.ClientSecret,
	})
	require.NoError(t, err)

	provider, err := authoidc.NewOpenIDProvider(storageFacade, serverURL, mockOIDC.Issuer(), keySet)
	require.NoError(t, err)

	ar := &oidc.AuthRequest{
		ClientID: staticClients[0].Id,
	}
	authReq, err := provider.Storage().CreateAuthRequest(context.Background(), ar, "")
	require.NoError(t, err)

	client, err := provider.Storage().GetClientByClientID(context.Background(), authReq.GetClientID())
	require.NoError(t, err)

	tokenResponse, err := op.CreateTokenResponse(context.Background(), authReq, client, provider, true, "", "")
	require.NoError(t, err)

	t.Run("unprotected route", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/any", nil)
		require.NoError(t, err)
		require.NoError(t, verifyAccessToken(req, provider))
	})

	t.Run("protected routes", func(t *testing.T) {
		t.Parallel()

		protectedRoutes := []string{"/clients", "/users"}
		for _, route := range protectedRoutes {

			t.Run("no token", func(t *testing.T) {
				t.Parallel()

				req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, route, nil)
				require.NoError(t, err)

				err = verifyAccessToken(req, provider)
				require.Error(t, err)
				require.EqualError(t, err, ErrMissingAuthHeader.Error())
			})

			t.Run("malformed token", func(t *testing.T) {
				t.Parallel()

				req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, route, nil)
				require.NoError(t, err)

				req.Header.Set("Authorization", "malformed")
				err = verifyAccessToken(req, provider)
				require.Error(t, err)
				require.EqualError(t, err, ErrMalformedAuthHeader.Error())
			})

			t.Run("unverified token", func(t *testing.T) {
				t.Parallel()

				req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, route, nil)
				require.NoError(t, err)

				req.Header.Set("Authorization", oidc.PrefixBearer+"unverified")
				err = verifyAccessToken(req, provider)
				require.Error(t, err)
				require.EqualError(t, err, ErrVerifyAuthToken.Error())
			})

			t.Run("verified token", func(t *testing.T) {
				t.Parallel()

				req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, route, nil)
				require.NoError(t, err)

				req.Header.Set("Authorization", oidc.PrefixBearer+tokenResponse.AccessToken)
				require.NoError(t, verifyAccessToken(req, provider))
			})
		}
	})
}
