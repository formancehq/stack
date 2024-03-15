package stack

import (
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
			users.NewCommand(),
		),
		fctl.WithPersistentPreRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}
			apiClient, err := fctl.NewMembershipClient(cmd, cfg)
			if err != nil {
				return err
			}

			organization, err := fctl.ResolveOrganizationID(cmd, cfg, apiClient.DefaultApi)
			if err != nil {
				return err
			}

			mbStore := &fctl.MembershipStore{
				Config:           cfg,
				MembershipClient: apiClient,
			}

			cmd.SetContext(store.ContextWithStore(cmd.Context(), store.StackNode(mbStore, organization)))
			return nil
		}),
	)
}
