package cmd

import (
	"context"
	"fmt"
	"net/http"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/go-libs/otlp"
	"github.com/formancehq/go-libs/otlp/otlptraces"
	"github.com/formancehq/orchestration/internal/storage"
	"github.com/formancehq/orchestration/internal/temporal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func stackClientModule() fx.Option {
	return fx.Options(
		fx.Provide(func() *sdk.APIClient {
			configuration := sdk.NewConfiguration()
			configuration.Servers = []sdk.ServerConfiguration{{
				URL: viper.GetString(stackURLFlag),
			}}
			configuration.Debug = viper.GetBool(debugFlag)
			oauthConfig := clientcredentials.Config{
				ClientID:     viper.GetString(stackClientIDFlag),
				ClientSecret: viper.GetString(stackClientSecretFlag),
				TokenURL:     fmt.Sprintf("%s/api/auth/oauth/token", viper.GetString(stackURLFlag)),
			}
			underlyingHTTPClient := &http.Client{
				Transport: otlp.NewRoundTripper(viper.GetBool(debugFlag)),
			}
			configuration.HTTPClient = oauthConfig.Client(context.WithValue(context.Background(),
				oauth2.HTTPClient, underlyingHTTPClient))
			return sdk.NewAPIClient(configuration)
		}),
	)
}

var workerCommand = &cobra.Command{
	Use: "worker",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return bindFlagsToViper(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		options := []fx.Option{
			fx.NopLogger,
			otlptraces.CLITracesModule(viper.GetViper()),
			storage.NewModule(viper.GetString(postgresDSNFlag), viper.GetBool(debugFlag)),
			stackClientModule(),
			temporal.NewClientModule(
				viper.GetString(temporalAddressFlag),
				viper.GetString(temporalNamespaceFlag),
				viper.GetString(temporalSSLClientCertFlag),
				viper.GetString(temporalSSLClientKeyFlag),
			),
			temporal.NewWorkerModule(viper.GetString(temporalTaskQueueFlag)),
		}

		app := fx.New(options...)
		err := app.Start(cmd.Context())
		if err != nil {
			return err
		}
		<-app.Done()
		return app.Err()
	},
}

func init() {
	rootCmd.AddCommand(workerCommand)
}
