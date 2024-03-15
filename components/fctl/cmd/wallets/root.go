package wallets

import (
	"github.com/formancehq/fctl/cmd/wallets/balances"
	"github.com/formancehq/fctl/cmd/wallets/holds"
	"github.com/formancehq/fctl/cmd/wallets/transactions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewStackCommand("wallets",
		fctl.WithAliases("wal", "wa", "wallet"),
		fctl.WithShortDescription("Wallets management"),
		fctl.WithChildCommands(
			NewCreateCommand(),
			NewUpdateCommand(),
			NewListCommand(),
			NewShowCommand(),
			NewCreditWalletCommand(),
			NewDebitWalletCommand(),
			transactions.NewCommand(),
			holds.NewCommand(),
			balances.NewCommand(),
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
			cmd.SetContext(fctl.ContextWithStackStore(cmd.Context(), fctl.StackNode(cfg, stack, organizationID, stackClient)))
			return nil
		}),
	)
}
