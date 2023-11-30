package views

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func DisplayMangopayConfig(cmd *cobra.Command, connectorConfig *shared.ConnectorConfigResponse) error {
	config := connectorConfig.Data.(map[string]interface{})

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Name:"), func() string {
		name, ok := config["name"].(string)
		if !ok {
			return ""
		}
		return name
	}()})
	tableData = append(tableData, []string{pterm.LightCyan("ClientID:"), func() string {
		clientID, ok := config["clientID"].(string)
		if !ok {
			return ""
		}
		return clientID
	}()})
	tableData = append(tableData, []string{pterm.LightCyan("API key:"), func() string {
		apiKey, ok := config["apiKey"].(string)
		if !ok {
			return ""
		}
		return apiKey
	}()})
	tableData = append(tableData, []string{pterm.LightCyan("Endpoint:"), func() string {
		endpoint, ok := config["endpoint"].(string)
		if !ok {
			return ""
		}
		return endpoint
	}()})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}
	return nil
}
