package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/formancehq/orchestration/internal/storage"
	"github.com/formancehq/orchestration/internal/temporal"
	_ "github.com/formancehq/orchestration/internal/workflow/stages/all"
	"github.com/formancehq/stack/libs/go-libs/otlp"
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
	stackURLFlag              = "stack-url"
	stackClientIDFlag         = "stack-client-id"
	stackClientSecretFlag     = "stack-client-secret"
	temporalAddressFlag       = "temporal-address"
	temporalNamespaceFlag     = "temporal-namespace"
	temporalSSLClientKeyFlag  = "temporal-ssl-client-key"
	temporalSSLClientCertFlag = "temporal-ssl-client-cert"
	temporalTaskQueueFlag     = "temporal-task-queue"
	listenFlag                = "listen"
	postgresDSNFlag           = "postgres-dsn"
	workerFlag                = "worker"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
	}
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmd.PersistentFlags().BoolP(service.DebugFlag, "d", false, "Debug mode")
	cmd.PersistentFlags().String(stackURLFlag, "", "Stack url")
	cmd.PersistentFlags().String(stackClientIDFlag, "", "Stack client ID")
	cmd.PersistentFlags().String(stackClientSecretFlag, "", "Stack client secret")
	cmd.PersistentFlags().String(temporalAddressFlag, "", "Temporal server address")
	cmd.PersistentFlags().String(temporalNamespaceFlag, "default", "Temporal namespace")
	cmd.PersistentFlags().String(temporalSSLClientKeyFlag, "", "Temporal client key")
	cmd.PersistentFlags().String(temporalSSLClientCertFlag, "", "Temporal client cert")
	cmd.PersistentFlags().String(temporalTaskQueueFlag, "default", "Temporal task queue name")
	cmd.PersistentFlags().String(postgresDSNFlag, "", "Postgres address")
	cmd.AddCommand(newServeCommand(), newWorkerCommand(), newVersionCommand())

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

func commonOptions(output io.Writer) fx.Option {
	return fx.Options(
		otlp.LoadResource(viper.GetString(otlp.OtelServiceName), viper.GetStringSlice(otlp.OtelResourceAttributes)),
		otlptraces.CLITracesModule(viper.GetViper()),
		temporal.NewClientModule(
			viper.GetString(temporalAddressFlag),
			viper.GetString(temporalNamespaceFlag),
			viper.GetString(temporalSSLClientCertFlag),
			viper.GetString(temporalSSLClientKeyFlag),
		),
		storage.NewModule(viper.GetString(postgresDSNFlag), viper.GetBool(service.DebugFlag), output),
	)
}
