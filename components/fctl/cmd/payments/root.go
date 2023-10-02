package payments

import (
	"github.com/formancehq/fctl/cmd/payments/accounts"
	"github.com/formancehq/fctl/cmd/payments/bankaccounts"
	"github.com/formancehq/fctl/cmd/payments/connectors"
	"github.com/formancehq/fctl/cmd/payments/payments"
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
		),
	)
}
