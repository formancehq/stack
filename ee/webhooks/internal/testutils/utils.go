package tests

import (
	"fmt"
	_ "fmt"
	_ "math/rand"
	"net/http"
	"time"

	"github.com/formancehq/webhooks/internal/migrations"
	"github.com/formancehq/webhooks/internal/services/httpclient"
	storage "github.com/formancehq/webhooks/internal/services/storage/postgres"

	"os"

	"database/sql"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"

	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/uptrace/bun"
)

func StartPostgresServer() {
	if err := pgtesting.CreatePostgresServer(); err != nil {
		logging.Error(err)
		os.Exit(1)
	}
}

func StopPostgresServer() {
	if err := pgtesting.DestroyPostgresServer(); err != nil {
		logging.Error(err)
		os.Exit(1)
	}
}

func GetStoreProvider() (storage.PostgresStore, error) {

	postgres := storage.NewPostgresStoreProvider(nil)
	db, err := sql.Open("postgres", pgtesting.Server().GetDSN())
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}

	bunDB := bun.NewDB(db, pgdialect.New())

	err = migrations.Migrate(logging.TestingContext(), bunDB)

	if err != nil {
		return postgres, err
	}

	return storage.NewPostgresStoreProvider(bunDB), nil
}

func NewHTTPServer(port int, routes ...[2]interface{}) *http.Server {
	mux := http.NewServeMux()

	for _, route := range routes {
		url, ok1 := route[0].(string)
		handler, ok2 := route[1].(http.HandlerFunc)
		if !ok1 || !ok2 {
			logging.Errorf("Invalid route format. Expected (string, http.HandlerFunc), got (%T, %T)", route[0], route[1])
		}
		mux.HandleFunc(url, handler)
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Errorf("Could not listen on port %d: %v", port, err)
			os.Exit(1)
		}
		logging.Debugf("TestServer is running on %s", server.Addr)
	}()

	return server
}

func NewHTTPClient() *httpclient.DefaultHttpClient {
	client := httpclient.NewDefaultHttpClient(&http.Client{})
	return &client
}
