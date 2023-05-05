package regions

import (
	"time"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List users"),
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

			regionsResponse, _, err := apiClient.DefaultApi.ListRegions(cmd.Context(), organizationID).Execute()
			if err != nil {
				return err
			}

			tableData := fctl.Map(regionsResponse.Data, func(i membershipclient.AnyRegion) []string {
				return []string{
					i.Id,
					i.Name,
					i.BaseUrl,
					fctl.BoolToString(i.Public),
					fctl.BoolToString(i.Active),
					func() string {
						if i.LastPing != nil {
							return i.LastPing.Format(time.RFC3339)
						}
						return ""
					}(),
					func() string {
						if i.Creator != nil {
							return i.Creator.Email
						}
						return ""
					}(),
				}
			})
			tableData = fctl.Prepend(tableData, []string{"ID", "Name", "Base url", "Public", "Active", "Last ping", "Creator"})
			return pterm.DefaultTable.
				WithHasHeader().
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render()
		}),
	)
}
