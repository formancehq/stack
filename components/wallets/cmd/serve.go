package cmd

import (
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/service"
	wallet "github.com/formancehq/wallets/pkg"
	"github.com/formancehq/wallets/pkg/api"
	"github.com/formancehq/wallets/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
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
		Use: "serve",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			options := []fx.Option{
				wallet.Module(
					viper.GetString(stackURLFlag)+"/api/ledger",
					viper.GetString(ledgerNameFlag),
					viper.GetString(accountPrefixFlag),
				),
				api.Module(sharedapi.ServiceInfo{
					Version: Version,
				}, viper.GetString(listenFlag)),
				client.NewModule(viper.GetString(stackClientIDFlag), viper.GetString(stackClientSecretFlag),
					viper.GetString(stackURLFlag), viper.GetBool(service.DebugFlag)),
				otlptraces.CLITracesModule(viper.GetViper()),
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
