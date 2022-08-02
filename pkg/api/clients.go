package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	auth "github.com/numary/auth/pkg"
	"github.com/numary/go-libs/sharedapi"
	_ "github.com/numary/go-libs/sharedapi"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/zitadel/oidc/pkg/oidc"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func validationError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(sharedapi.ErrorResponse{
		ErrorCode:    "VALIDATION",
		ErrorMessage: err.Error(),
	}); err != nil {
		sharedlogging.GetLogger(r.Context()).Info("Error validating request: %s", err)
	}
}

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(sharedapi.ErrorResponse{
		ErrorCode:    "INTERNAL",
		ErrorMessage: err.Error(),
	}); err != nil {
		trace.SpanFromContext(r.Context()).RecordError(err)
	}
}

func writeObject[T any](w http.ResponseWriter, r *http.Request, v T) {
	if err := json.NewEncoder(w).Encode(sharedapi.BaseResponse[T]{
		Data: &v,
	}); err != nil {
		trace.SpanFromContext(r.Context()).RecordError(err)
	}
}

type client struct {
	auth.ClientOptions
	ID string `json:"id"`
}

func deleteSecret(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &auth.Client{}
		if err := db.Find(client, "id = ?", mux.Vars(r)["clientId"]).Error; err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				w.WriteHeader(http.StatusNotFound)
			default:
				internalServerError(w, r, err)
			}
			return
		}

		if !client.DeleteSecret(mux.Vars(r)["secretId"]) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := db.Save(client).Error; err != nil {
			internalServerError(w, r, err)
			return
		}
	}
}

type secretCreate struct {
	Name string `json:"name"`
}

type secretCreateResult struct {
	ID         string `json:"id"`
	LastDigits string `json:"lastDigits"`
	Name       string `json:"name"`
	Clear      string `json:"clear"`
}

func createSecret(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &auth.Client{}
		if err := db.Find(client, "id = ?", mux.Vars(r)["clientId"]).Error; err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				w.WriteHeader(http.StatusNotFound)
			default:
				internalServerError(w, r, err)
			}
			return
		}

		sc := secretCreate{}
		if err := json.NewDecoder(r.Body).Decode(&sc); err != nil {
			validationError(w, r, err)
			return
		}

		secret, clear := client.GenerateNewSecret(sc.Name)

		if err := db.Save(client).Error; err != nil {
			internalServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusOK)

		writeObject(w, r, secretCreateResult{
			ID:         secret.ID,
			LastDigits: secret.LastDigits,
			Name:       secret.Name,
			Clear:      clear,
		})
	}
}

func readClient(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &auth.Client{}
		if err := db.Find(client, "id = ?", mux.Vars(r)["clientId"]).Error; err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				w.WriteHeader(http.StatusNotFound)
			default:
				internalServerError(w, r, err)
			}
			return
		}
		writeObject(w, r, client)
	}
}

func listClients(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clients := make([]auth.Client, 0)
		if err := db.Find(&clients).Error; err != nil {
			internalServerError(w, r, err)
			return
		}
		writeObject(w, r, clients)
	}
}

func updateClient(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		c := auth.Client{}
		if err := db.First(&c, "id = ?", mux.Vars(r)["clientId"]).Error; err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				w.WriteHeader(http.StatusNotFound)
			default:
				internalServerError(w, r, err)
			}
			return
		}

		opts := auth.ClientOptions{}
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			validationError(w, r, err)
			return
		}

		c.Update(opts)

		if err := db.Save(c).Error; err != nil {
			internalServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusOK)

		writeObject(w, r, client{
			ClientOptions: opts,
			ID:            c.Id,
		})
	}
}

func createClient(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts := auth.ClientOptions{}
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			validationError(w, r, err)
			return
		}

		grantTypes := []oidc.GrantType{
			oidc.GrantTypeCode,
			oidc.GrantTypeRefreshToken,
		}
		if !opts.Public {
			grantTypes = append(grantTypes, oidc.GrantTypeClientCredentials)
		}

		c := auth.NewClient(opts)
		if err := db.Create(c).Error; err != nil {
			internalServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Location", "./"+c.Id)

		writeObject(w, r, client{
			ClientOptions: opts,
			ID:            c.Id,
		})
	}
}
