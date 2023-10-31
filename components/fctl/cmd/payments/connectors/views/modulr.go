package views

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func DisplayModulrConfig(cmd *cobra.Command, connectorConfig *shared.ConnectorConfigResponse) error {
	config := connectorConfig.Data.ModulrConfig

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("API key:"), config.APIKey})
	tableData = append(tableData, []string{pterm.LightCyan("API secret:"), config.APISecret})
	tableData = append(tableData, []string{pterm.LightCyan("Endpoint:"), func() string {
		if config.Endpoint == nil {
			return ""
		}
		return *config.Endpoint
	}()})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}
	return nil
}
