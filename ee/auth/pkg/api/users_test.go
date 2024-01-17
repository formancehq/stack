package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/uptrace/bun"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/auth/pkg/storage/sqlstorage"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

var (
	user1 = &auth.User{
		ID:      uuid.NewString(),
		Subject: "alice",
		Email:   "alice@formance.com",
	}

	user2 = &auth.User{
		ID:      uuid.NewString(),
		Subject: "bob",
		Email:   "bob@formance.com",
	}
)

func TestListUsers(t *testing.T) {
	withDbAndUserRouter(t, func(router *mux.Router, db *bun.DB) {
		_, err := db.NewInsert().Model(user1).Exec(context.Background())
		require.NoError(t, err)

		_, err = db.NewInsert().Model(user2).Exec(context.Background())
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		require.Equal(t, http.StatusOK, res.Code)

		users := readTestResponse[[]auth.User](t, res)
		require.Len(t, users, 2)
	})
}

func TestReadUser(t *testing.T) {
	withDbAndUserRouter(t, func(router *mux.Router, db *bun.DB) {
		_, err := db.NewInsert().Model(user1).Exec(context.Background())
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/users/"+user1.ID, nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		require.Equal(t, http.StatusOK, res.Code)

		user := readTestResponse[auth.User](t, res)
		require.Equal(t, *user1, user)
	})
}

func withDbAndUserRouter(t *testing.T, callback func(router *mux.Router, db *bun.DB)) {
	t.Parallel()

	pgDatabase := pgtesting.NewPostgresDatabase(t)
	db, err := bunconnect.OpenSQLDB(bunconnect.ConnectionOptions{
		DatabaseSourceName: pgDatabase.ConnString(),
		Debug:              testing.Verbose(),
	})
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, sqlstorage.Migrate(context.Background(), db))

	router := mux.NewRouter()
	addUserRoutes(db, router)

	callback(router, db)
}
