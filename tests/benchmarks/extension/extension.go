package extension

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
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
	CPUCount    int64  `json:"cpuCount"`
	TestID      string `json:"testID"`
}

type Ledger struct {
	URL string `json:"url"`
}

func (c *LedgerConfiguration) resolve() error {
	if c.TestID == "" {
		c.TestID = os.Getenv("TEST_ID")
		if c.TestID == "" {
			return errors.New("missing test id")
		}
	}
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
			if c.pool != nil {
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

	var envVars []string
	if semver.IsValid(configuration.Version) && semver.Compare(configuration.Version, "v2.0.0") < 0 {
		envVars = v1EnvVars(c.testID, configuration)
	} else {
		envVars = v2EnvVars(c.testID, configuration)
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
	if configuration.CPUCount != 0 {
		options = append(options, func(config *docker.HostConfig) {
			config.CPUCount = configuration.CPUCount
		})
	}

	logger.Infof("Starting ledger container...")
	c.resource, err = c.pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "ghcr.io/formancehq/ledger",
		Tag:          configuration.Version,
		Env:          envVars,
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

func (c *Extension) ExportResults() {
	client, err := api.NewClient(api.Config{
		Address: "http://localhost:9090", // TODO: Use env var
	})
	if err != nil {
		panic(err)
	}

	<-time.After(5 * time.Second) // TODO: Check if the delay is enough

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	labels, warnings, err := v1api.Series(ctx, []string{fmt.Sprintf(`{testID="%s"}`, c.testID)}, time.Now().Add(-3*time.Hour), time.Now())
	if err != nil {
		panic(err)
	}
	if len(warnings) > 0 {
		c.logger.Warnf("Warnings: %v", warnings)
	}

	visitedLabel := map[model.LabelValue]struct{}{}

	if err := os.MkdirAll(filepath.Join(".", "results"), 0755); err != nil {
		panic(err)
	}

	f, err := os.Create(filepath.Join("results", fmt.Sprintf("%s.txt", c.testID)))
	if err != nil {
		panic(err)
	}

	enc := json.NewEncoder(f)

	for _, labelSet := range labels {
		name := labelSet["__name__"]
		_, alreadyVisited := visitedLabel[name]
		if alreadyVisited {
			continue
		}
		visitedLabel[name] = struct{}{}

		timeSeries, warnings, err := v1api.Query(ctx, fmt.Sprintf(`%s{testID="%s"}[1h]`, name, c.testID), time.Time{}, v1.WithTimeout(5*time.Second))
		if err != nil {
			panic(err)
		}
		if len(warnings) > 0 {
			c.logger.Warnf("Warnings: %v", warnings)
		}

		if err := enc.Encode(timeSeries); err != nil {
			panic(err)
		}
	}
}

func v1EnvVars(testID string, configuration LedgerConfiguration) []string {
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
		"NUMARY_OTEL_RESOURCE_ATTRIBUTES=testid=" + testID,
	}
}

func v2EnvVars(testID string, configuration LedgerConfiguration) []string {
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
		"OTEL_RESOURCE_ATTRIBUTES=testid=" + testID,
		"DEBUG=true",
	}
}

func init() {
	ext := &Extension{
		logger: logrus.New(),
	}
	modules.Register("k6/x/formancehq/benchmarks", ext)
}
