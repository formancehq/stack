package views

import (
	"errors"
	"io"

	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
)

func DisplayWiseConfig(writer io.Writer, connectorConfig *shared.ConnectorConfigResponse) error {
	config, ok := connectorConfig.Data.(*shared.WiseConfig)
	if !ok {
		return errors.New("invalid wise connector config")
	}

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("API key:"), config.APIKey})

	if err := pterm.DefaultTable.
		WithWriter(writer).
		WithData(tableData).
		Render(); err != nil {
		return err
	}
	return nil
}
