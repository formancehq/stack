package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version   = "develop"
	BuildDate = "-"
	Commit    = "-"
)

const (
	debugFlag = "debug"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
	}
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmd.PersistentFlags().BoolP(debugFlag, "d", false, "Debug mode")
	cmd.AddCommand(newServeCommand(), newVersionCommand())

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
