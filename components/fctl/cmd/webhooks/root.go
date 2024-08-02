package webhooks

import (
	"github.com/formancehq/fctl/cmd/webhooks/attempts"
	"github.com/formancehq/fctl/cmd/webhooks/hooks"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewStackCommand("webhooks",
		fctl.WithAliases("web", "wh"),
		fctl.WithShortDescription("Webhooks management"),
		fctl.WithChildCommands(
			hooks.NewHooksCommand(),
			attempts.NewAttemptsCommand(),
		),
		fctl.WithPersistentPreRunE(func(cmd *cobra.Command, args []string) error {
			return fctl.NewStackStore(cmd)
		}),
	)
}
