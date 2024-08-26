package pgtesting

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	dockerlib "github.com/ory/dockertest/v3/docker"

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
	url string
}

func (s *Database) ConnString() string {
	return s.url
}

func (s *Database) ConnectionOptions() bunconnect.ConnectionOptions {
	return bunconnect.ConnectionOptions{
		DatabaseSourceName: s.ConnString(),
	}
}

type PostgresServer struct {
	port   string
	config config
	t      TestingT
}

func (s *PostgresServer) GetPort() int {
	v, err := strconv.ParseInt(s.port, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(v)
}

func (s *PostgresServer) GetHost() string {
	return "127.0.0.1"
}

func (s *PostgresServer) GetUsername() string {
	return s.config.initialUsername
}

func (s *PostgresServer) GetPassword() string {
	return s.config.initialUserPassword
}

func (s *PostgresServer) GetDSN() string {
	return s.GetDatabaseDSN(s.config.initialDatabaseName)
}

func (s *PostgresServer) GetDatabaseDSN(databaseName string) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", s.config.initialUsername,
		s.config.initialUserPassword, s.GetHost(), s.port, databaseName)
}

func (s *PostgresServer) NewDatabase() *Database {
	db, err := sql.Open("postgres", s.GetDSN())
	require.NoError(s.t, err)
	defer func() {
		require.Nil(s.t, db.Close())
	}()

	databaseName := uuid.NewString()
	_, err = db.ExecContext(sharedlogging.TestingContext(), fmt.Sprintf(`CREATE DATABASE "%s"`, databaseName))
	require.NoError(s.t, err)

	if os.Getenv("NO_CLEANUP") != "true" {
		s.t.Cleanup(func() {
			db, err := sql.Open("postgres", s.GetDSN())
			require.NoError(s.t, err)
			defer func() {
				require.Nil(s.t, db.Close())
			}()

			_, err = db.ExecContext(sharedlogging.TestingContext(), fmt.Sprintf(`DROP DATABASE "%s"`, databaseName))
			if err != nil {
				panic(err)
			}
		})
	}

	return &Database{
		url: s.GetDatabaseDSN(databaseName),
	}
}

type config struct {
	initialDatabaseName string
	initialUserPassword string
	initialUsername     string
	statusCheckInterval time.Duration
	maximumWaitingTime  time.Duration
	hostConfigOptions   []func(hostConfig *dockerlib.HostConfig)
	logger              sharedlogging.Logger
}

func (c config) validate() error {
	if c.statusCheckInterval == 0 {
		return errors.New("status check interval must be greater than 0")
	}
	if c.initialUsername == "" {
		return errors.New("initial username must be defined")
	}
	if c.initialUserPassword == "" {
		return errors.New("initial user password must be defined")
	}
	if c.initialDatabaseName == "" {
		return errors.New("initial database name must be defined")
	}
	return nil
}

type option func(opts *config)

func WithInitialDatabaseName(name string) option {
	return func(opts *config) {
		opts.initialDatabaseName = name
	}
}

func WithInitialUser(username, pwd string) option {
	return func(opts *config) {
		opts.initialUserPassword = pwd
		opts.initialUsername = username
	}
}

func WithStatusCheckInterval(d time.Duration) option {
	return func(opts *config) {
		opts.statusCheckInterval = d
	}
}

func WithMaximumWaitingTime(d time.Duration) option {
	return func(opts *config) {
		opts.maximumWaitingTime = d
	}
}

func WithDockerHostConfigOption(opt func(hostConfig *dockerlib.HostConfig)) option {
	return func(opts *config) {
		opts.hostConfigOptions = append(opts.hostConfigOptions, opt)
	}
}

func WithLogger(logger sharedlogging.Logger) option {
	return func(opts *config) {
		opts.logger = logger
	}
}

var defaultOptions = []option{
	WithStatusCheckInterval(200 * time.Millisecond),
	WithInitialUser("root", "root"),
	WithMaximumWaitingTime(time.Minute),
	WithInitialDatabaseName("formance"),
	WithLogger(sharedlogging.NewDefaultLogger(os.Stdout, false, false)),
}

func CreatePostgresServer(t TestingT, pool *docker.Pool, opts ...option) *PostgresServer {
	cfg := config{}
	for _, opt := range append(defaultOptions, opts...) {
		opt(&cfg)
	}

	require.NoError(t, cfg.validate())

	resource := pool.Run(docker.Configuration{
		RunOptions: &dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "15-alpine",
			Env: []string{
				fmt.Sprintf("POSTGRES_USER=%s", cfg.initialUsername),
				fmt.Sprintf("POSTGRES_PASSWORD=%s", cfg.initialUserPassword),
				fmt.Sprintf("POSTGRES_DB=%s", cfg.initialDatabaseName),
			},
			Cmd: []string{
				"-c", "superuser-reserved-connections=0",
				"-c", "enable_partition_pruning=on",
				"-c", "enable_partitionwise_join=on",
				"-c", "enable_partitionwise_aggregate=on",
			},
		},
		HostConfigOptions: cfg.hostConfigOptions,
		CheckFn: func(ctx context.Context, resource *dockertest.Resource) error {
			dsn := fmt.Sprintf(
				"postgresql://%s:%s@127.0.0.1:%s/%s?sslmode=disable",
				cfg.initialUsername,
				cfg.initialUserPassword,
				resource.GetPort("5432/tcp"),
				cfg.initialDatabaseName,
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
		port:   resource.GetPort("5432/tcp"),
		config: cfg,
		t:      t,
	}
}
