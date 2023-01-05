package wallets

import (
	"github.com/formancehq/fctl/cmd/wallets/holds"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("wallets",
		fctl.WithAliases("wal", "wa", "wallet"),
		fctl.WithShortDescription("Wallets management"),
		fctl.WithChildCommands(
			NewCreateCommand(),
			NewUpdateCommand(),
			NewListCommand(),
			NewShowCommand(),
			NewCreditWalletCommand(),
			NewDebitWalletCommand(),
			holds.NewCommand(),
		),
	)
}
