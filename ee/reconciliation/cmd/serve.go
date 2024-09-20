package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/formancehq/go-libs/aws/iam"

	"github.com/formancehq/go-libs/bun/bunconnect"
	"github.com/formancehq/go-libs/licence"

	sdk "github.com/formancehq/formance-sdk-go/v2"
	sharedapi "github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/auth"
	"github.com/formancehq/go-libs/otlp"
	"github.com/formancehq/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/go-libs/otlp/otlptraces"
	"github.com/formancehq/go-libs/service"
	"github.com/formancehq/reconciliation/internal/api"
	"github.com/formancehq/reconciliation/internal/storage"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func stackClientModule(cmd *cobra.Command) fx.Option {
	return fx.Options(
		fx.Provide(func() *sdk.Formance {
			stackClientID, _ := cmd.Flags().GetString(stackClientIDFlag)
			stackClientSecret, _ := cmd.Flags().GetString(stackClientSecretFlag)
			stackURL, _ := cmd.Flags().GetString(stackURLFlag)

			oauthConfig := clientcredentials.Config{
				ClientID:     stackClientID,
				ClientSecret: stackClientSecret,
				TokenURL:     fmt.Sprintf("%s/api/auth/oauth/token", stackURL),
				Scopes:       []string{"openid", "ledger:read", "ledger:write", "payments:read", "payments:write"},
			}
			underlyingHTTPClient := &http.Client{
				Transport: otlp.NewRoundTripper(http.DefaultTransport, service.IsDebug(cmd)),
			}
			return sdk.New(
				sdk.WithClient(
					oauthConfig.Client(context.WithValue(context.Background(),
						oauth2.HTTPClient, underlyingHTTPClient)),
				),
				sdk.WithServerURL(stackURL),
			)
		}),
	)
}

func newServeCommand(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "serve",
		RunE: runServer(version),
	}
	cmd.Flags().String(listenFlag, ":8080", "Listening address")
	cmd.Flags().String(stackURLFlag, "", "Stack url")
	cmd.Flags().String(stackClientIDFlag, "", "Stack client ID")
	cmd.Flags().String(stackClientSecretFlag, "", "Stack client secret")

	otlpmetrics.AddFlags(cmd.Flags())
	otlptraces.AddFlags(cmd.Flags())
	auth.AddFlags(cmd.Flags())
	bunconnect.AddFlags(cmd.Flags())
	iam.AddFlags(cmd.Flags())
	service.AddFlags(cmd.Flags())
	licence.AddFlags(cmd.Flags())

	return cmd
}

func runServer(version string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		databaseOptions, err := prepareDatabaseOptions(cmd)
		if err != nil {
			return err
		}

		options := make([]fx.Option, 0)
		options = append(options, databaseOptions)

		options = append(options,
			otlptraces.FXModuleFromFlags(cmd),
			otlpmetrics.FXModuleFromFlags(cmd),
			auth.FXModuleFromFlags(cmd),
		)

		listen, _ := cmd.Flags().GetString(listenFlag)
		options = append(options,
			stackClientModule(cmd),
			api.HTTPModule(sharedapi.ServiceInfo{
				Version: version,
				Debug:   service.IsDebug(cmd),
			}, listen),
			licence.FXModuleFromFlags(cmd, ServiceName),
		)

		return service.New(cmd.OutOrStdout(), options...).Run(cmd)
	}
}

func prepareDatabaseOptions(cmd *cobra.Command) (fx.Option, error) {
	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd)
	if err != nil {
		return nil, err
	}

	return storage.Module(*connectionOptions, service.IsDebug(cmd)), nil
}
