package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/formancehq/go-libs/bun/bundebug"

	"github.com/formancehq/go-libs/logging"

	"github.com/formancehq/go-libs/bun/bunconnect"
	"github.com/uptrace/bun"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/auth/pkg/storage/sqlstorage"
	"github.com/google/uuid"
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
	withDbAndUserRouter(t, func(router chi.Router, db *bun.DB) {
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
	withDbAndUserRouter(t, func(router chi.Router, db *bun.DB) {
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

func withDbAndUserRouter(t *testing.T, callback func(router chi.Router, db *bun.DB)) {
	t.Parallel()

	hooks := make([]bun.QueryHook, 0)
	if testing.Verbose() {
		hooks = append(hooks, bundebug.NewQueryHook())
	}

	pgDatabase := srv.NewDatabase(t)
	db, err := bunconnect.OpenSQLDB(logging.TestingContext(), bunconnect.ConnectionOptions{
		DatabaseSourceName: pgDatabase.ConnString(),
	}, hooks...)
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, sqlstorage.Migrate(context.Background(), db))

	router := chi.NewRouter()
	addUserRoutes(db, router)

	callback(router, db)
}
