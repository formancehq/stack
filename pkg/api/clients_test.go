package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	auth "github.com/numary/auth/pkg"
	"github.com/numary/auth/pkg/storage"
	"github.com/numary/go-libs/sharedapi"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/oidc/pkg/op"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func withDbAndClientRouter(t *testing.T, callback func(router *mux.Router, db *gorm.DB)) {
	db, err := storage.LoadGorm(sqlite.Open(":memory:"))
	require.NoError(t, err)
	require.NoError(t, storage.MigrateTables(context.Background(), db))

	router := mux.NewRouter()
	addClientRoutes(db, router)

	callback(router, db)
}

func createJSONBuffer(t *testing.T, v any) io.Reader {
	data, err := json.Marshal(v)
	require.NoError(t, err)

	return bytes.NewBuffer(data)
}

func readObject[T any](t *testing.T, recorder *httptest.ResponseRecorder) T {
	body := sharedapi.BaseResponse[T]{}
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	return *body.Data
}

func TestCreateClient(t *testing.T) {

	type testCase struct {
		name           string
		options        auth.ClientOptions
		expectedClient auth.Client
	}
	for _, tc := range []testCase{
		{
			name: "confidential client",
			options: auth.ClientOptions{
				Name:                   "confidential client",
				Scopes:                 auth.Scopes,
				RedirectUris:           []string{"http://localhost:8080"},
				Description:            "abc",
				PostLogoutRedirectUris: []string{"http://localhost:8080/logout"},
			},
			expectedClient: auth.Client{
				GrantTypes: auth.Array[oidc.GrantType]{
					oidc.GrantTypeCode,
					oidc.GrantTypeRefreshToken,
					oidc.GrantTypeClientCredentials,
				},
				AccessTokenType: op.AccessTokenTypeJWT,
				ResponseTypes: auth.Array[oidc.ResponseType]{
					oidc.ResponseTypeCode,
				},
				AuthMethod:             oidc.AuthMethodNone,
				Scopes:                 auth.Scopes,
				RedirectURIs:           []string{"http://localhost:8080"},
				Description:            "abc",
				PostLogoutRedirectUris: []string{"http://localhost:8080/logout"},
			},
		},
		{
			name: "public client",
			options: auth.ClientOptions{
				Name:   "public client",
				Public: true,
			},
			expectedClient: auth.Client{
				GrantTypes: auth.Array[oidc.GrantType]{
					oidc.GrantTypeCode,
					oidc.GrantTypeRefreshToken,
				},
				AccessTokenType: op.AccessTokenTypeJWT,
				ResponseTypes: auth.Array[oidc.ResponseType]{
					oidc.ResponseTypeCode,
				},
				AuthMethod: oidc.AuthMethodNone,
			},
		},
	} {
		withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {
			req := httptest.NewRequest(http.MethodPost, "/clients", createJSONBuffer(t, tc.options))
			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)

			require.Equal(t, http.StatusCreated, res.Code)

			createdClient := readObject[client](t, res)
			require.NotEmpty(t, createdClient.ID)
			require.Equal(t, tc.options, createdClient.ClientOptions)

			tc.expectedClient.Id = createdClient.ID

			clientFromDatabase := auth.Client{}
			require.NoError(t, db.Find(&clientFromDatabase, "id = ?", createdClient.ID).Error)
			require.Equal(t, tc.expectedClient, clientFromDatabase)
		})
	}
}

func TestUpdateClient(t *testing.T) {

	type testCase struct {
		name           string
		options        auth.ClientOptions
		expectedClient auth.Client
	}
	for _, tc := range []testCase{
		{
			name: "confidential client",
			options: auth.ClientOptions{
				Name:                   "confidential client",
				Scopes:                 auth.Scopes,
				RedirectUris:           []string{"http://localhost:8080"},
				Description:            "abc",
				PostLogoutRedirectUris: []string{"http://localhost:8080/logout"},
			},
			expectedClient: auth.Client{
				GrantTypes: auth.Array[oidc.GrantType]{
					oidc.GrantTypeCode,
					oidc.GrantTypeRefreshToken,
					oidc.GrantTypeClientCredentials,
				},
				AccessTokenType: op.AccessTokenTypeJWT,
				ResponseTypes: auth.Array[oidc.ResponseType]{
					oidc.ResponseTypeCode,
				},
				AuthMethod:             oidc.AuthMethodNone,
				Scopes:                 auth.Scopes,
				RedirectURIs:           []string{"http://localhost:8080"},
				Description:            "abc",
				PostLogoutRedirectUris: []string{"http://localhost:8080/logout"},
			},
		},
		{
			name: "public client",
			options: auth.ClientOptions{
				Name:   "public client",
				Public: true,
			},
			expectedClient: auth.Client{
				GrantTypes: auth.Array[oidc.GrantType]{
					oidc.GrantTypeCode,
					oidc.GrantTypeRefreshToken,
				},
				AccessTokenType: op.AccessTokenTypeJWT,
				ResponseTypes: auth.Array[oidc.ResponseType]{
					oidc.ResponseTypeCode,
				},
				AuthMethod: oidc.AuthMethodNone,
			},
		},
	} {
		withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {

			initialClient := auth.NewClient(auth.ClientOptions{})
			require.NoError(t, db.Create(initialClient).Error)

			req := httptest.NewRequest(http.MethodPut, "/clients/"+initialClient.Id, createJSONBuffer(t, tc.options))
			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)

			require.Equal(t, http.StatusOK, res.Code)

			updatedClient := readObject[client](t, res)
			require.NotEmpty(t, updatedClient.ID)
			require.Equal(t, tc.options, updatedClient.ClientOptions)

			tc.expectedClient.Id = updatedClient.ID

			clientFromDatabase := auth.Client{}
			require.NoError(t, db.Find(&clientFromDatabase, "id = ?", updatedClient.ID).Error)
			require.Equal(t, tc.expectedClient, clientFromDatabase)
		})
	}
}

func TestListClients(t *testing.T) {
	withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {
		client1 := auth.NewClient(auth.ClientOptions{})
		require.NoError(t, db.Create(client1).Error)

		client2 := auth.NewClient(auth.ClientOptions{})
		require.NoError(t, db.Create(client2).Error)

		req := httptest.NewRequest(http.MethodGet, "/clients", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusOK, res.Code)

		clients := readObject[[]client](t, res)
		require.Len(t, clients, 2)
	})
}

func TestReadClient(t *testing.T) {
	withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {
		client := auth.NewClient(auth.ClientOptions{})
		require.NoError(t, db.Create(client).Error)

		req := httptest.NewRequest(http.MethodGet, "/clients/"+client.Id, nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusOK, res.Code)
	})
}

func TestGenerateNewSecret(t *testing.T) {
	withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {
		client := auth.NewClient(auth.ClientOptions{})
		require.NoError(t, db.Create(client).Error)

		req := httptest.NewRequest(http.MethodPost, "/clients/"+client.Id+"/secrets", createJSONBuffer(t, secretCreate{
			Name: "secret1",
		}))
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		result := readObject[secretCreateResult](t, res)
		require.NotEmpty(t, result.Clear)
		require.Equal(t, result.LastDigits, result.Clear[len(result.Clear)-4:])
		require.Equal(t, result.Name, "secret1")

		require.Equal(t, http.StatusOK, res.Code)
		require.NoError(t, db.First(client, "id = ?", client.Id).Error)
		require.Len(t, client.Secrets, 1)
		require.True(t, client.Secrets[0].Check(result.Clear))
	})
}

func TestDeleteSecret(t *testing.T) {
	withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {
		client := auth.NewClient(auth.ClientOptions{})
		secret, _ := client.GenerateNewSecret("testing")
		require.NoError(t, db.Create(client).Error)

		req := httptest.NewRequest(http.MethodDelete, "/clients/"+client.Id+"/secrets/"+secret.ID, nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusOK, res.Code)
		require.NoError(t, db.First(client, "id = ?", client.Id).Error)
		require.Len(t, client.Secrets, 0)
	})
}
