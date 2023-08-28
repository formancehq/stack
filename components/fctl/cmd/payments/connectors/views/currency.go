package views

import (
	"errors"
	"io"

	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
)

func DisplayCurrencyCloudConfig(writer io.Writer, connectorConfig *shared.ConnectorConfigResponse) error {
	config, ok := connectorConfig.Data.(*shared.CurrencyCloudConfig)
	if !ok {
		return errors.New("invalid currency cloud connector config")
	}

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("API key:"), config.APIKey})
	tableData = append(tableData, []string{pterm.LightCyan("Login ID:"), config.LoginID})
	tableData = append(tableData, []string{pterm.LightCyan("Endpoint:"), func() string {
		if config.Endpoint == nil {
			return ""
		}
		return *config.Endpoint
	}()})

	if err := pterm.DefaultTable.
		WithWriter(writer).
		WithData(tableData).
		Render(); err != nil {
		return err
	}
	return nil
}
