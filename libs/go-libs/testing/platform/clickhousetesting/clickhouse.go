package clickhousetesting

import (
	"context"
	"fmt"

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
	port string
}

func (s *Server) GetHost() string {
	return "127.0.0.1"
}

func (s *Server) GetDSN() string {
	return s.GetDatabaseDSN("")
}

func (s *Server) GetDatabaseDSN(databaseName string) string {
	return fmt.Sprintf("clickhouse://%s:%s/%s", s.GetHost(), s.port, databaseName)
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
		port: resource.GetPort("9000/tcp"),
	}
}
