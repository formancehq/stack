package regions

import (
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("show <region-id>",
		fctl.WithAliases("sh", "s"),
		fctl.WithShortDescription("Show region details"),
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

			response, _, err := apiClient.DefaultApi.GetRegion(cmd.Context(), organizationID, args[0]).Execute()
			if err != nil {
				return err
			}

			fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
			tableData := pterm.TableData{}
			tableData = append(tableData, []string{pterm.LightCyan("ID"), response.Data.Id})
			tableData = append(tableData, []string{pterm.LightCyan("Name"), response.Data.Name})
			tableData = append(tableData, []string{pterm.LightCyan("Base URL"), response.Data.BaseUrl})
			tableData = append(tableData, []string{pterm.LightCyan("Active: "), fctl.BoolToString(response.Data.Active)})
			tableData = append(tableData, []string{pterm.LightCyan("Public: "), fctl.BoolToString(response.Data.Public)})
			if response.Data.Creator != nil {
				tableData = append(tableData, []string{pterm.LightCyan("Creator"), response.Data.Creator.Email})
			}
			if response.Data.LastPing != nil {
				tableData = append(tableData, []string{pterm.LightCyan("Base URL"), response.Data.LastPing.Format(time.RFC3339)})
			}

			return pterm.DefaultTable.
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render()
		}),
	)
}
