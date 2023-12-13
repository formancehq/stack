package stacks

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewMembershipCommand("stacks",
		fctl.WithAliases("s"),
		fctl.WithShortDescription("Stack users management within an organization"),
		fctl.WithChildCommands(
			NewUpsertStackAccessRolesCommand(),
			NewListStackAccessRolesCommand(),
		),
	)
}
