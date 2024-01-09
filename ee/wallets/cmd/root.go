package cmd

import (
	"fmt"
	"os"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
)

var (
	ServiceName = "wallets"
	Version     = "develop"
	BuildDate   = "-"
	Commit      = "-"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := bindFlagsToViper(cmd); err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmd.PersistentFlags().BoolP(service.DebugFlag, "d", false, "Debug mode")
	serverCmd := newServeCommand()
	auth.InitAuthFlags(serverCmd.Flags())
	otlptraces.InitOTLPTracesFlags(serverCmd.Flags())
	cmd.AddCommand(serverCmd)
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
