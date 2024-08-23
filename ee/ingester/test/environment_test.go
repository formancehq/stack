//go:build it

package test_suite

import (
	"encoding/json"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/testing/docker"
	"github.com/formancehq/stack/libs/go-libs/testing/platform/clickhousetesting"
	. "github.com/formancehq/stack/libs/go-libs/testing/platform/natstesting"
	. "github.com/formancehq/stack/libs/go-libs/testing/platform/pgtesting"
	. "github.com/formancehq/stack/libs/go-libs/testing/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"slices"
)

var (
	dockerPool       = NewDeferred[*docker.Pool]()
	pgServer         = NewDeferred[*PostgresServer]()
	clickhouseServer = NewDeferred[*clickhousetesting.Server]()
	natsServer       = NewDeferred[*NatsServer]()
	debug            = os.Getenv("DEBUG") == "true"
	logger           = logging.NewDefaultLogger(GinkgoWriter, debug, false)
)

type ParallelExecutionContext struct {
	PostgresServer   *PostgresServer
	ClickhouseServer *clickhousetesting.Server
	NatsServer       *NatsServer
}

var _ = SynchronizedBeforeSuite(func() []byte {
	By("Initializing docker pool")
	dockerPool.SetValue(docker.NewPool(GinkgoT(), logger))

	pgServer.LoadAsync(func() *PostgresServer {
		By("Initializing postgres server")
		return CreatePostgresServer(GinkgoT(), dockerPool.GetValue())
	})

	natsServer.LoadAsync(func() *NatsServer {
		By("Initializing nats server")
		return CreateServer(GinkgoT(), true, logger)
	})

	if slices.Contains(EnabledConnectors(), "clickhouse") {
		clickhouseServer.LoadAsync(func() *clickhousetesting.Server {
			By("Initializing clickhouse server")
			return clickhousetesting.CreateServer(dockerPool.GetValue())
		})
	} else {
		clickhouseServer.SetValue(&clickhousetesting.Server{})
	}

	By("Waiting services alive")
	Wait(pgServer, natsServer, clickhouseServer)
	By("All services ready.")

	data, err := json.Marshal(ParallelExecutionContext{
		PostgresServer:   pgServer.GetValue(),
		ClickhouseServer: clickhouseServer.GetValue(),
		NatsServer:       natsServer.GetValue(),
	})
	Expect(err).To(BeNil())

	return data
}, func(data []byte) {
	select {
	case <-pgServer.Done():
		// Process #1, setup is terminated
		return
	default:
	}
	pec := ParallelExecutionContext{}
	err := json.Unmarshal(data, &pec)
	Expect(err).To(BeNil())

	pgServer.SetValue(pec.PostgresServer)
	clickhouseServer.SetValue(pec.ClickhouseServer)
	natsServer.SetValue(pec.NatsServer)
})
