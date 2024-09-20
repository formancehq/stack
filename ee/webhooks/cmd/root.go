package cmd

import (
	"github.com/formancehq/go-libs/service"

	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use: "webhooks",
	}

	root.AddCommand(newServeCommand())
	root.AddCommand(newWorkerCommand())
	root.AddCommand(newVersionCommand())
	root.AddCommand(newMigrateCommand())

	return root
}

func Execute() {
	service.Execute(NewRootCommand())
}
