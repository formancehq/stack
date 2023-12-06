package workflows

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewStackCommand("workflows",
		fctl.WithAliases("w", "work"),
		fctl.WithShortDescription("Workflows management"),
		fctl.WithChildCommands(
			NewListCommand(),
			NewCreateCommand(),
			NewRunCommand(),
			NewShowCommand(),
			NewDeleteCommand(),
		),
	)
}
