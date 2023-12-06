package instances

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewStackCommand("instances",
		fctl.WithAliases("ins", "i"),
		fctl.WithShortDescription("Instances management"),
		fctl.WithChildCommands(
			NewListCommand(),
			NewShowCommand(),
			NewDescribeCommand(),
			NewSendEventCommand(),
			NewStopCommand(),
		),
	)
}
