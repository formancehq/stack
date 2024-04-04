package ledger

import (
	"github.com/formancehq/fctl/cmd/ledger/accounts"
	"github.com/formancehq/fctl/cmd/ledger/internal"
	"github.com/formancehq/fctl/cmd/ledger/transactions"
	"github.com/formancehq/fctl/cmd/ledger/volumes"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewStackCommand("ledger",
		fctl.WithAliases("l"),
		fctl.WithPersistentStringFlag(internal.LedgerFlag, "default", "Specific ledger"),
		fctl.WithShortDescription("Ledger management"),
		fctl.WithChildCommands(
			NewCreateCommand(),
			NewSendCommand(),
			NewStatsCommand(),
			NewServerInfoCommand(),
			NewListCommand(),
			NewSetMetadataCommand(),
			NewDeleteMetadataCommand(),
			transactions.NewLedgerTransactionsCommand(),
			accounts.NewLedgerAccountsCommand(),
			volumes.NewLedgerVolumesCommand(),
		),
		fctl.WithPersistentPreRunE(func(cmd *cobra.Command, args []string) error {
			return fctl.NewStackStore(cmd)
		}),
	)
}
