package connectors

import (
	"flag"
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
	connectorsAvailable  = []string{internal.StripeConnector} //internal.ModulrConnector, internal.BankingCircleConnector, internal.CurrencyCloudConnector, internal.WiseConnector}
	useGetConfig         = "get-config <connector-name>"
	descriptionGetConfig = fmt.Sprintf("Read a connector config (Connectors available: %s)", connectorsAvailable)
	shortGetConfig       = fmt.Sprintf("Read a connector config (Connectors available: %s)", connectorsAvailable)
)

type GetConfigStore struct {
	ConnectorConfig *shared.ConnectorConfigResponse `json:"connectorConfig"`
}

func NewGetConfigStore() *GetConfigStore {
	return &GetConfigStore{}
}
func NewGetConfigConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useGetConfig, flag.ExitOnError)

	return fctl.NewControllerConfig(
		useGetConfig,
		descriptionGetConfig,
		shortGetConfig,
		[]string{
			"getconfig", "getconf", "gc", "get", "g",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

var _ fctl.Controller[*GetConfigStore] = (*GetConfigController)(nil)

type GetConfigController struct {
	store  *GetConfigStore
	config *fctl.ControllerConfig
}

func NewGetConfigController(config *fctl.ControllerConfig) *GetConfigController {
	return &GetConfigController{
		store:  NewGetConfigStore(),
		config: config,
	}
}

func (c *GetConfigController) GetStore() *GetConfigStore {
	return c.store
}

func (c *GetConfigController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *GetConfigController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()
	out := c.config.GetOut()
	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	response, err := client.Payments.ReadConnectorConfig(ctx, operations.ReadConnectorConfigRequest{
		Connector: shared.Connector(args[0]),
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.ConnectorConfig = response.ConnectorConfigResponse

	return c, err

}

func (c *GetConfigController) Render() error {
	var err error

	out := c.config.GetOut()

	switch c.config.GetArgs()[0] {
	case internal.StripeConnector:
		err = views.DisplayStripeConfig(out, c.store.ConnectorConfig)
	case internal.ModulrConnector:
		err = views.DisplayModulrConfig(out, c.store.ConnectorConfig)
	case internal.BankingCircleConnector:
		err = views.DisplayBankingCircleConfig(out, c.store.ConnectorConfig)
	case internal.CurrencyCloudConnector:
		err = views.DisplayCurrencyCloudConfig(out, c.store.ConnectorConfig)
	case internal.WiseConnector:
		err = views.DisplayWiseConfig(out, c.store.ConnectorConfig)
	case internal.MoneycorpConnector:
		err = views.DisplayMoneycorpConfig(out, c.store.ConnectorConfig)
	case internal.MangoPayConnector:
		err = views.DisplayMangoPayConfig(out, c.store.ConnectorConfig)
	default:
		pterm.Error.WithWriter(out).Printfln("Connection unknown.")
	}

	return err

}

func NewGetConfigCommand() *cobra.Command {
	config := NewGetConfigConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithValidArgs(connectorsAvailable...),
		fctl.WithController[*GetConfigStore](NewGetConfigController(config)),
	)
}
