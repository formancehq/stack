package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/formancehq/orchestration/internal/storage"
	"github.com/formancehq/orchestration/internal/temporalclient"
	"github.com/formancehq/orchestration/internal/triggers"
	"github.com/formancehq/orchestration/internal/workflow"
	_ "github.com/formancehq/orchestration/internal/workflow/stages/all"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/bun/bunmigrate"
	"github.com/formancehq/stack/libs/go-libs/licence"
	"github.com/formancehq/stack/libs/go-libs/otlp"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/formancehq/stack/libs/go-libs/temporal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	ServiceName = "orchestration"
	Version     = "develop"
	BuildDate   = "-"
	Commit      = "-"
)

const (
	stackFlag             = "stack"
	stackURLFlag          = "stack-url"
	stackClientIDFlag     = "stack-client-id"
	stackClientSecretFlag = "stack-client-secret"
	topicsFlag            = "topics"
	listenFlag            = "listen"
	workerFlag            = "worker"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
	}

	cobra.EnableTraverseRunHooks = true

	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmd.PersistentFlags().String(stackURLFlag, "", "Stack url")
	cmd.PersistentFlags().String(stackClientIDFlag, "", "Stack client ID")
	cmd.PersistentFlags().String(stackClientSecretFlag, "", "Stack client secret")
	cmd.PersistentFlags().StringSlice(topicsFlag, []string{}, "Topics to listen")
	cmd.PersistentFlags().String(stackFlag, "", "Stack")
	cmd.AddCommand(
		newServeCommand(),
		newVersionCommand(),
		newWorkerCommand(),
		bunmigrate.NewDefaultCommand(func(cmd *cobra.Command, args []string, db *bun.DB) error {
			return storage.Migrate(cmd.Context(), db)
		}),
	)

	publish.InitCLIFlags(cmd)
	auth.InitAuthFlags(cmd.PersistentFlags())
	bunconnect.InitFlags(cmd.PersistentFlags())
	iam.InitFlags(cmd.PersistentFlags())
	service.BindFlags(cmd)
	licence.InitCLIFlags(cmd)
	temporal.InitCLIFlags(cmd)

	return cmd
}

func exitWithCode(code int, v ...any) {
	fmt.Fprintln(os.Stdout, v...)
	os.Exit(code)
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		exitWithCode(1, err)
	}
}

func commonOptions(cmd *cobra.Command) (fx.Option, error) {
	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd.Context())
	if err != nil {
		return nil, err
	}
	return fx.Options(
		otlptraces.CLITracesModule(),
		temporal.NewModule(
			workflow.Tracer,
			temporal.SearchAttributes{
				SearchAttributes: temporalclient.SearchAttributes,
			},
		),
		bunconnect.Module(*connectionOptions),
		publish.CLIPublisherModule("orchestration"),
		auth.CLIAuthModule(),
		licence.CLIModule(ServiceName),
		workflow.NewModule(viper.GetString(temporal.TemporalTaskQueueFlag)),
		triggers.NewModule(viper.GetString(temporal.TemporalTaskQueueFlag)),
		fx.Provide(func() *bunconnect.ConnectionOptions {
			return connectionOptions
		}),
		fx.Provide(func() *http.Client {
			httpClient := &http.Client{
				Transport: otlp.NewRoundTripper(http.DefaultTransport, viper.GetBool(service.DebugFlag)),
			}

			if viper.GetString(stackClientIDFlag) == "" {
				return httpClient
			}
			oauthConfig := clientcredentials.Config{
				ClientID:     viper.GetString(stackClientIDFlag),
				ClientSecret: viper.GetString(stackClientSecretFlag),
				TokenURL:     fmt.Sprintf("%s/api/auth/oauth/token", viper.GetString(stackURLFlag)),
				Scopes:       []string{"openid", "ledger:read", "ledger:write", "wallets:read", "wallets:write", "payments:read", "payments:write"},
			}
			return oauthConfig.Client(context.WithValue(context.Background(),
				oauth2.HTTPClient, httpClient))
		}),
	), nil
}
