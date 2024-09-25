package views

import (
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func DisplayAdyenConfig(cmd *cobra.Command, connectorConfig *shared.ConnectorConfigResponse) error {
	config := connectorConfig.Data.AdyenConfig

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Name:"), config.Name})
	tableData = append(tableData, []string{pterm.LightCyan("ApiKey:"), config.APIKey})
	tableData = append(tableData, []string{pterm.LightCyan("HMACKey:"), config.HmacKey})
	tableData = append(tableData, []string{pterm.LightCyan("LiveEndpointPrefix:"), func() string {
		if config.LiveEndpointPrefix == nil {
			return ""
		}

		return *config.LiveEndpointPrefix
	}()})
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
