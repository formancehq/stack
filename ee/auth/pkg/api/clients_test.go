package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/auth/pkg/storage/sqlstorage"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func withDbAndClientRouter(t *testing.T, callback func(router *mux.Router, db *gorm.DB)) {
	t.Parallel()

	pgDatabase := pgtesting.NewPostgresDatabase(t)
	dialector := postgres.Open(pgDatabase.ConnString())

	db, err := sqlstorage.LoadGorm(dialector, &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	defer sqlDB.Close()

	require.NoError(t, sqlstorage.MigrateTables(context.Background(), db))
	require.NoError(t, sqlstorage.MigrateData(context.Background(), db))

	router := mux.NewRouter()
	addClientRoutes(db, router)

	callback(router, db)
}

func TestCreateClient(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name    string
		options auth.ClientOptions
	}
	for _, tc := range []testCase{
		{
			name: "confidential client",
			options: auth.ClientOptions{
				Name:                   "confidential client",
				RedirectURIs:           []string{"http://localhost:8080"},
				Description:            "abc",
				PostLogoutRedirectUris: []string{"http://localhost:8080/logout"},
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name: "public client",
			options: auth.ClientOptions{
				Name:   "public client",
				Public: true,
			},
		},
		{
			name: "confidential client",
			options: auth.ClientOptions{
				Name:                   "confidential client",
				RedirectURIs:           []string{"http://localhost:8080"},
				Description:            "abc",
				PostLogoutRedirectUris: []string{"http://localhost:8080/logout"},
				Metadata: map[string]string{
					"foo": "bar",
				},
				Scopes: []string{"ledger:read", "ledger:write", "formance:test"},
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {
				req := httptest.NewRequest(http.MethodPost, "/clients", createJSONBuffer(t, tc.options))
				res := httptest.NewRecorder()

				router.ServeHTTP(res, req)

				require.Equal(t, http.StatusCreated, res.Code)

				createdClient := readTestResponse[clientView](t, res)
				require.NotEmpty(t, createdClient.ID)
				tcScopes := tc.options.Scopes
				tc.options.Scopes = nil
				require.Equal(t, tc.options, createdClient.ClientOptions)
				require.True(t, func() bool {
					for _, scope := range tcScopes {
						contain := collectionutils.Contains[string](createdClient.Scopes, scope)
						if !contain {
							t.Logf("scope %s not found in created client scopes", scope)
							return false
						}
					}

					return true
				}())
				tc.options.Id = createdClient.ID
				tc.options.Scopes = tcScopes
				clientFromDatabase := auth.Client{}
				require.NoError(t, db.Find(&clientFromDatabase, "id = ?", createdClient.ID).Error)
				require.Equal(t, auth.Client{
					ClientOptions: tc.options,
				}, clientFromDatabase)
			})
		})

	}
}

func TestUpdateClient(t *testing.T) {

	t.Parallel()
	type testCase struct {
		name    string
		options auth.ClientOptions
	}
	for _, tc := range []testCase{
		{
			name: "confidential client",
			options: auth.ClientOptions{
				Name:                   "confidential client",
				RedirectURIs:           []string{"http://localhost:8080"},
				Description:            "abc",
				PostLogoutRedirectUris: []string{"http://localhost:8080/logout"},
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name: "public client",
			options: auth.ClientOptions{
				Name:   "public client",
				Public: true,
			},
		},
		{
			name: "confidential client",
			options: auth.ClientOptions{
				Name:                   "confidential client",
				RedirectURIs:           []string{"http://localhost:8080"},
				Description:            "abc",
				PostLogoutRedirectUris: []string{"http://localhost:8080/logout"},
				Metadata: map[string]string{
					"foo": "bar",
				},
				Scopes: []string{"ledger:read", "ledger:write", "formance:test"},
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {

				initialClient := auth.NewClient(auth.ClientOptions{})
				require.NoError(t, db.Create(initialClient).Error)

				req := httptest.NewRequest(http.MethodPut, "/clients/"+initialClient.Id, createJSONBuffer(t, tc.options))
				res := httptest.NewRecorder()

				router.ServeHTTP(res, req)

				require.Equal(t, http.StatusOK, res.Code)

				updatedClient := readTestResponse[clientView](t, res)
				tcScopes := tc.options.Scopes
				tc.options.Scopes = nil
				require.NotEmpty(t, updatedClient.ID)
				require.Equal(t, tc.options, updatedClient.ClientOptions)

				tc.options.Id = updatedClient.ID
				tc.options.Scopes = tcScopes
				clientFromDatabase := auth.Client{}
				require.NoError(t, db.Find(&clientFromDatabase, "id = ?", updatedClient.ID).Error)
				require.Equal(t, auth.Client{
					ClientOptions: tc.options,
				}, clientFromDatabase)
			})
		})
	}
}

func TestListClients(t *testing.T) {
	withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {
		client1 := auth.NewClient(auth.ClientOptions{})
		require.NoError(t, db.Create(client1).Error)

		client2 := auth.NewClient(auth.ClientOptions{
			Metadata: map[string]string{
				"foo": "bar",
			},
		})
		require.NoError(t, db.Create(client2).Error)

		req := httptest.NewRequest(http.MethodGet, "/clients", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusOK, res.Code)

		clients := readTestResponse[[]clientView](t, res)
		require.Len(t, clients, 2)
		require.Len(t, clients[1].Metadata, 1)
		require.Equal(t, clients[1].Metadata["foo"], "bar")
	})
}

func TestReadClient(t *testing.T) {
	withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {

		opts := auth.ClientOptions{
			Metadata: map[string]string{
				"foo": "bar",
			},
		}
		client1 := auth.NewClient(opts)
		client1.Scopes = append(client1.Scopes, "XXX")
		secret, _ := client1.GenerateNewSecret(auth.SecretCreate{
			Name: "testing",
		})
		require.NoError(t, db.Create(client1).Error)

		req := httptest.NewRequest(http.MethodGet, "/clients/"+client1.Id, nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusOK, res.Code)

		ret := readTestResponse[clientView](t, res)
		require.Equal(t, clientView{
			ClientOptions: opts,
			ID:            client1.Id,
			Scopes:        []string{"XXX"},
			Secrets: []clientSecretView{{
				ClientSecret: secret,
			}},
		}, ret)
	})
}

func TestDeleteClient(t *testing.T) {
	withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {

		opts := auth.ClientOptions{
			Metadata: map[string]string{
				"foo": "bar",
			},
		}
		client1 := auth.NewClient(opts)
		client1.Scopes = append(client1.Scopes, "XXX")
		require.NoError(t, db.Create(client1).Error)

		req := httptest.NewRequest(http.MethodDelete, "/clients/"+client1.Id, nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusNoContent, res.Code)
	})
}

func TestGenerateNewSecret(t *testing.T) {
	withDbAndClientRouter(t, func(router *mux.Router, db *gorm.DB) {
		client := auth.NewClient(auth.ClientOptions{})
		require.NoError(t, db.Create(client).Error)

		req := httptest.NewRequest(http.MethodPost, "/clients/"+client.Id+"/secrets", createJSONBuffer(t, auth.SecretCreate{
			Name: "secret1",
		}))
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		result := readTestResponse[secretCreateResult](t, res)
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
		secret, _ := client.GenerateNewSecret(auth.SecretCreate{
			Name: "testing",
		})
		require.NoError(t, db.Create(client).Error)

		req := httptest.NewRequest(http.MethodDelete, "/clients/"+client.Id+"/secrets/"+secret.ID, nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusNoContent, res.Code)
		require.NoError(t, db.First(client, "id = ?", client.Id).Error)
		require.Len(t, client.Secrets, 0)
	})
}
