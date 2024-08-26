package cmd

import (
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
	cmd := &cobra.Command{}

	cobra.EnableTraverseRunHooks = true

	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd := newServeCommand()
	cmd.AddCommand(serverCmd)
	return cmd
}

func Execute() {
	service.Execute(NewRootCommand())
}
