package cmd

import (
	"context"
	"net/http"

	sharedapi "github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/auth"
	"github.com/formancehq/go-libs/licence"
	"github.com/formancehq/go-libs/otlp"
	"github.com/formancehq/go-libs/otlp/otlptraces"
	"github.com/formancehq/go-libs/service"
	wallet "github.com/formancehq/wallets/pkg"
	"github.com/formancehq/wallets/pkg/api"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	stackClientIDFlag     = "stack-client-id"
	stackClientSecretFlag = "stack-client-secret"
	stackURLFlag          = "stack-url"
	ledgerNameFlag        = "ledger"
	accountPrefixFlag     = "account-prefix"
	listenFlag            = "listen"
)

func newServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"server"},
		RunE: func(cmd *cobra.Command, args []string) error {
			stackClientID, _ := cmd.Flags().GetString(stackClientIDFlag)
			stackClientSecret, _ := cmd.Flags().GetString(stackClientSecretFlag)
			stackURL, _ := cmd.Flags().GetString(stackURLFlag)
			ledgerName, _ := cmd.Flags().GetString(ledgerNameFlag)
			accountPrefix, _ := cmd.Flags().GetString(accountPrefixFlag)
			listen, _ := cmd.Flags().GetString(listenFlag)

			options := []fx.Option{
				fx.Provide(func() (*http.Client, error) {
					return GetHTTPClient(
						cmd.Context(),
						stackClientID,
						stackClientSecret,
						stackURL,
						service.IsDebug(cmd),
					)
				}),
				wallet.Module(
					stackURL,
					ledgerName,
					accountPrefix,
				),
				api.Module(sharedapi.ServiceInfo{
					Version: Version,
					Debug:   service.IsDebug(cmd),
				}, listen),
				otlptraces.FXModuleFromFlags(cmd),
				auth.FXModuleFromFlags(cmd),
				licence.FXModuleFromFlags(cmd, ServiceName),
			}

			return service.New(cmd.OutOrStdout(), options...).Run(cmd)
		},
	}
	cmd.Flags().String(stackClientIDFlag, "", "Client ID")
	cmd.Flags().String(stackClientSecretFlag, "", "Client Secret")
	cmd.Flags().String(stackURLFlag, "", "Token URL")
	cmd.Flags().String(ledgerNameFlag, "wallets-002", "Target ledger")
	cmd.Flags().String(accountPrefixFlag, "", "Account prefix flag")
	cmd.Flags().String(listenFlag, ":8080", "Listen address")

	service.AddFlags(cmd.Flags())
	licence.AddFlags(cmd.Flags())
	auth.AddFlags(cmd.Flags())
	otlptraces.AddFlags(cmd.Flags())

	return cmd
}

func GetHTTPClient(ctx context.Context, clientID, clientSecret, stackURL string, debug bool) (*http.Client, error) {
	httpClient := &http.Client{
		Transport: otlp.NewRoundTripper(http.DefaultTransport, debug),
	}

	if clientID == "" {
		return httpClient, nil
	}

	clientCredentialsConfig := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     stackURL + "/api/auth/oauth/token",
		Scopes:       []string{"openid ledger:read ledger:write"},
	}

	return clientCredentialsConfig.Client(context.WithValue(ctx, oauth2.HTTPClient, httpClient)), nil
}
