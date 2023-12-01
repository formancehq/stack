package views

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func DisplayBankingCircleConfig(cmd *cobra.Command, connectorConfig *shared.ConnectorConfigResponse) error {
	config := connectorConfig.Data.(map[string]interface{})

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Name:"), func() string {
		name, ok := config["name"].(string)
		if !ok {
			return ""
		}
		return name
	}()})
	tableData = append(tableData, []string{pterm.LightCyan("Username:"), func() string {
		username, ok := config["username"].(string)
		if !ok {
			return ""
		}
		return username
	}()})
	tableData = append(tableData, []string{pterm.LightCyan("Password:"), func() string {
		password, ok := config["password"].(string)
		if !ok {
			return ""
		}
		return password
	}()})
	tableData = append(tableData, []string{pterm.LightCyan("Endpoint:"), func() string {
		endpoint, ok := config["endpoint"].(string)
		if !ok {
			return ""
		}
		return endpoint
	}()})
	tableData = append(tableData, []string{pterm.LightCyan("Authorization endpoint:"), func() string {
		authorizationEndpoint, ok := config["authorizationEndpoint"].(string)
		if !ok {
			return ""
		}
		return authorizationEndpoint
	}()})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}
	return nil
}
