package extension

import (
	"errors"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/js/modules"
	"golang.org/x/mod/semver"
)

type OTLP struct {
	Endpoint string `json:"endpoint"`
}

type LedgerConfiguration struct {
	OTLP        OTLP   `json:"otlp"`
	Version     string `json:"version"`
	PostgresDSN string `json:"postgresDSN"`
	Network     string `json:"network"`
}

type Ledger struct {
	URL string `json:"url"`
}

func (c *LedgerConfiguration) resolve() error {
	if c.Version == "" {
		c.Version = os.Getenv("LEDGER_VERSION")
		if c.Version == "" {
			return errors.New("missing ledger version")
		}
	}
	if c.PostgresDSN == "" {
		c.PostgresDSN = os.Getenv("POSTGRES_DSN")
		if c.PostgresDSN == "" {
			return errors.New("missing postgres dsn")
		}
	}
	if c.OTLP.Endpoint == "" {
		c.OTLP.Endpoint = os.Getenv("OTLP_ENDPOINT")
		if c.OTLP.Endpoint == "" {
			return errors.New("missing otlp endpoint")
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
	runID    string `json:"runID"`
}

func (c *Extension) StartLedger(configuration LedgerConfiguration) *Ledger {

	c.runID = uuid.New().String()
	logger := c.logger.WithFields(map[string]interface{}{
		"configuration": configuration,
	})

	defer func() {
		if e := recover(); e != nil {
			if c.pool != nil {
				if err := c.pool.Purge(c.resource); err != nil {
					logger.Errorf("unable to clean docker resource: %s", err)
				}
			}
			panic(e)
		}
	}()
	var err error

	if err := configuration.resolve(); err != nil {
		panic(err)
	}

	logger.Infof("Connecting to docker server...")
	c.pool, err = dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	err = c.pool.Client.Ping()
	if err != nil {
		panic(err)
	}

	var envVars []string
	if semver.IsValid(configuration.Version) && semver.Compare(configuration.Version, "v2.0.0") < 0 {
		envVars = v1EnvVars(c.runID, configuration)
	} else {
		envVars = v2EnvVars(c.runID, configuration)
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

	logger.Infof("Starting ledger container...")
	c.resource, err = c.pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "ghcr.io/formancehq/ledger",
		Tag:          configuration.Version,
		Env:          envVars,
		NetworkID:    networkID,
		ExposedPorts: []string{"3068/tcp"},
	})
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
	c.logger.Infof("Shutting down ledger container...")
	if err := c.pool.Client.KillContainer(docker.KillContainerOptions{
		ID:     c.resource.Container.ID,
		Signal: 15,
	}); err != nil {
		panic(err)
	}
	if err := c.pool.Purge(c.resource); err != nil {
		panic(err)
	}

	c.logger.Infof("Ledger stopped!")
}

func v1EnvVars(runID string, configuration LedgerConfiguration) []string {
	return []string{
		"NUMARY_SERVER_HTTP_BIND_ADDRESS=:3068",
		"NUMARY_STORAGE_DRIVER=postgres",
		"NUMARY_STORAGE_POSTGRES_CONN_STRING=" + configuration.PostgresDSN,
		"NUMARY_OTEL_TRACES=true",
		"NUMARY_OTEL_TRACES_EXPORTER=otlp",
		"NUMARY_OTEL_TRACES_EXPORTER_OTLP_ENDPOINT=" + configuration.OTLP.Endpoint,
		"NUMARY_OTEL_TRACES_EXPORTER_OTLP_INSECURE=true",
		"NUMARY_OTEL_METRICS=true",
		"NUMARY_OTEL_METRICS_EXPORTER=otlp",
		"NUMARY_OTEL_METRICS_EXPORTER_OTLP_ENDPOINT=" + configuration.OTLP.Endpoint,
		"NUMARY_OTEL_METRICS_EXPORTER_OTLP_INSECURE=true",
		"NUMARY_OTEL_RESOURCE_ATTRIBUTES=runID=" + runID,
	}
}

func v2EnvVars(runID string, configuration LedgerConfiguration) []string {
	return []string{
		"BIND=:3068",
		"STORAGE_DRIVER=postgres",
		"STORAGE_POSTGRES_CONN_STRING=" + configuration.PostgresDSN,
		"OTEL_TRACES=true",
		"OTEL_TRACES_EXPORTER=otlp",
		"OTEL_TRACES_EXPORTER_OTLP_ENDPOINT=" + configuration.OTLP.Endpoint,
		"OTEL_TRACES_EXPORTER_OTLP_INSECURE=true",
		"OTEL_METRICS=true",
		"OTEL_METRICS_EXPORTER=otlp",
		"OTEL_METRICS_EXPORTER_OTLP_ENDPOINT=" + configuration.OTLP.Endpoint,
		"OTEL_METRICS_EXPORTER_OTLP_INSECURE=true",
		"OTEL_METRICS_RUNTIME=true",
		"OTEL_RESOURCE_ATTRIBUTES=runID=" + runID,
		"DEBUG=true",
	}
}

func init() {
	ext := &Extension{
		logger: logrus.New(),
	}
	modules.Register("k6/x/formancehq/benchmarks", ext)
}
