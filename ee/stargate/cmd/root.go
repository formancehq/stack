package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/formancehq/stack/libs/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version   = "develop"
	BuildDate = "-"
	Commit    = "-"
)

func NewRootCommand() *cobra.Command {
	viper.SetDefault("version", Version)

	root := &cobra.Command{
		Use:               "stargate",
		Short:             "stargate",
		DisableAutoGenTag: true,
	}

	version := newVersion()
	root.AddCommand(version)

	server := newServer()
	root.AddCommand(server)

	client := newClient()
	root.AddCommand(client)

	publish.InitCLIFlags(server, func(cd *publish.ConfigDefault) {
		// We want to override the default values of flags here in order for the
		// Max reconnect flag to be set to -1, which means infinite reconnects.
		// NOTE(polo): the value provided in the env vars will still override
		// this one, so make sure to check if the value is set or not in the env
		// vars
		cd.PublisherNatsMaxReconnect = -1
	})
	server.Flags().String(serviceHttpAddrFlag, "localhost:8080", "Listen address for http API")
	server.Flags().String(serviceGrpcAddrFlag, "localhost:3068", "Listen address for grpc API")
	server.Flags().String(authIssuerURLFlag, "", "JWKS URL")
	server.Flags().Int(maxRetriesJWKSFetchingFlag, 3, "Max retries for fetching JWKS")
	server.Flags().Duration(natsRequestTimeout, 10*time.Second, "NATS request timeout (in seconds)")
	if err := viper.BindPFlags(server.Flags()); err != nil {
		panic(err)
	}
	if err := viper.BindPFlags(server.PersistentFlags()); err != nil {
		panic(err)
	}

	client.Flags().String(organizationIDFlag, "", "Organization ID")
	client.Flags().String(stackIDFlag, "", "Stack ID")
	client.Flags().String(bindFlag, "0.0.0.0:8080", "Listen address for http API")
	client.Flags().String(stargateServerURLFlag, "", "Stargate server URL")
	client.Flags().String(gatewayURLFlag, "", "Gateway URL")
	client.Flags().Int(workerPoolMaxWorkersFlag, 100, "Max worker pool size")
	client.Flags().Int(workerPoolMaxTasksFlag, 10000, "Max worker pool tasks")
	client.Flags().Int(ClientChanSizeFlag, 1024, "Client chan size")
	client.Flags().Duration(HTTPClientTimeoutFlag, 10*time.Second, "HTTP client timeout")
	client.Flags().Int(HTTPClientMaxIdleConnsFlag, 100, "HTTP client max idle conns")
	client.Flags().Int(HTTPClientMaxIdleConnsPerHostFlag, 2, "HTTP client max idle conns per host")
	client.Flags().Duration(AuthRefreshTokenDurationBeforeExpireTimeFlag, 30*time.Second, "Auth refresh token duration")
	client.Flags().String(StargateAuthClientIDFlag, "", "Stargate auth client ID")
	client.Flags().String(StargateAuthClientSecretFlag, "", "Stargate auth client secret")
	client.Flags().String(StargateAuthIssuerURLFlag, "", "Stargate auth issuer")
	client.Flags().Bool(TlsEnabledFlag, true, "TLS enabled")
	client.Flags().String(TlsCACertificateFlag, "", "TLS cert file")
	client.Flags().Bool(TlsInsecureSkipVerifyFlag, false, "TLS insecure skip verify")
	if err := viper.BindPFlags(client.Flags()); err != nil {
		panic(err)
	}
	if err := viper.BindPFlags(client.PersistentFlags()); err != nil {
		panic(err)
	}

	root.PersistentFlags().Bool(service.DebugFlag, false, "Debug mode")
	otlptraces.InitOTLPTracesFlags(root.PersistentFlags())
	otlpmetrics.InitOTLPMetricsFlags(root.PersistentFlags())

	if err := viper.BindPFlags(root.PersistentFlags()); err != nil {
		panic(err)
	}
	if err := viper.BindPFlags(root.Flags()); err != nil {
		panic(err)
	}

	BindEnv(viper.GetViper())

	return root
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		if _, err = fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}

		os.Exit(1)
	}
}

var EnvVarReplacer = strings.NewReplacer(".", "_", "-", "_")

func BindEnv(v *viper.Viper) {
	v.SetEnvKeyReplacer(EnvVarReplacer)
	v.AutomaticEnv()
}
