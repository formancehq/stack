package connectors

import (
	"errors"
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/connectors/internal"
	"github.com/formancehq/fctl/cmd/payments/connectors/views"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	connectorsAvailable = []string{internal.StripeConnector} //internal.ModulrConnector, internal.BankingCircleConnector, internal.CurrencyCloudConnector, internal.WiseConnector}
)

type PaymentsGetConfigStore struct {
	ConnectorConfig *shared.ConnectorConfigResponse `json:"connector_config"`
}
type PaymentsGetConfigController struct {
	store *PaymentsGetConfigStore
	args  []string
}

var _ fctl.Controller[*PaymentsGetConfigStore] = (*PaymentsGetConfigController)(nil)

func NewDefaultPaymentsGetConfigStore() *PaymentsGetConfigStore {
	return &PaymentsGetConfigStore{}
}

func NewPaymentsGetConfigController() *PaymentsGetConfigController {
	return &PaymentsGetConfigController{
		store: NewDefaultPaymentsGetConfigStore(),
	}
}

func NewGetConfigCommand() *cobra.Command {
	return fctl.NewCommand("get-config <connector-name>",
		fctl.WithAliases("getconfig", "getconf", "gc", "get", "g"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithValidArgs(connectorsAvailable...),
		fctl.WithShortDescription(fmt.Sprintf("Read a connector config (Connectors available: %s)", connectorsAvailable)),
		fctl.WithController[*PaymentsGetConfigStore](NewPaymentsGetConfigController()),
	)
}

func (c *PaymentsGetConfigController) GetStore() *PaymentsGetConfigStore {
	return c.store
}

func (c *PaymentsGetConfigController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	response, err := client.Payments.ReadConnectorConfig(cmd.Context(), operations.ReadConnectorConfigRequest{
		Connector: shared.Connector(args[0]),
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.args = args
	c.store.ConnectorConfig = response.ConnectorConfigResponse

	return c, err

}

// TODO: This need to use the ui.NewListModel
func (c *PaymentsGetConfigController) Render(cmd *cobra.Command, args []string) error {
	var err error

	switch c.args[0] {
	case internal.StripeConnector:
		err = views.DisplayStripeConfig(cmd, c.store.ConnectorConfig)
	case internal.ModulrConnector:
		err = views.DisplayModulrConfig(cmd, c.store.ConnectorConfig)
	case internal.BankingCircleConnector:
		err = views.DisplayBankingCircleConfig(cmd, c.store.ConnectorConfig)
	case internal.CurrencyCloudConnector:
		err = views.DisplayCurrencyCloudConfig(cmd, c.store.ConnectorConfig)
	case internal.WiseConnector:
		err = views.DisplayWiseConfig(cmd, c.store.ConnectorConfig)
	default:
		pterm.Error.WithWriter(cmd.OutOrStderr()).Printfln("Connection unknown.")
	}

	return err

}

func displayMangoPayConfig(cmd *cobra.Command, connectorConfig *shared.ConnectorConfigResponse) error {
	config, ok := connectorConfig.Data.(*shared.MangoPayConfig)
	if !ok {
		return errors.New("invalid currency cloud connector config")
	}

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("API key:"), config.APIKey})
	tableData = append(tableData, []string{pterm.LightCyan("Client ID:"), config.ClientID})
	tableData = append(tableData, []string{pterm.LightCyan("Endpoint:"), config.Endpoint})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}
	return nil
}

func displayMoneycorpConfig(cmd *cobra.Command, connectorConfig *shared.ConnectorConfigResponse) error {
	config, ok := connectorConfig.Data.(*shared.MoneycorpConfig)
	if !ok {
		return errors.New("invalid currency cloud connector config")
	}

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("API key:"), config.APIKey})
	tableData = append(tableData, []string{pterm.LightCyan("Client ID:"), config.ClientID})
	tableData = append(tableData, []string{pterm.LightCyan("Endpoint:"), config.Endpoint})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}
	return nil
}
