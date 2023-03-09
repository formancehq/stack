// TODO(gfyrag): The code of this file is pretty dirty. I have other priority when writing this.
// I will clean this later.
// It works as expected.
package internal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/docker/client"
	authCmd "github.com/formancehq/auth/cmd"
	auth "github.com/formancehq/auth/pkg"
	orchestrationCmd "github.com/formancehq/orchestration/cmd"
	paymentsCmd "github.com/formancehq/payments/cmd"
	searchCmd "github.com/formancehq/search/cmd"
	"github.com/formancehq/search/pkg/searchengine"
	"github.com/formancehq/stack/libs/go-libs/httpclient"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/logging"
	walletsCmd "github.com/formancehq/wallets/cmd"
	webhooksCmd "github.com/formancehq/webhooks/cmd"
	"github.com/google/uuid"
	"github.com/opensearch-project/opensearch-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/egymgmbh/go-prefix-writer/prefixer"
	"github.com/formancehq/ledger/cmd"
	natsServer "github.com/nats-io/nats-server/v2/server"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/spf13/cobra"
)

var (
	server       *natsServer.Server
	actualTestID string
	ctx          context.Context
	dockerPool   *dockertest.Pool
	dockerClient *client.Client
)

func TestContext() context.Context {
	return ctx
}

var _ = BeforeSuite(func() {
	// Some defaults
	SetDefaultEventuallyTimeout(10 * time.Second)
	SetDefaultEventuallyPollingInterval(500 * time.Millisecond)

	var err error
	dockerPool, err = dockertest.NewPool("")
	Expect(err).To(BeNil())

	dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	Expect(err).To(BeNil())

	// uses pool to try to connect to Docker
	err = dockerPool.Client.Ping()
	Expect(err).To(BeNil())
})

func runDockerResource(options *dockertest.RunOptions) *dockertest.Resource {
	resource, err := dockerPool.RunWithOptions(options, func(config *docker.HostConfig) {
		config.AutoRemove = true
	})
	Expect(err).To(BeNil())

	return resource
}

var _ = BeforeEach(func() {
	ctx = context.TODO()
	actualTestID = uuid.NewString()
	l := logrus.New()
	l.Out = GinkgoWriter
	l.Level = logrus.DebugLevel
	ctx = logging.ContextWithLogger(ctx, logging.NewLogrus(l))
	openSearchClient, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{"http://" + getOpenSearchUrl()},
		Transport: httpclient.NewDebugHTTPTransport(http.DefaultTransport),
	})
	Eventually(func() error {
		return searchengine.CreateIndex(ctx, openSearchClient, actualTestID)
	}).WithTimeout(10 * time.Second).Should(BeNil())

	createDatabases() // TODO: drop databases

	startFakeGateway()

	startSearch()
	startLedger()
	startAuth()
	startWallets()
	startPayments()
	startWebhooks()
	startOrchestration()

	// TODO: Wait search has properly configured mapping before trying to ingest any data
	startBenthosServer()

	// Start the gateway
	ledgerUrl, err := url.Parse(fmt.Sprintf("http://localhost:%d", ledgerPort))
	if err != nil {
		panic(err)
	}
	registerService("ledger", ledgerUrl)

	searchUrl, err := url.Parse(fmt.Sprintf("http://localhost:%d", searchPort))
	if err != nil {
		panic(err)
	}
	registerService("search", searchUrl)

	authUrl, err := url.Parse(fmt.Sprintf("http://localhost:%d", authPort))
	if err != nil {
		panic(err)
	}
	registerService("auth", authUrl)

	walletsUrl, err := url.Parse(fmt.Sprintf("http://localhost:%d", walletsPort))
	if err != nil {
		panic(err)
	}
	registerService("wallets", walletsUrl)

	paymentsUrl, err := url.Parse(fmt.Sprintf("http://localhost:%d", paymentsPort))
	if err != nil {
		panic(err)
	}
	registerService("payments", paymentsUrl)

	webhooksUrl, err := url.Parse(fmt.Sprintf("http://localhost:%d", webhooksPort))
	if err != nil {
		panic(err)
	}
	registerService("webhooks", webhooksUrl)

	orchestrationUrl, err := url.Parse(fmt.Sprintf("http://localhost:%d", orchestrationPort))
	if err != nil {
		panic(err)
	}
	registerService("orchestration", orchestrationUrl)

	// Start services
	// Configure the sdk with a preconfigured auth client
	configureSDK()
})

var _ = AfterEach(func() {
	stopBenthosServer()
	stopLedger()
	stopSearch()
	stopAuth()
	stopWallets()
	stopPayments()
	stopWebhooks()
	stopOrchestration()
	stopFakeGateway() // TODO: Wait for gateway to be shutdown
})

var (
	ledgerPort   int
	ledgerErrCh  chan error
	ledgerCancel func()
)

func startLedger() {
	dsn, err := getPostgresDSN()
	Expect(err).To(BeNil())
	dsn.Path = fmt.Sprintf("%s-ledger", actualTestID)

	ledgerCmd := cmd.NewRootCommand()
	args := []string{
		"serve",
		"--publisher-nats-enabled",
		"--publisher-nats-client-id=ledger",
		"--publisher-nats-url=" + natsAddress(),
		fmt.Sprintf("--publisher-topic-mapping=*:%s-ledger", actualTestID),
		"--storage.postgres.conn_string=" + dsn.String(),
		"--bind=0.0.0.0:0", // Random port
	}
	if testing.Verbose() {
		args = append(args, "--debug")
	}
	ledgerCmd.SetArgs(args)
	ledgerPort, ledgerCancel, ledgerErrCh = runAndWaitPort("ledger", ledgerCmd)
}

func stopLedger() {
	ledgerCancel()
	select {
	case <-ledgerErrCh:
	case <-time.After(5 * time.Second):
		Fail("timeout waiting for ledger stopped")
	}
}

var (
	searchPort   int
	searchErrCh  chan error
	searchCancel func()
)

func startSearch() {
	searchCmd := searchCmd.NewRootCommand()
	args := make([]string, 0)
	args = append(args,
		"serve",
		"--open-search-service="+getOpenSearchUrl(),
		"--open-search-scheme=http",
		"--open-search-username=admin",
		"--open-search-password=admin",
		"--bind=0.0.0.0:0",
		fmt.Sprintf("--es-indices=%s", actualTestID),
		"--mapping-init-disabled",
	)
	searchCmd.SetArgs(args)
	searchPort, searchCancel, searchErrCh = runAndWaitPort("search", searchCmd)
}

func stopSearch() {
	searchCancel()
	select {
	case <-searchErrCh:
	case <-time.After(5 * time.Second):
		Fail("timeout waiting for search stopped")
	}
}

var (
	paymentsPort   int
	paymentsErrCh  chan error
	paymentsCancel func()
)

func startPayments() {
	dsn, err := getPostgresDSN()
	Expect(err).To(BeNil())
	dsn.Path = fmt.Sprintf("%s-payments", actualTestID)

	paymentsCmd := paymentsCmd.NewRootCommand()
	if testing.Verbose() {
		paymentsCmd.SetOut(os.Stdout)
		paymentsCmd.SetErr(os.Stderr)
	}

	args := make([]string, 0)
	args = append(args,
		"serve",
		"--postgres-uri="+dsn.String(),
		"--config-encryption-key=encryption-key",
		"--publisher-nats-enabled",
		"--publisher-nats-client-id=payments",
		"--publisher-nats-url="+natsAddress(),
		fmt.Sprintf("--publisher-topic-mapping=*:%s-payments", actualTestID),
		"--listen=0.0.0.0:0",
		"--auto-migrate",
	)
	paymentsCmd.SetArgs(args)
	paymentsPort, paymentsCancel, paymentsErrCh = runAndWaitPort("payments", paymentsCmd)
}

func stopPayments() {
	paymentsCancel()
	select {
	case <-paymentsErrCh:
	case <-time.After(5 * time.Second):
		Fail("timeout waiting for payments stopped")
	}
}

var (
	orchestrationPort   int
	orchestrationErrCh  chan error
	orchestrationCancel func()
)

func startOrchestration() {
	dsn, err := getPostgresDSN()
	Expect(err).To(BeNil())
	dsn.Path = fmt.Sprintf("%s-orchestration", actualTestID)

	orchestrationCmd := orchestrationCmd.NewRootCommand()
	if testing.Verbose() {
		orchestrationCmd.SetOut(os.Stdout)
		orchestrationCmd.SetErr(os.Stderr)
	}

	args := make([]string, 0)
	args = append(args,
		"serve",
		"--listen=0.0.0.0:0",
		"--postgres-dsn="+dsn.String(),
		"--stack-client-id=global",
		"--stack-client-secret=global",
		"--stack-url="+gatewayServer.URL,
		"--temporal-address="+getTemporalAddress(),
		"--temporal-task-queue="+actualTestID,
		"--worker",
	)
	orchestrationCmd.SetArgs(args)
	orchestrationPort, orchestrationCancel, orchestrationErrCh = runAndWaitPort("orchestration", orchestrationCmd)
}

func stopOrchestration() {
	orchestrationCancel()
	select {
	case <-orchestrationErrCh:
	case <-time.After(5 * time.Second):
		Fail("timeout waiting for orchestration stopped")
	}
}

var (
	authPort   int
	authErrCh  chan error
	authCancel func()
)

func startAuth() {
	dsn, err := getPostgresDSN()
	Expect(err).To(BeNil())
	dsn.Path = fmt.Sprintf("%s-auth", actualTestID)

	authCmd := authCmd.NewRootCommand()
	if testing.Verbose() {
		authCmd.SetOut(os.Stdout)
		authCmd.SetErr(os.Stderr)
	}

	authDir := filepath.Join(os.TempDir(), uuid.NewString())
	Expect(os.MkdirAll(authDir, 0777)).To(BeNil())
	type configuration struct {
		Clients []auth.StaticClient `yaml:"clients"`
	}
	cfg := &configuration{
		Clients: []auth.StaticClient{{
			ClientOptions: auth.ClientOptions{
				Name:    "global",
				Id:      "global",
				Trusted: true,
			},
			Secrets: []string{"global"},
		}},
	}
	configFile := filepath.Join(authDir, "config.yaml")
	f, err := os.Create(configFile)
	Expect(err).To(BeNil())
	Expect(yaml.NewEncoder(f).Encode(cfg)).To(BeNil())

	os.Setenv("CAOS_OIDC_DEV", "1")
	args := make([]string, 0)
	args = append(args,
		"serve",
		"--config="+configFile,
		"--postgres-uri="+dsn.String(),
		"--delegated-client-id=noop",
		"--delegated-client-secret=noop",
		"--delegated-issuer=https://accounts.google.com",
		"--base-url=http://localhost/api/auth",
		"--listen=0.0.0.0:0",
	)
	if testing.Verbose() {
		args = append(args, "--debug")
	}
	authCmd.SetArgs(args)
	authPort, authCancel, authErrCh = runAndWaitPort("auth", authCmd)
}

func stopAuth() {
	authCancel()
	select {
	case <-authErrCh:
	case <-time.After(5 * time.Second):
		Fail("timeout waiting for auth stopped")
	}
}

var (
	walletsPort   int
	walletsErrCh  chan error
	walletsCancel func()
)

func startWallets() {
	walletCmd := walletsCmd.NewRootCommand()
	if testing.Verbose() {
		walletCmd.SetOut(os.Stdout)
		walletCmd.SetErr(os.Stderr)
	}

	args := make([]string, 0)
	args = append(args,
		"serve",
		"--stack-client-id=global",
		"--stack-client-secret=global",
		"--stack-url="+gatewayServer.URL,
		"--listen=0.0.0.0:0",
	)
	if testing.Verbose() {
		args = append(args, "--debug")
	}
	walletCmd.SetArgs(args)
	walletsPort, walletsCancel, walletsErrCh = runAndWaitPort("wallets", walletCmd)
}

func stopWallets() {
	walletsCancel()
	select {
	case <-walletsErrCh:
	case <-time.After(5 * time.Second):
		Fail("timeout waiting for wallet stopped")
	}
}

var (
	webhooksPort   int
	webhooksErrCh  chan error
	webhooksCancel func()
)

func startWebhooks() {

	dsn, err := getPostgresDSN()
	Expect(err).To(BeNil())
	dsn.Path = fmt.Sprintf("%s-webhooks", actualTestID)

	webhooksCmd := webhooksCmd.NewRootCommand()
	if testing.Verbose() {
		webhooksCmd.SetOut(os.Stdout)
		webhooksCmd.SetErr(os.Stderr)
	}

	args := make([]string, 0)
	args = append(args,
		"serve",
		"--storage-postgres-conn-string="+dsn.String(),
		"--listen=0.0.0.0:0",
		"--worker",
		"--publisher-nats-enabled",
		"--publisher-nats-client-id=webhooks",
		"--publisher-nats-url="+natsAddress(),
		fmt.Sprintf("--kafka-topics=%s-ledger", actualTestID),
	)
	if testing.Verbose() {
		args = append(args, "--debug")
	}
	webhooksCmd.SetArgs(args)
	webhooksPort, webhooksCancel, webhooksErrCh = runAndWaitPort("webhooks", webhooksCmd)
}

func stopWebhooks() {
	webhooksCancel()
	select {
	case <-webhooksErrCh:
	case <-time.After(5 * time.Second):
		Fail("timeout waiting for webhooks stopped")
	}
}

func runAndWaitPort(service string, cmd *cobra.Command) (int, context.CancelFunc, chan error) {

	writer := prefixer.New(GinkgoWriter, func() string {
		return service + " | "
	})
	cmd.SetOut(writer)
	cmd.SetErr(writer)

	ctx := httpserver.ContextWithServerInfo(TestContext())
	ctx, cancel := context.WithCancel(ctx)
	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.ExecuteContext(ctx)
	}()
	select {
	case <-httpserver.Started(ctx):
	case err := <-errCh:
		By("starting service " + service)
		Expect(err).To(BeNil())
	case <-time.After(5 * time.Second):
		Fail("timeout waiting for service to be properly started")
	}
	port := httpserver.Port(ctx)

	return port, cancel, errCh
}
