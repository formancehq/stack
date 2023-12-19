package triggers

import (
	"github.com/formancehq/fctl/cmd/orchestration/triggers/occurrences"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewStackCommand("triggers",
		fctl.WithAliases("trig", "t"),
		fctl.WithShortDescription("Triggers management"),
		fctl.WithChildCommands(
			NewListCommand(),
			NewShowCommand(),
			NewDeleteCommand(),
			NewCreateCommand(),
			NewTestCommand(),
			occurrences.NewCommand(),
		),
	)
}
