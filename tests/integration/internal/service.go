package internal

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/ory/dockertest/v3"
	"io"
	"time"

	"github.com/egymgmbh/go-prefix-writer/prefixer"
	"github.com/formancehq/go-libs/httpserver"
	serviceutils "github.com/formancehq/go-libs/service"
	"github.com/spf13/cobra"
)

type service interface {
	load(context.Context, *Test) error
	unload(context.Context, *Test) error
}

type cobraCommandService struct {
	args        func(*Test) []string
	routingName string
	routingFunc func(path, method string) bool
	command     func() *cobra.Command
	name        string
	cancel      func()
	started     bool
	appContext  context.Context
	running     bool
}

func (c *cobraCommandService) load(ctx context.Context, t *Test) error {
	command := c.command()
	command.SetArgs(c.args(t))

	writer := prefixer.New(t.env.writer, func() string {
		return c.name + " | "
	})
	command.SetOut(writer)
	command.SetErr(writer)

	c.appContext = serviceutils.ContextWithLifecycle(ctx)
	c.appContext = httpserver.ContextWithServerInfo(c.appContext)
	c.appContext, c.cancel = context.WithCancel(c.appContext)
	errCh := make(chan error, 1)
	go func() {
		errCh <- command.ExecuteContext(c.appContext)
	}()
	select {
	case <-serviceutils.Ready(c.appContext):
		if c.routingFunc != nil {
			t.registerServiceToRoute(c.routingName, routing{
				port:        uint16(httpserver.Port(c.appContext)),
				routingFunc: c.routingFunc,
			})
		}
		t.registerServiceToRoute(c.name, routing{port: uint16(httpserver.Port(c.appContext))})
		c.running = true

		return nil
	case err := <-errCh:
		return err
	case <-time.After(30 * time.Second):
		return fmt.Errorf("timeout waiting for service '%s' to be properly started", c.name)
	}
}

func (c cobraCommandService) unload(context.Context, *Test) error {
	if !c.running {
		return nil
	}
	c.cancel()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-serviceutils.Stopped(c.appContext):
		return nil
	}
}

func (c *cobraCommandService) WithArgs(fn func(*Test) []string) *cobraCommandService {
	c.args = fn

	return c
}

func (c *cobraCommandService) WithRoutingFunc(name string, fn func(path, method string) bool) *cobraCommandService {
	c.routingName = name
	c.routingFunc = fn

	return c
}

var _ service = (*cobraCommandService)(nil)

func NewCommandService(name string, factory func() *cobra.Command) *cobraCommandService {
	return &cobraCommandService{
		command: factory,
		name:    name,
	}
}

type HealthCheck func(test *Test, container *dockertest.Resource) bool

type dockerContainerService struct {
	entrypoint  []string
	repository  string
	tag         string
	mounts      func(*Test) []string
	healthCheck HealthCheck
	env         func(*Test) []string
	resource    *dockertest.Resource
}

func (d *dockerContainerService) load(ctx context.Context, test *Test) error {

	var err error
	d.resource, err = test.env.dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: d.repository,
		Tag:        d.tag,
		Mounts:     d.mounts(test),
		Tty:        true,
		Entrypoint: d.entrypoint,
		Env:        d.env(test),
	})
	if err != nil {
		return err
	}

	go func() {
		reader, _ := test.env.dockerClient.ContainerLogs(TestContext(), d.resource.Container.ID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
			Details:    false,
		})
		if reader != nil {
			_, _ = io.Copy(prefixer.New(test.env.writer, func() string {
				return "benthos | "
			}), reader)
		}
	}()

	healthCheckContext, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if d.healthCheck != nil {
		for {
			if d.healthCheck(test, d.resource) {
				break
			}
			select {
			case <-healthCheckContext.Done():
				return healthCheckContext.Err()
			case <-time.After(100 * time.Millisecond):
			}
		}
	}

	return nil
}

func (d *dockerContainerService) unload(ctx context.Context, test *Test) error {
	if d.resource == nil {
		return nil
	}
	return test.env.dockerPool.Purge(d.resource)
}

func (d *dockerContainerService) WithEntrypoint(entrypoint []string) *dockerContainerService {
	d.entrypoint = entrypoint

	return d
}

func (d *dockerContainerService) WithEnv(f func(*Test) []string) *dockerContainerService {
	d.env = f

	return d
}

func (d *dockerContainerService) WithMounts(mounts func(test *Test) []string) *dockerContainerService {
	d.mounts = mounts

	return d
}

func (d *dockerContainerService) WithHealthCheck(fn HealthCheck) *dockerContainerService {
	d.healthCheck = fn

	return d
}

var _ service = (*dockerContainerService)(nil)

func NewDockerService(repository, tag string) *dockerContainerService {
	return &dockerContainerService{
		repository: repository,
		tag:        tag,
	}
}
