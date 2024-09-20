package views

import (
	"fmt"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func DisplayAtlarConfig(cmd *cobra.Command, connectorConfig *shared.ConnectorConfigResponse) error {
	config := connectorConfig.Data.AtlarConfig

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Name:"), config.Name})
	tableData = append(tableData, []string{pterm.LightCyan("AccessKey:"), config.AccessKey})
	tableData = append(tableData, []string{pterm.LightCyan("Secret:"), config.Secret})
	tableData = append(tableData, []string{pterm.LightCyan("BaseUrl:"), func() string {
		if config.BaseURL == nil {
			return ""
		}
		return *config.BaseURL
	}()})
	tableData = append(tableData, []string{pterm.LightCyan("PageSize:"), func() string {
		if config.PageSize == nil {
			return ""
		}
		return fmt.Sprintf("%d", *config.PageSize)
	}()})
	tableData = append(tableData, []string{pterm.LightCyan("Polling Period:"), func() string {
		if config.PollingPeriod == nil {
			return ""
		}
		return *config.PollingPeriod
	}()})
	tableData = append(tableData, []string{pterm.LightCyan("Transfer Initiation Status Polling Period:"), func() string {
		if config.TransferInitiationStatusPollingPeriod == nil {
			return ""
		}
		return *config.TransferInitiationStatusPollingPeriod
	}()})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}
	return nil
}
