package organizations

import (
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List organizations"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}

			apiClient, err := fctl.NewMembershipClient(cmd, cfg)
			if err != nil {
				return err
			}

			organizations, _, err := apiClient.DefaultApi.ListOrganizationsExpanded(cmd.Context()).Execute()
			if err != nil {
				return err
			}

			currentProfile := fctl.GetCurrentProfile(cmd, cfg)
			claims, err := currentProfile.GetClaims()
			if err != nil {
				return err
			}

			tableData := fctl.Map(organizations.Data, func(o membershipclient.ListOrganizationExpandedResponseDataInner) []string {
				isMine := fctl.BoolToString(o.OwnerId == claims["sub"].(string))
				if isMine == "yes" {
					isMine = pterm.LightMagenta(isMine)
				} else {
					isMine = pterm.LightGreen(isMine)
				}
				return []string{
					o.Id, o.Name, o.OwnerId, o.Owner.Email, isMine,
				}
			})
			tableData = fctl.Prepend(tableData, []string{"ID", "Name", "Owner ID", "Owner email", "Is mine?"})
			return pterm.DefaultTable.
				WithHasHeader().
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render()
		}),
	)
}
