package modules

import (
	"errors"

	"github.com/formancehq/fctl/cmd/stack/store"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewMembershipCommand("modules",
		fctl.WithShortDescription("Manage your modules"),
		fctl.WithAliases("module", "mod"),
		fctl.WithPersistentStringFlag(stackFlag, "", "Stack"),
		fctl.WithPersistentPreRunE(func(cmd *cobra.Command, args []string) error {
			stack := fctl.GetString(cmd, stackFlag)
			if stack == "" {
				return errors.New("--stack=<stack-id> is required")
			}

			store := store.GetStore(cmd.Context())
			if err := store.CheckRegionCapability(string(membershipclient.MODULE_LIST), func(s []any) bool {
				return len(s) > 0
			})(cmd, args); err != nil {
				return err
			}

			if err := fctl.CheckMembershipCapabilities(membershipclient.MODULE_SELECTION)(cmd, args); err != nil {
				return err
			}
			return nil
		}),
		fctl.WithChildCommands(
			NewDisableCommand(),
			NewEnableCommand(),
			NewListCommand(),
		),
	)
}
