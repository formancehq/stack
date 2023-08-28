package ledger

import (
	"github.com/formancehq/fctl/cmd/ledger/accounts"
	"github.com/formancehq/fctl/cmd/ledger/transactions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("ledger",
		fctl.WithAliases("l"),
		fctl.WithShortDescription("Ledger management"),
		fctl.WithChildCommands(
			NewBalancesCommand(),
			NewSendCommand(),
			NewStatsCommand(),
			NewServerInfoCommand(),
			NewListCommand(),
			transactions.NewCommand(),
			accounts.NewCommand(),
		),
		fctl.WithCommandScopesFlags(fctl.Organization, fctl.Stack, fctl.Ledger),
	)
}
