package ledger

import (
	"github.com/formancehq/fctl/cmd/ledger/accounts"
	"github.com/formancehq/fctl/cmd/ledger/internal"
	"github.com/formancehq/fctl/cmd/ledger/store"
	"github.com/formancehq/fctl/cmd/ledger/transactions"
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
			transactions.NewLedgerTransactionsCommand(),
			accounts.NewLedgerAccountsCommand(),
		),
		fctl.WithPersistentPreRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}
			apiClient, err := fctl.NewMembershipClient(cmd, cfg)
			if err != nil {
				return err
			}
			organizationID, err := fctl.ResolveOrganizationID(cmd, cfg, apiClient.DefaultApi)
			if err != nil {
				return err
			}

			stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
			if err != nil {
				return err
			}

			stackClient, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return err
			}
			cmd.SetContext(store.ContextWithStore(cmd.Context(), store.LedgerNode(cfg, stack, organizationID, stackClient)))
			return nil
		}),
	)
}
