package bankaccounts

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewBankAccountsCommand() *cobra.Command {
	return fctl.NewCommand("bank_accounts",
		fctl.WithAliases("acc", "a", "ac", "account"),
		fctl.WithShortDescription("Bank Accounts management"),
		fctl.WithChildCommands(
			NewCreateCommand(),
			NewShowCommand(),
			NewListCommand(),
		),
	)
}
