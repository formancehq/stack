package payments

import (
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
			payments.NewListPaymentsCommand(),
			transferinitiation.NewShowCommand(),
			transferinitiation.NewListTransferInitiationCommand(),
			transferinitiation.NewCreateCommand(),
			transferinitiation.NewUpdateStatusCommand(),
			transferinitiation.NewDeleteCommand(),
		),
	)
}
