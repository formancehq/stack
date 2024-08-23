package pgtesting

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/ory/dockertest/v3"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"

	"github.com/formancehq/stack/libs/go-libs/testing/docker"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type TestingT interface {
	require.TestingT
	Cleanup(func())
}

type Database struct {
	Url string
}

func (s *Database) ConnString() string {
	return s.Url
}

func (s *Database) ConnectionOptions() bunconnect.ConnectionOptions {
	return bunconnect.ConnectionOptions{
		DatabaseSourceName: s.ConnString(),
	}
}

type PostgresServer struct {
	Port   string
	Config Config
}

func (s *PostgresServer) GetPort() int {
	v, err := strconv.ParseInt(s.Port, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(v)
}

func (s *PostgresServer) GetHost() string {
	return "127.0.0.1"
}

func (s *PostgresServer) GetUsername() string {
	return s.Config.InitialUsername
}

func (s *PostgresServer) GetPassword() string {
	return s.Config.InitialUserPassword
}

func (s *PostgresServer) GetDSN() string {
	return s.GetDatabaseDSN(s.Config.InitialDatabaseName)
}

func (s *PostgresServer) GetDatabaseDSN(databaseName string) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", s.Config.InitialUsername,
		s.Config.InitialUserPassword, s.GetHost(), s.Port, databaseName)
}

func (s *PostgresServer) NewDatabase(t TestingT) *Database {
	db, err := sql.Open("postgres", s.GetDSN())
	require.NoError(t, err)
	defer func() {
		require.Nil(t, db.Close())
	}()

	databaseName := uuid.NewString()
	_, err = db.ExecContext(sharedlogging.TestingContext(), fmt.Sprintf(`CREATE DATABASE "%s"`, databaseName))
	require.NoError(t, err)

	if os.Getenv("NO_CLEANUP") != "true" {
		t.Cleanup(func() {
			db, err := sql.Open("postgres", s.GetDSN())
			require.NoError(t, err)
			defer func() {
				require.Nil(t, db.Close())
			}()

			_, err = db.ExecContext(sharedlogging.TestingContext(), fmt.Sprintf(`DROP DATABASE "%s"`, databaseName))
			if err != nil {
				panic(err)
			}
		})
	}

	return &Database{
		Url: s.GetDatabaseDSN(databaseName),
	}
}

type Config struct {
	InitialDatabaseName string
	InitialUserPassword string
	InitialUsername     string
	StatusCheckInterval time.Duration
	MaximumWaitingTime  time.Duration
}

func (c Config) validate() error {
	if c.StatusCheckInterval == 0 {
		return errors.New("status check interval must be greater than 0")
	}
	if c.InitialUsername == "" {
		return errors.New("initial username must be defined")
	}
	if c.InitialUserPassword == "" {
		return errors.New("initial user password must be defined")
	}
	if c.InitialDatabaseName == "" {
		return errors.New("initial database name must be defined")
	}
	return nil
}

type option func(opts *Config)

func WithInitialDatabaseName(name string) option {
	return func(opts *Config) {
		opts.InitialDatabaseName = name
	}
}

func WithInitialUser(username, pwd string) option {
	return func(opts *Config) {
		opts.InitialUserPassword = pwd
		opts.InitialUsername = username
	}
}

func WithStatusCheckInterval(d time.Duration) option {
	return func(opts *Config) {
		opts.StatusCheckInterval = d
	}
}

func WithMaximumWaitingTime(d time.Duration) option {
	return func(opts *Config) {
		opts.MaximumWaitingTime = d
	}
}

var defaultOptions = []option{
	WithStatusCheckInterval(200 * time.Millisecond),
	WithInitialUser("root", "root"),
	WithMaximumWaitingTime(time.Minute),
	WithInitialDatabaseName("formance"),
}

func CreatePostgresServer(t TestingT, pool *docker.Pool, opts ...option) *PostgresServer {
	cfg := Config{}
	for _, opt := range append(defaultOptions, opts...) {
		opt(&cfg)
	}

	require.NoError(t, cfg.validate())

	resource := pool.Run(docker.Configuration{
		RunOptions: &dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "15-alpine",
			Env: []string{
				fmt.Sprintf("POSTGRES_USER=%s", cfg.InitialUsername),
				fmt.Sprintf("POSTGRES_PASSWORD=%s", cfg.InitialUserPassword),
				fmt.Sprintf("POSTGRES_DB=%s", cfg.InitialDatabaseName),
			},
			Cmd: []string{
				"-c", "superuser-reserved-connections=0",
				"-c", "enable_partition_pruning=on",
				"-c", "enable_partitionwise_join=on",
				"-c", "enable_partitionwise_aggregate=on",
			},
		},
		CheckFn: func(ctx context.Context, resource *dockertest.Resource) error {
			dsn := fmt.Sprintf(
				"postgresql://%s:%s@127.0.0.1:%s/%s?sslmode=disable",
				cfg.InitialUsername,
				cfg.InitialUserPassword,
				resource.GetPort("5432/tcp"),
				cfg.InitialDatabaseName,
			)
			db, err := sql.Open("postgres", dsn)
			if err != nil {
				return errors.Wrap(err, "opening database")
			}
			defer func() {
				_ = db.Close()
			}()

			if err := db.Ping(); err != nil {
				return errors.Wrap(err, "pinging database")
			}

			return nil
		},
	})

	return &PostgresServer{
		Port:   resource.GetPort("5432/tcp"),
		Config: cfg,
	}
}
