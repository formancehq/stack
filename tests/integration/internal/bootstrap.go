package internal

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/benthosdev/benthos/v4/public/components/pure"
	"github.com/docker/docker/client"
	authCmd "github.com/formancehq/auth/cmd"
	auth "github.com/formancehq/auth/pkg"
	paymentsCmd "github.com/formancehq/payments/cmd"
	searchCmd "github.com/formancehq/search/cmd"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/logging"
	walletsCmd "github.com/formancehq/wallets/cmd"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/egymgmbh/go-prefix-writer/prefixer"
	natsServer "github.com/nats-io/nats-server/v2/server"
	"github.com/numary/ledger/cmd"
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
	cancel       func()
	dockerPool   *dockertest.Pool
	dockerClient *client.Client
)

func TestContext() context.Context {
	return ctx
}

var _ = BeforeSuite(func() {
	// Some defaults
	SetDefaultEventuallyTimeout(5 * time.Second)

	var err error
	dockerPool, err = dockertest.NewPool("")
	Expect(err).To(BeNil())

	dockerClient, err = client.NewClientWithOpts(client.FromEnv)
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
	actualTestID = uuid.NewString()
	ctx, cancel = context.WithCancel(context.TODO())
	l := logrus.New()
	l.Out = io.Discard
	if testing.Verbose() {
		l.Level = logrus.DebugLevel
		l.Out = os.Stdout
	}
	ctx = logging.ContextWithLogger(ctx, logging.New(l))

	startBenthosServer()
	createDatabases() // TODO: drop databases

	startFakeGateway()

	startLedger()
	startSearch()
	startAuth()
	startWallets()
	startPayments()

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
		"server", "start",
		"--publisher-nats-enabled",
		"--publisher-nats-client-id=ledger",
		"--publisher-nats-url=" + natsAddress(),
		fmt.Sprintf("--publisher-topic-mapping=*:%s-ledger", actualTestID),
		"--storage.driver=postgres",
		"--storage.postgres.conn_string=" + dsn.String(),
		"--server.http.bind_address=0.0.0.0:0", // Random port
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
		"server",
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
		"server",
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

func runAndWaitPort(service string, cmd *cobra.Command) (int, context.CancelFunc, chan error) {

	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	if testing.Verbose() {
		cmd.SetOut(prefixer.New(os.Stdout, func() string {
			return service + " | "
		}))
		cmd.SetErr(prefixer.New(os.Stderr, func() string {
			return service + " | "
		}))
	}
	ctx := httpserver.ContextWithServerInfo(TestContext())
	ctx, cancel := context.WithCancel(ctx)
	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.ExecuteContext(ctx)
	}()
	select {
	case <-httpserver.Started(ctx):
	case err := <-errCh:
		Expect(err).To(BeNil())
	case <-time.After(5 * time.Second):
		Fail("timeout waiting for service to be properly started")
	}
	port := httpserver.Port(ctx)

	return port, cancel, errCh
}
