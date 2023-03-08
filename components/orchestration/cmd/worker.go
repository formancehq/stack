package cmd

import (
	"context"
	"fmt"
	"net/http"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/orchestration/internal/temporal"
	"github.com/formancehq/stack/libs/go-libs/otlp"
	"github.com/formancehq/stack/libs/go-libs/service"
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
			configuration.Debug = viper.GetBool(service.DebugFlag)
			oauthConfig := clientcredentials.Config{
				ClientID:     viper.GetString(stackClientIDFlag),
				ClientSecret: viper.GetString(stackClientSecretFlag),
				TokenURL:     fmt.Sprintf("%s/api/auth/oauth/token", viper.GetString(stackURLFlag)),
			}
			underlyingHTTPClient := &http.Client{
				Transport: otlp.NewRoundTripper(viper.GetBool(service.DebugFlag)),
			}
			configuration.HTTPClient = oauthConfig.Client(context.WithValue(context.Background(),
				oauth2.HTTPClient, underlyingHTTPClient))
			return sdk.NewAPIClient(configuration)
		}),
	)
}

func workerOptions() fx.Option {
	return fx.Options(
		stackClientModule(),
		temporal.NewWorkerModule(viper.GetString(temporalTaskQueueFlag)),
	)
}

func newWorkerCommand() *cobra.Command {
	return &cobra.Command{
		Use: "worker",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			options := []fx.Option{
				commonOptions(),
				workerOptions(),
			}

			return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
		},
	}
}
