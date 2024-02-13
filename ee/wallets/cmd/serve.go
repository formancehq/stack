package cmd

import (
	"context"
	"net/http"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/otlp"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/service"
	wallet "github.com/formancehq/wallets/pkg"
	"github.com/formancehq/wallets/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			options := []fx.Option{
				fx.Provide(func() (*http.Client, error) {
					return GetHTTPClient(cmd.Context(), viper.GetString(stackClientIDFlag), viper.GetString(stackClientSecretFlag),
						viper.GetString(stackURLFlag), viper.GetBool(service.DebugFlag))
				}),
				wallet.Module(
					viper.GetString(stackURLFlag),
					viper.GetString(ledgerNameFlag),
					viper.GetString(accountPrefixFlag),
				),
				api.Module(sharedapi.ServiceInfo{
					Version: Version,
				}, viper.GetString(listenFlag)),
				otlptraces.CLITracesModule(),
				auth.CLIAuthModule(),
			}

			return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
		},
	}
	cmd.Flags().String(stackClientIDFlag, "", "Client ID")
	cmd.Flags().String(stackClientSecretFlag, "", "Client Secret")
	cmd.Flags().String(stackURLFlag, "", "Token URL")
	cmd.Flags().String(ledgerNameFlag, "wallets-002", "Target ledger")
	cmd.Flags().String(accountPrefixFlag, "", "Account prefix flag")
	cmd.Flags().String(listenFlag, ":8080", "Listen address")
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
