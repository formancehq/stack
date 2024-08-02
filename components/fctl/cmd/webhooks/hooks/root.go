package hooks

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewHooksCommand() *cobra.Command {
	return fctl.NewCommand("hooks",
		fctl.WithAliases("hks", "hk"),
		fctl.WithShortDescription("Hooks Management"),
		fctl.WithChildCommands(
			NewCreateCommand(),
			NewListCommand(),
			NewDeactivateCommand(),
			NewActivateCommand(),
			NewDeleteCommand(),
			NewChangeSecretCommand(),
			NewChangeEndpointCommand(),
		),
		fctl.WithPersistentPreRunE(func(cmd *cobra.Command, args []string) error {
			return fctl.NewStackStore(cmd)
		}),
	)
}
