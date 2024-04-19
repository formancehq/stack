package cmd

import (
	"fmt"
	"os"

	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/licence"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
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

	otlpmetrics.InitOTLPMetricsFlags(cmd.PersistentFlags())
	otlptraces.InitOTLPTracesFlags(cmd.PersistentFlags())
	auth.InitAuthFlags(cmd.PersistentFlags())
	bunconnect.InitFlags(cmd.PersistentFlags())
	iam.InitFlags(cmd.PersistentFlags())
	service.BindFlags(cmd)
	licence.InitCLIFlags(cmd)

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
