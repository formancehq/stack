package ledger

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewServerInfoCommand() *cobra.Command {
	return fctl.NewCommand("server-infos",
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithAliases("si"),
		fctl.WithShortDescription("Read server info"),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}

			organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
			if err != nil {
				return err
			}

			stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
			if err != nil {
				return err
			}

			ledgerClient, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return err
			}

			response, err := ledgerClient.Ledger.GetInfo(cmd.Context())
			if err != nil {
				return err
			}

			infoResponse := response.ConfigInfoResponse
			tableData := pterm.TableData{}
			tableData = append(tableData, []string{pterm.LightCyan("Server"), fmt.Sprint(infoResponse.Server)})
			tableData = append(tableData, []string{pterm.LightCyan("Version"), fmt.Sprint(infoResponse.Version)})
			tableData = append(tableData, []string{pterm.LightCyan("Storage driver"), fmt.Sprint(infoResponse.Config.Storage.Driver)})

			if err := pterm.DefaultTable.
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render(); err != nil {
				return err
			}

			fctl.BasicTextCyan.WithWriter(cmd.OutOrStdout()).Printfln("Ledgers :")
			if err := pterm.DefaultBulletList.
				WithWriter(cmd.OutOrStdout()).
				WithItems(fctl.Map(infoResponse.Config.Storage.Ledgers, func(ledger string) pterm.BulletListItem {
					return pterm.BulletListItem{
						Text:        ledger,
						TextStyle:   pterm.NewStyle(pterm.FgDefault),
						BulletStyle: pterm.NewStyle(pterm.FgLightCyan),
					}
				})).
				Render(); err != nil {
				return err
			}

			return nil
		}),
	)
}
