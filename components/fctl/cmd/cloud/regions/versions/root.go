package versions

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewMembershipCommand("versions",
		fctl.WithAliases("v"),
		fctl.WithShortDescription("Versions management by region"),
		fctl.WithChildCommands(
			NewListCommand(),
		),
	)
}
