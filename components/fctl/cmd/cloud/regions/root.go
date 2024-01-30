package regions

import (
	"github.com/formancehq/fctl/cmd/cloud/regions/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewMembershipCommand("regions",
		fctl.WithAliases("region", "reg"),
		fctl.WithShortDescription("Regions management"),
		fctl.WithChildCommands(
			NewListCommand(),
			NewShowCommand(),
			NewCreateCommand(),
			NewDeleteCommand(),
			versions.NewCommand(),
		),
	)
}
