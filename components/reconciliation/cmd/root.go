package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
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
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
	}

	cmd.PersistentFlags().BoolP(service.DebugFlag, "d", false, "Debug mode")
	cmd.PersistentFlags().String(stackURLFlag, "", "Stack url")
	cmd.PersistentFlags().String(stackClientIDFlag, "", "Stack client ID")
	cmd.PersistentFlags().String(stackClientSecretFlag, "", "Stack client secret")

	serveCmd := newServeCommand(Version)
	cmd.AddCommand(serveCmd)
	versionCmd := newVersionCommand()
	cmd.AddCommand(versionCmd)
	migrate := newMigrate()
	cmd.AddCommand(migrate)

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
		otlptraces.CLITracesModule(viper.GetViper()),
		// storage.NewModule(viper.GetString(postgresDSNFlag), viper.GetBool(service.DebugFlag), output),
	)
}
