package occurrences

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewStackCommand("occurrences",
		fctl.WithAliases("oc", "o"),
		fctl.WithShortDescription("Triggers occurrences management"),
		fctl.WithChildCommands(
			NewListCommand(),
		),
	)
}
