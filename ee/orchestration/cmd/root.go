package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/formancehq/orchestration/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunmigrate"
	"github.com/formancehq/stack/libs/go-libs/licence"
	"github.com/uptrace/bun"

	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/otlp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/formancehq/orchestration/internal/triggers"
	"github.com/formancehq/orchestration/internal/workflow"

	"github.com/formancehq/orchestration/internal/temporalclient"
	"github.com/formancehq/stack/libs/go-libs/publish"

	_ "github.com/formancehq/orchestration/internal/workflow/stages/all"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var (
	ServiceName = "orchestration"
	Version     = "develop"
	BuildDate   = "-"
	Commit      = "-"
)

const (
	stackFlag                    = "stack"
	stackURLFlag                 = "stack-url"
	stackClientIDFlag            = "stack-client-id"
	stackClientSecretFlag        = "stack-client-secret"
	temporalAddressFlag          = "temporal-address"
	temporalNamespaceFlag        = "temporal-namespace"
	temporalSSLClientKeyFlag     = "temporal-ssl-client-key"
	temporalSSLClientCertFlag    = "temporal-ssl-client-cert"
	temporalTaskQueueFlag        = "temporal-task-queue"
	temporalInitSearchAttributes = "temporal-init-search-attributes"
	topicsFlag                   = "topics"
	listenFlag                   = "listen"
	workerFlag                   = "worker"
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
	cmd.PersistentFlags().String(temporalAddressFlag, "", "Temporal server address")
	cmd.PersistentFlags().String(temporalNamespaceFlag, "default", "Temporal namespace")
	cmd.PersistentFlags().String(temporalSSLClientKeyFlag, "", "Temporal client key")
	cmd.PersistentFlags().String(temporalSSLClientCertFlag, "", "Temporal client cert")
	cmd.PersistentFlags().String(temporalTaskQueueFlag, "default", "Temporal task queue name")
	cmd.PersistentFlags().Bool(temporalInitSearchAttributes, false, "Init temporal search attributes")
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
		temporalclient.NewModule(
			viper.GetString(temporalAddressFlag),
			viper.GetString(temporalNamespaceFlag),
			viper.GetString(temporalSSLClientCertFlag),
			viper.GetString(temporalSSLClientKeyFlag),
			viper.GetBool(temporalInitSearchAttributes),
		),
		bunconnect.Module(*connectionOptions),
		publish.CLIPublisherModule("orchestration"),
		auth.CLIAuthModule(),
		licence.CLIModule(ServiceName),
		workflow.NewModule(viper.GetString(temporalTaskQueueFlag)),
		triggers.NewModule(viper.GetString(temporalTaskQueueFlag)),
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
