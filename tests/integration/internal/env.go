package internal

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ory/dockertest/v3"
	"io"
	"os"
)

type env struct {
	sqlConn      *pgx.Conn
	writer       io.Writer
	dockerPool   *dockertest.Pool
	dockerClient *client.Client
	workdir      string
}

func (e *env) Setup(ctx context.Context) error {

	var err error
	e.workdir, err = os.Getwd()
	if err != nil {
		return err
	}

	e.sqlConn, err = pgx.Connect(ctx, GetPostgresDSNString())
	if err != nil {
		return err
	}

	e.dockerPool, err = dockertest.NewPool("")
	if err != nil {
		return err
	}

	e.dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	// uses pool to try to connect to Docker
	err = e.dockerPool.Client.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (e *env) Teardown(ctx context.Context) error {
	if err := e.dockerClient.Close(); err != nil {
		return err
	}
	// TODO: Purge docker pool resources
	if e.sqlConn != nil {
		return e.sqlConn.Close(ctx)
	}
	return nil
}

func (e *env) newTest() *Test {
	return &Test{
		env:             e,
		id:              uuid.NewString(),
		loadedModules:   collectionutils.NewSet[string](),
		servicesToRoute: make(map[string]uint16),
	}
}

func newEnv(w io.Writer) *env {
	return &env{
		writer: w,
	}
}
