package cloud

import (
	"github.com/formancehq/fctl/cmd/cloud/me"
	"github.com/formancehq/fctl/cmd/cloud/organizations"
	"github.com/formancehq/fctl/cmd/cloud/regions"
	"github.com/formancehq/fctl/cmd/cloud/store"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return fctl.NewCommand("cloud",
		fctl.WithAliases("c"),
		fctl.WithShortDescription("Cloud management"),
		fctl.WithChildCommands(
			organizations.NewCommand(),
			me.NewCommand(),
			regions.NewCommand(),
			NewGeneratePersonalTokenCommand(),
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
			cmd.SetContext(store.ContextWithStore(cmd.Context(), store.CloudNode(cfg, apiClient)))
			return nil
		}),
	)
}
