package regions

import (
	"errors"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("show",
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

			regionsResponse, _, err := apiClient.DefaultApi.ListRegions(cmd.Context()).Execute()
			if err != nil {
				return err
			}

			//TODO: Add GET /regions/<id> on fctl

			var region *membershipclient.Region
			for _, r := range regionsResponse.Data {
				if r.Id == args[0] {
					region = &r
				}
			}
			if region == nil {
				return errors.New("region not found")
			}

			tableData := pterm.TableData{}
			tableData = append(tableData, []string{pterm.LightCyan("ID"), region.Id})
			tableData = append(tableData, []string{pterm.LightCyan("BaseUrl"), region.BaseUrl})

			if err := pterm.DefaultTable.
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render(); err != nil {
				return err
			}

			fctl.Highlightln(cmd.OutOrStdout(), "Tags :")
			tableData = pterm.TableData{}
			for k, v := range region.Tags {
				tableData = append(tableData, []string{pterm.LightCyan(k), v})
			}

			return pterm.DefaultTable.
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render()
		}),
	)
}
