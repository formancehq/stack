package users

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewUnlinkCommand() *cobra.Command {
	return fctl.NewCommand("unlink <user-id>",
		fctl.WithAliases("u", "un"),
		fctl.WithShortDescription("unlink user from organization"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}

			apiClient, err := fctl.NewMembershipClient(cmd, cfg)
			if err != nil {
				return err
			}

			organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
			if err != nil {
				return err
			}

			_, err = apiClient.DefaultApi.UnlinkUserFromOrganization(cmd.Context(), organizationID, args[0]).Execute()
			if err != nil {
				return err
			}

			pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("User '%s' unlinked from organization '%s'", args[0], organizationID)

			return nil
		}),
	)
}
