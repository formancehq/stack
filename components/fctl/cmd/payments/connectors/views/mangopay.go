package views

import (
	"errors"
	"io"

	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
)

func DisplayMangoPayConfig(out io.Writer, connectorConfig *shared.ConnectorConfigResponse) error {
	config, ok := connectorConfig.Data.(*shared.MangoPayConfig)
	if !ok {
		return errors.New("invalid currency cloud connector config")
	}

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("API key:"), config.APIKey})
	tableData = append(tableData, []string{pterm.LightCyan("Client ID:"), config.ClientID})
	tableData = append(tableData, []string{pterm.LightCyan("Endpoint:"), config.Endpoint})

	if err := pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}
	return nil
}
