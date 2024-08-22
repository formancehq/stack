package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/bun/bundebug"

	"github.com/go-chi/chi/v5"

	"github.com/formancehq/stack/libs/go-libs/logging"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/uptrace/bun"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/auth/pkg/storage/sqlstorage"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/stretchr/testify/require"
)

func withDbAndClientRouter(t *testing.T, callback func(router chi.Router, db *bun.DB)) {
	t.Parallel()

	pgDatabase := srv.NewDatabase()

	hooks := make([]bun.QueryHook, 0)
	if testing.Verbose() {
		hooks = append(hooks, bundebug.NewQueryHook())
	}

	db, err := bunconnect.OpenSQLDB(logging.TestingContext(), bunconnect.ConnectionOptions{
		DatabaseSourceName: pgDatabase.ConnString(),
	}, hooks...)
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, sqlstorage.Migrate(context.Background(), db))

	router := chi.NewRouter()
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
				Scopes: []string{},
			},
		},
		{
			name: "public client",
			options: auth.ClientOptions{
				Name:                   "public client",
				Public:                 true,
				RedirectURIs:           []string{},
				PostLogoutRedirectUris: []string{},
				Metadata:               map[string]string{},
				Scopes:                 []string{},
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

			withDbAndClientRouter(t, func(router chi.Router, db *bun.DB) {
				req := httptest.NewRequest(http.MethodPost, "/clients", createJSONBuffer(t, tc.options))
				res := httptest.NewRecorder()

				router.ServeHTTP(res, req)

				require.Equal(t, http.StatusCreated, res.Code)

				createdClient := readTestResponse[clientView](t, res)
				require.NotEmpty(t, createdClient.ID)
				require.Equal(t, tc.options, createdClient.ClientOptions)
				require.True(t, func() bool {
					for _, scope := range tc.options.Scopes {
						contain := collectionutils.Contains[string](createdClient.Scopes, scope)
						if !contain {
							t.Logf("scope %s not found in created client scopes", scope)
							return false
						}
					}

					return true
				}())
				tc.options.Id = createdClient.ID
				clientFromDatabase := auth.Client{}
				err := db.NewSelect().
					Model(&clientFromDatabase).
					Where("id = ?", createdClient.ID).
					Scan(context.Background())
				require.NoError(t, err)
				require.Equal(t, auth.Client{
					ClientOptions: tc.options,
					Secrets:       []auth.ClientSecret{},
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
				Scopes: []string{},
			},
		},
		{
			name: "public client",
			options: auth.ClientOptions{
				Name:                   "public client",
				Public:                 true,
				RedirectURIs:           []string{},
				PostLogoutRedirectUris: []string{},
				Metadata:               map[string]string{},
				Scopes:                 []string{},
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
			withDbAndClientRouter(t, func(router chi.Router, db *bun.DB) {

				initialClient := auth.NewClient(auth.ClientOptions{})
				_, err := db.NewInsert().Model(initialClient).Exec(context.Background())
				require.NoError(t, err)

				req := httptest.NewRequest(http.MethodPut, "/clients/"+initialClient.Id, createJSONBuffer(t, tc.options))
				res := httptest.NewRecorder()

				router.ServeHTTP(res, req)

				require.Equal(t, http.StatusOK, res.Code)

				updatedClient := readTestResponse[clientView](t, res)
				require.NotEmpty(t, updatedClient.ID)
				require.Equal(t, tc.options, updatedClient.ClientOptions)

				tc.options.Id = updatedClient.ID
				clientFromDatabase := auth.Client{}
				err = db.NewSelect().
					Model(&clientFromDatabase).
					Where("id = ?", updatedClient.ID).
					Scan(context.Background())
				require.Equal(t, auth.Client{
					ClientOptions: tc.options,
					Secrets:       []auth.ClientSecret{},
				}, clientFromDatabase)
			})
		})
	}
}

func TestListClients(t *testing.T) {
	withDbAndClientRouter(t, func(router chi.Router, db *bun.DB) {
		client1 := auth.NewClient(auth.ClientOptions{})
		_, err := db.NewInsert().Model(client1).Exec(context.Background())
		require.NoError(t, err)

		client2 := auth.NewClient(auth.ClientOptions{
			Metadata: map[string]string{
				"foo": "bar",
			},
		})
		_, err = db.NewInsert().Model(client2).Exec(context.Background())
		require.NoError(t, err)

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
	withDbAndClientRouter(t, func(router chi.Router, db *bun.DB) {

		opts := auth.ClientOptions{
			Metadata: map[string]string{
				"foo": "bar",
			},
			Scopes:                 []string{"XXX"},
			RedirectURIs:           []string{},
			PostLogoutRedirectUris: []string{},
		}
		client1 := auth.NewClient(opts)
		secret, _ := client1.GenerateNewSecret(auth.SecretCreate{
			Name: "testing",
		})
		_, err := db.NewInsert().Model(client1).Exec(context.Background())
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/clients/"+client1.Id, nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusOK, res.Code)

		ret := readTestResponse[clientView](t, res)
		require.Equal(t, clientView{
			ClientOptions: opts,
			ID:            client1.Id,
			Secrets: []clientSecretView{{
				ClientSecret: secret,
			}},
		}, ret)
	})
}

func TestDeleteClient(t *testing.T) {
	withDbAndClientRouter(t, func(router chi.Router, db *bun.DB) {

		opts := auth.ClientOptions{
			Metadata: map[string]string{
				"foo": "bar",
			},
		}
		client1 := auth.NewClient(opts)
		client1.Scopes = append(client1.Scopes, "XXX")
		_, err := db.NewInsert().Model(client1).Exec(context.Background())
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodDelete, "/clients/"+client1.Id, nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusNoContent, res.Code)
	})
}

func TestGenerateNewSecret(t *testing.T) {
	withDbAndClientRouter(t, func(router chi.Router, db *bun.DB) {
		client := auth.NewClient(auth.ClientOptions{})
		_, err := db.NewInsert().Model(client).Exec(context.Background())
		require.NoError(t, err)

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

		err = db.NewSelect().
			Model(client).
			Limit(1).
			Where("id = ?", client.Id).
			Scan(context.Background())
		require.NoError(t, err)
		require.Len(t, client.Secrets, 1)
		require.True(t, client.Secrets[0].Check(result.Clear))
	})
}

func TestDeleteSecret(t *testing.T) {
	withDbAndClientRouter(t, func(router chi.Router, db *bun.DB) {
		client := auth.NewClient(auth.ClientOptions{})
		secret, _ := client.GenerateNewSecret(auth.SecretCreate{
			Name: "testing",
		})
		_, err := db.NewInsert().Model(client).Exec(context.Background())
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodDelete, "/clients/"+client.Id+"/secrets/"+secret.ID, nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		require.Equal(t, http.StatusNoContent, res.Code)

		err = db.NewSelect().
			Model(client).
			Where("id = ?", client.Id).
			Scan(context.Background())
		require.NoError(t, err)
		require.Len(t, client.Secrets, 0)
	})
}
