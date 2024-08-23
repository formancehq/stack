package clickhousetesting

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ory/dockertest/v3"
	"github.com/pkg/errors"

	"github.com/formancehq/stack/libs/go-libs/testing/docker"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type TestingT interface {
	require.TestingT
	Cleanup(func())
}

type Server struct {
	Port string
}

func (s *Server) GetHost() string {
	return "127.0.0.1"
}

func (s *Server) GetDSN() string {
	return s.GetDatabaseDSN("")
}

func (s *Server) GetDatabaseDSN(databaseName string) string {
	return fmt.Sprintf("clickhouse://%s:%s/%s", s.GetHost(), s.Port, databaseName)
}

func (s *Server) NewDatabase(t TestingT) *Database {

	options, err := clickhouse.ParseDSN(s.GetDSN())
	require.NoError(t, err)

	db, err := clickhouse.Open(options)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	databaseName := uuid.NewString()
	err = db.Exec(context.Background(), fmt.Sprintf(`CREATE DATABASE "%s"`, databaseName))
	require.NoError(t, err)

	if os.Getenv("NO_CLEANUP") != "true" {
		t.Cleanup(func() {
			err = db.Exec(context.Background(), fmt.Sprintf(`DROP DATABASE "%s"`, databaseName))
			require.NoError(t, err)
		})
	}

	return &Database{
		url: s.GetDatabaseDSN(databaseName),
	}
}

func CreateServer(pool *docker.Pool) *Server {
	resource := pool.Run(docker.Configuration{
		RunOptions: &dockertest.RunOptions{
			Repository: "clickhouse/clickhouse-server",
			Tag:        "head",
		},
		CheckFn: func(ctx context.Context, resource *dockertest.Resource) error {
			dsn := fmt.Sprintf("clickhouse://127.0.0.1:%s", resource.GetPort("9000/tcp"))
			options, _ := clickhouse.ParseDSN(dsn)

			db, err := clickhouse.Open(options)
			if err != nil {
				return errors.Wrap(err, "opening database")
			}
			defer func() {
				_ = db.Close()
			}()

			if err := db.Ping(context.Background()); err != nil {
				return errors.Wrap(err, "pinging database")
			}

			return nil
		},
	})

	return &Server{
		Port: resource.GetPort("9000/tcp"),
	}
}

type Database struct {
	url string
}

func (d *Database) ConnString() string {
	return d.url
}
