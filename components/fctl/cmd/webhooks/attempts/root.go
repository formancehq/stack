package attempts

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewAttemptsCommand() *cobra.Command {
	return fctl.NewCommand("attempts",
		fctl.WithAliases("att", "ats"),
		fctl.WithShortDescription("Attempts Management"),
		fctl.WithChildCommands(
			NewListWaitingCommand(),
			NewListAbortedCommand(),
			NewAbortCommand(),
			NewRetryCommand(),
			NewRetryAllCommand(),
		),
		fctl.WithPersistentPreRunE(func(cmd *cobra.Command, args []string) error {
			return fctl.NewStackStore(cmd)
		}),
	)
}
