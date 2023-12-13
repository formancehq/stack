package roles

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewMembershipCommand("roles",
		fctl.WithAliases("s"),
		fctl.WithShortDescription("Stack users management within an organization"),
		fctl.WithChildCommands(
			NewUpsertCommand(),
			NewListCommand(),
			NewDeleteCommand(),
		),
	)
}
