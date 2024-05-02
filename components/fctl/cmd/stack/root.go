package stack

import (
	"github.com/formancehq/fctl/cmd/stack/modules"
	"github.com/formancehq/fctl/cmd/stack/store"
	"github.com/formancehq/fctl/cmd/stack/users"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewMembershipCommand("stack",
		fctl.WithShortDescription("Manage your stack"),
		fctl.WithAliases("stack", "stacks", "st"),
		fctl.WithChildCommands(
			NewCreateCommand(),
			NewListCommand(),
			NewDeleteCommand(),
			NewShowCommand(),
			NewDisableCommand(),
			NewEnableCommand(),
			NewRestoreStackCommand(),
			NewUpdateCommand(),
			NewUpgradeCommand(),
			NewHistoryCommand(),
			users.NewCommand(),
			modules.NewCommand(),
		),
		fctl.WithPersistentPreRunE(func(cmd *cobra.Command, args []string) error {
			return store.NewMembershipStackStore(cmd)
		}),
	)
}
