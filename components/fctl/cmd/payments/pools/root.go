package pools

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewPoolsCommand() *cobra.Command {
	return fctl.NewCommand("pools",
		fctl.WithAliases("p"),
		fctl.WithShortDescription("Pools management"),
		fctl.WithChildCommands(
			NewListCommand(),
			NewCreateCommand(),
			NewShowCommand(),
			NewDeleteCommand(),
			NewBalancesCommand(),
			NewAddAccountCommand(),
			NewRemoveAccountCommand(),
		),
	)
}
