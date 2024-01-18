package api

import (
	"net/http"

	"github.com/uptrace/bun"

	auth "github.com/formancehq/auth/pkg"
	"github.com/gorilla/mux"
)

func addUserRoutes(db *bun.DB, router *mux.Router) {
	router.Path("/users").Methods(http.MethodGet).HandlerFunc(listUsers(db))
	router.Path("/users/{userId}").Methods(http.MethodGet).HandlerFunc(readUser(db))
}

func listUsers(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := make([]auth.User, 0)
		if err := db.
			NewSelect().
			Model(&users).
			Scan(r.Context()); err != nil {
			internalServerError(w, r, err)
			return
		}
		writeJSONObject(w, r, users)
	}
}

func readUser(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := findById[*auth.User](w, r, db, "userId")
		if user == nil {
			return
		}
		writeJSONObject(w, r, user)
	}
}
