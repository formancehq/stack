package payments

import (
	"github.com/formancehq/fctl/cmd/payments/accounts"
	"github.com/formancehq/fctl/cmd/payments/bankaccounts"
	"github.com/formancehq/fctl/cmd/payments/connectors"
	"github.com/formancehq/fctl/cmd/payments/payments"
	"github.com/formancehq/fctl/cmd/payments/pools"
	"github.com/formancehq/fctl/cmd/payments/store"
	"github.com/formancehq/fctl/cmd/payments/transferinitiation"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewStackCommand("payments",
		fctl.WithShortDescription("Payments management"),
		fctl.WithChildCommands(
			connectors.NewConnectorsCommand(),
			payments.NewPaymentsCommand(),
			transferinitiation.NewTransferInitiationCommand(),
			bankaccounts.NewBankAccountsCommand(),
			accounts.NewAccountsCommand(),
			pools.NewPoolsCommand(),
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
			cmd.SetContext(store.ContextWithStore(cmd.Context(), store.PaymentsNode(cfg, stack, organizationID, stackClient)))
			return nil
		}),
	)
}
