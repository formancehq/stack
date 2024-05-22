package cmd

import (
	"fmt"
	"os"

	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/bun/bunmigrate"
	"github.com/formancehq/stack/libs/go-libs/licence"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"go.uber.org/fx"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ServiceName = "reconciliation"
	Version     = "develop"
	BuildDate   = "-"
	Commit      = "-"
)

const (
	stackURLFlag          = "stack-url"
	stackClientIDFlag     = "stack-client-id"
	stackClientSecretFlag = "stack-client-secret"
	listenFlag            = "listen"
	topicsFlag            = "topics"
	workerFlag            = "worker"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
	}

	cobra.EnableTraverseRunHooks = true

	cmd.PersistentFlags().String(stackURLFlag, "", "Stack url")
	cmd.PersistentFlags().String(stackClientIDFlag, "", "Stack client ID")
	cmd.PersistentFlags().String(stackClientSecretFlag, "", "Stack client secret")
	cmd.PersistentFlags().StringSlice(topicsFlag, []string{}, "Topics to listen")

	otlpmetrics.InitOTLPMetricsFlags(cmd.PersistentFlags())
	otlptraces.InitOTLPTracesFlags(cmd.PersistentFlags())
	auth.InitAuthFlags(cmd.PersistentFlags())
	bunconnect.InitFlags(cmd.PersistentFlags())
	iam.InitFlags(cmd.PersistentFlags())
	service.BindFlags(cmd)
	licence.InitCLIFlags(cmd)
	publish.InitCLIFlags(cmd)

	serveCmd := newServeCommand(Version)
	addAutoMigrateCommand(serveCmd)
	cmd.AddCommand(serveCmd)
	versionCmd := newVersionCommand()
	cmd.AddCommand(versionCmd)
	migrate := newMigrate()
	cmd.AddCommand(migrate)
	workerCmd := newWorkerCommand()
	cmd.AddCommand(workerCmd)
	return cmd
}

func commonOptions(cmd *cobra.Command) (fx.Option, error) {
	databaseOptions, err := prepareDatabaseOptions(cmd)
	if err != nil {
		return nil, err
	}

	options := make([]fx.Option, 0)
	options = append(options, databaseOptions)

	options = append(options,
		otlptraces.CLITracesModule(),
		otlpmetrics.CLIMetricsModule(),
		auth.CLIAuthModule(),
		licence.CLIModule(ServiceName),
	)

	return fx.Options(options...), nil
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

func addAutoMigrateCommand(cmd *cobra.Command) {
	cmd.Flags().Bool(autoMigrateFlag, false, "Auto migrate database")
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if viper.GetBool(autoMigrateFlag) {
			return bunmigrate.Run(cmd, args, Migrate)
		}
		return nil
	}
}
