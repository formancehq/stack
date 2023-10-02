package payments

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewPaymentsCommand() *cobra.Command {
	return fctl.NewCommand("payments",
		fctl.WithAliases("p"),
		fctl.WithShortDescription("Payments management"),
		fctl.WithChildCommands(
			NewListCommand(),
		),
	)
}
