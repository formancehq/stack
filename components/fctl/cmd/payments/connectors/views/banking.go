package views

import (
	"errors"

	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func DisplayBankingCircleConfig(cmd *cobra.Command, connectorConfig *shared.ConnectorConfigResponse) error {
	config, ok := connectorConfig.Data.(*shared.BankingCircleConfig)
	if !ok {
		return errors.New("invalid banking circle connector config")
	}

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Username:"), config.Username})
	tableData = append(tableData, []string{pterm.LightCyan("Password:"), config.Password})
	tableData = append(tableData, []string{pterm.LightCyan("Endpoint:"), config.Endpoint})
	tableData = append(tableData, []string{pterm.LightCyan("Authorization endpoint:"), config.AuthorizationEndpoint})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}
	return nil
}
