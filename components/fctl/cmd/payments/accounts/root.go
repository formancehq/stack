package accounts

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewAccountsCommand() *cobra.Command {
	return fctl.NewCommand("accounts",
		fctl.WithAliases("acc", "a", "ac", "account"),
		fctl.WithShortDescription("Accounts management"),
		fctl.WithChildCommands(
			NewListCommand(),
			NewShowCommand(),
			NewListBalanceCommand(),
		),
	)
}
