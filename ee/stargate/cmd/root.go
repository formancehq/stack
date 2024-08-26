package cmd

import (
	"time"

	"github.com/formancehq/stack/libs/go-libs/licence"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
)

var (
	Version   = "develop"
	BuildDate = "-"
	Commit    = "-"
)

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:               "stargate",
		Short:             "stargate",
		DisableAutoGenTag: true,
		Version:           Version,
	}

	version := newVersion()
	root.AddCommand(version)

	client := newClient()
	root.AddCommand(client)

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

	service.AddFlags(client.PersistentFlags())
	licence.AddFlags(client.PersistentFlags())
	otlptraces.AddFlags(root.PersistentFlags())
	otlpmetrics.AddFlags(root.PersistentFlags())

	return root
}

func Execute() {
	service.Execute(NewRootCommand())
}
