package api

import (
	"net/http"

	"github.com/uptrace/bun"

	auth "github.com/formancehq/auth/pkg"
	_ "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/gorilla/mux"
)

func addClientRoutes(db *bun.DB, router *mux.Router) {
	router.Path("/clients").Methods(http.MethodPost).HandlerFunc(createClient(db))
	router.Path("/clients").Methods(http.MethodGet).HandlerFunc(listClients(db))
	router.Path("/clients/{clientId}").Methods(http.MethodPut).HandlerFunc(updateClient(db))
	router.Path("/clients/{clientId}").Methods(http.MethodGet).HandlerFunc(readClient(db))
	router.Path("/clients/{clientId}").Methods(http.MethodDelete).HandlerFunc(deleteClient(db))
	router.Path("/clients/{clientId}/secrets").Methods(http.MethodPost).HandlerFunc(createSecret(db))
	router.Path("/clients/{clientId}/secrets/{secretId}").Methods(http.MethodDelete).HandlerFunc(deleteSecret(db))
}

type clientSecretView struct {
	auth.ClientSecret
	Hash string `json:"-"`
}

type clientView struct {
	auth.ClientOptions
	ID      string                       `json:"id"`
	Secrets auth.Array[clientSecretView] `json:"secrets" bun:"type:text"`
}

func mapBusinessClient(c auth.Client) clientView {
	return clientView{
		ClientOptions: auth.ClientOptions{
			Public:                 c.Public,
			RedirectURIs:           c.RedirectURIs,
			Description:            c.Description,
			Name:                   c.Name,
			PostLogoutRedirectUris: c.PostLogoutRedirectUris,
			Metadata:               c.Metadata,
			Scopes:                 c.Scopes,
		},
		ID: c.Id,
		Secrets: mapList(c.Secrets, func(i auth.ClientSecret) clientSecretView {
			return clientSecretView{
				ClientSecret: i,
			}
		}),
	}
}

type secretCreateResult struct {
	ID         string `json:"id"`
	LastDigits string `json:"lastDigits"`
	Name       string `json:"name"`
	Clear      string `json:"clear"`
}

func deleteSecret(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := findById[*auth.Client](w, r, db, "clientId")
		if client == nil {
			return
		}

		if !client.DeleteSecret(mux.Vars(r)["secretId"]) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, err := db.NewUpdate().
			Model(client).
			Where("id = ?", client.Id).
			Exec(r.Context())
		if err != nil {
			internalServerError(w, r, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func createSecret(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := findById[*auth.Client](w, r, db, "clientId")
		if client == nil {
			return
		}

		sc := readJSONObject[auth.SecretCreate](w, r)
		if sc == nil {
			return
		}

		secret, clear := client.GenerateNewSecret(*sc)

		_, err := db.NewUpdate().
			Model(client).
			Where("id = ?", client.Id).
			Exec(r.Context())
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		writeJSONObject(w, r, secretCreateResult{
			ID:         secret.ID,
			LastDigits: secret.LastDigits,
			Name:       secret.Name,
			Clear:      clear,
		})
	}
}

func readClient(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := findById[*auth.Client](w, r, db, "clientId")
		if client == nil {
			return
		}
		writeJSONObject(w, r, mapBusinessClient(*client))
	}
}

func deleteClient(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := db.
			NewDelete().
			Model(&auth.Client{}).
			Where("id = ?", mux.Vars(r)["clientId"]).
			Exec(r.Context())
		if err != nil {
			internalServerError(w, r, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func listClients(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clients := make([]auth.Client, 0)
		if err := db.
			NewSelect().
			Model(&clients).
			Scan(r.Context()); err != nil {
			internalServerError(w, r, err)
			return
		}
		writeJSONObject(w, r, mapList(clients, mapBusinessClient))
	}
}

func updateClient(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		client := findById[*auth.Client](w, r, db, "clientId")
		if client == nil {
			return
		}

		opts := readJSONObject[auth.ClientOptions](w, r)
		if opts == nil {
			return
		}

		client.Update(*opts)

		_, err := db.NewUpdate().
			Model(client).
			Where("id = ?", client.Id).
			Exec(r.Context())
		if err != nil {
			internalServerError(w, r, err)
			return
		}

		writeJSONObject(w, r, mapBusinessClient(*client))
	}
}

func createClient(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts := readJSONObject[auth.ClientOptions](w, r)
		if opts == nil {
			return
		}

		c := auth.NewClient(*opts)
		if err := createObject(w, r, db, c); err != nil {
			return
		}

		writeCreatedJSONObject(w, r, mapBusinessClient(*c), c.Id)
	}
}
