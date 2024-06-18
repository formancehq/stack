package extension

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	dockertest "github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/js/modules"
	"net/http"
	"os"
)

type LedgerConfiguration struct {
	Version     string `json:"version"`
	PostgresDSN string `json:"postgresDSN"`
	Network     string `json:"network"`
	TestID      string `json:"testID"`
}

type Ledger struct {
	URL string `json:"url"`
}

func (c *LedgerConfiguration) resolve() error {
	if c.TestID == "" {
		c.TestID = os.Getenv("TEST_ID")
		if c.TestID == "" {
			c.TestID = uuid.NewString()
		}
	}
	if c.Version == "" {
		c.Version = os.Getenv("LEDGER_VERSION")
		if c.Version == "" {
			c.Version = "latest"
		}
	}
	if c.PostgresDSN == "" {
		c.PostgresDSN = os.Getenv("POSTGRES_DSN")
		if c.PostgresDSN == "" {
			return errors.New("missing postgres dsn")
		}
	}
	if c.Network == "" {
		c.Network = os.Getenv("DOCKER_NETWORK")
	}
	return nil
}

type Extension struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
	logger   *logrus.Logger
	testID   string `json:"testID"`
}

func (c *Extension) StartLedger(configuration LedgerConfiguration) *Ledger {
	var err error
	if err := configuration.resolve(); err != nil {
		panic(err)
	}

	logger := c.logger.WithFields(map[string]interface{}{
		"version": configuration.Version,
		"testid":  configuration.TestID,
	})

	defer func() {
		if e := recover(); e != nil {
			if c.pool != nil && c.resource != nil {
				if err := c.pool.Purge(c.resource); err != nil {
					logger.Errorf("unable to clean docker resource: %s", err)
				}
			}
			panic(e)
		}
	}()

	logger.Infof("Connecting to docker server...")
	c.pool, err = dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	err = c.pool.Client.Ping()
	if err != nil {
		panic(err)
	}

	var networkID string
	if configuration.Network != "" {
		logger.Infof("Searching network named '%s'...", configuration.Network)
		network, err := c.pool.NetworksByName(configuration.Network)
		if err != nil {
			panic(err)
		}
		networkID = network[0].Network.ID
	}

	options := make([]func(config *docker.HostConfig), 0)
	options = append(options, func(config *docker.HostConfig) {
		config.AutoRemove = false
		config.CPUShares = 4
		config.Memory = 1024 * 1024 * 1024 * 8
	})

	logger.Infof("Starting ledger container...")
	c.resource, err = c.pool.RunWithOptions(&dockertest.RunOptions{
		Name:         fmt.Sprintf("ledger-%s", configuration.TestID),
		Repository:   "ghcr.io/formancehq/ledger",
		Tag:          configuration.Version,
		Env:          envVars(configuration),
		NetworkID:    networkID,
		ExposedPorts: []string{"3068/tcp"},
	}, options...)
	if err != nil {
		panic(err)
	}

	logger.Infof("Ledger container started with id : %s", c.resource.Container.ID)

	logger.Infof("Checking ledger is alive...")
	if err := c.pool.Retry(func() error {
		_, err := http.Get("http://localhost:" + c.resource.GetPort("3068/tcp"))
		return err
	}); err != nil {
		panic(err)
	}

	logger.Infof("Ledger properly started!")

	return &Ledger{
		URL: "http://localhost:" + c.resource.GetPort("3068/tcp"),
	}
}

func (c *Extension) StopLedger() {
	if os.Getenv("NO_CLEANUP") == "true" {
		return
	}
	c.logger.Infof("Shutting down ledger container...")
	if err := c.pool.Client.KillContainer(docker.KillContainerOptions{
		ID:     c.resource.Container.ID,
		Signal: 15,
	}); err != nil {
		panic(err)
	}

	exitCode, err := c.pool.Client.WaitContainer(c.resource.Container.ID)
	if err != nil {
		panic(err)
	}
	if exitCode != 0 {
		panic(fmt.Errorf("unexpected status code %d when stopping the ledger", exitCode))
	}

	if err := c.pool.Purge(c.resource); err != nil {
		panic(err)
	}

	c.logger.Infof("Ledger stopped!")
}

func envVars(configuration LedgerConfiguration) []string {
	return []string{
		"BIND=:3068",
		"POSTGRES_URI=" + configuration.PostgresDSN,
		"POSTGRES_MAX_OPEN_CONNS=50",
		"POSTGRES_MAX_IDLE_CONNS=50",
		"DEBUG=false",
	}
}

func init() {
	ext := &Extension{
		logger: logrus.New(),
	}
	modules.Register("k6/x/formancehq/benchmarks", ext)
}
