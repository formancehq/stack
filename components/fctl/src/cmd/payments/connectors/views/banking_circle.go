package views

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func DisplayBankingCircleConfig(cmd *cobra.Command, connectorConfig *shared.ConnectorConfigResponse) error {
	config := connectorConfig.Data.BankingCircleConfig

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Name:"), config.Name})
	tableData = append(tableData, []string{pterm.LightCyan("Username:"), config.Username})
	tableData = append(tableData, []string{pterm.LightCyan("Password:"), config.Password})
	tableData = append(tableData, []string{pterm.LightCyan("Endpoint:"), config.Endpoint})
	tableData = append(tableData, []string{pterm.LightCyan("Authorization endpoint:"), config.AuthorizationEndpoint})
	tableData = append(tableData, []string{pterm.LightCyan("Polling Period:"), func() string {
		if config.PollingPeriod == nil {
			return ""
		}
		return *config.PollingPeriod
	}()})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}
	return nil
}
