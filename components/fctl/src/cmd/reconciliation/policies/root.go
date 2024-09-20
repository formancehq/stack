package policies

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewPoliciesCommand() *cobra.Command {
	return fctl.NewCommand("policies",
		fctl.WithAliases("p"),
		fctl.WithShortDescription("Policies management"),
		fctl.WithChildCommands(
			NewListCommand(),
			NewShowCommand(),
			NewCreateCommand(),
			NewDeleteCommand(),
			NewReconciliationCommand(),
		),
	)
}
