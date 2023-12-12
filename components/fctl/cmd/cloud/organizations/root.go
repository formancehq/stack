package organizations

import (
	"github.com/formancehq/fctl/cmd/cloud/organizations/invitations"
	"github.com/formancehq/fctl/cmd/cloud/organizations/users"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewStackCommand("organizations",
		fctl.WithAliases("org", "o"),
		fctl.WithShortDescription("Organizations management"),
		fctl.WithChildCommands(
			NewListCommand(),
			NewCreateCommand(),
			NewDeleteCommand(),
			users.NewCommand(),
			invitations.NewCommand(),
		),
	)
}
