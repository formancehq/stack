package views

import (
	"errors"

	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func DisplayStripeConfig(cmd *cobra.Command, connectorConfig *shared.ConnectorConfigResponse) error {
	config, ok := connectorConfig.Data.(*shared.StripeConfig)
	if !ok {
		return errors.New("invalid stripe connector config")
	}

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("API key:"), config.APIKey})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}
	return nil
}
