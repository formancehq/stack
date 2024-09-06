//go:build it

package test_suite

import (
	"fmt"
	. "github.com/formancehq/stack/ee/ingester/pkg/testserver"
	"github.com/formancehq/stack/libs/go-libs/testing/platform/clickhousetesting"
	. "github.com/formancehq/stack/libs/go-libs/testing/utils"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	"os"
	"strings"
	"time"
)

var (
	defaultEnabledConnectors = []string{"http", "clickhouse"}
)

func EnabledConnectors() []string {
	fromEnv := os.Getenv("CONNECTORS")
	if fromEnv == "" {
		return defaultEnabledConnectors
	}
	return strings.Split(fromEnv, ",")
}

// connectorsSetup allow to define a ginkgo node factory function for each connector
// This allows to configure the environment for the connector
var connectorsSetup = map[string]func(func(d *Deferred[Connector])){
	"http": func(fn func(d *Deferred[Connector])) {
		WithConnector("http", func() Connector {
			return NewHTTPConnector(GinkgoT(), NewCollector())
		}, fn)
	},
	"clickhouse": func(fn func(d *Deferred[Connector])) {
		clickhousetesting.WithNewDatabase(clickhouseServer, func(db *Deferred[*clickhousetesting.Database]) {
			WithConnector("clickhouse", func() Connector {
				return NewClickhouseConnector(logger, db.GetValue().ConnString())
			}, fn)
		})
	},
}

// todo: test opentelemetry integration
var _ = Context("Data ingestion integration tests", func() {
	stackName := NewDeferred[string]()
	BeforeEach(func() {
		stackName.Reset()
		stackName.SetValue(uuid.NewString()[:8])

		natsServer.GetValue().CreateStream(GinkgoT(), stackName.GetValue())
		natsServer.GetValue().CreateConsumer(GinkgoT(), stackName.GetValue(), &nats.ConsumerConfig{
			Durable: stackName.GetValue() + "-ingester",
			Name:    stackName.GetValue() + "-ingester",
			FilterSubjects: []string{
				fmt.Sprintf("%s.module1", stackName.GetValue()),
				fmt.Sprintf("%s.module2", stackName.GetValue()),
			},
			DeliverSubject: stackName.GetValue() + "-delivery",
			DeliverGroup:   stackName.GetValue() + "-ingester",
			AckPolicy:      nats.AckExplicitPolicy,
			MemoryStorage:  true,
			AckWait:        time.Second,
		})
	})
	for _, connectorName := range EnabledConnectors() {
		setup, ok := connectorsSetup[connectorName]
		if !ok {
			Fail(fmt.Sprintf("Driver '%s' not exists", connectorName))
		}
		setup(func(connector *Deferred[Connector]) {
			runTest(stackName, natsServer, connectorName, connector)
		})
	}
})
