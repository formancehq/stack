package invitations

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewDeleteCommand() *cobra.Command {
	return fctl.NewCommand("delete <id>",
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithShortDescription("Delete an invitation"),
		fctl.WithAliases("del"),
		fctl.WithConfirmFlag(),
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

			if !fctl.CheckOrganizationApprobation(cmd, "You are about to delete an invitation") {
				return fctl.ErrMissingApproval
			}

			_, err = apiClient.DefaultApi.
				DeleteInvitation(cmd.Context(), organizationID, args[0]).
				Execute()
			if err != nil {
				return err
			}

			pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Invitation deleted")
			return nil
		}),
	)
}
