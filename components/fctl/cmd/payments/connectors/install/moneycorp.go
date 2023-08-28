package install

import (
	"flag"
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/connectors/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useMoneycorpConnector    = internal.MoneycorpConnector + " <clientID> <apiKey>"
	shortMoneycorpConnector  = "Install Moneycorp connector"
	defaultEndpointMoneyCorp = "https://sandbox-corpapi.moneycorp.com"
)

type MoneycorpStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
}

func NewMoneycorpStore() *MoneycorpStore {
	return &MoneycorpStore{
		Success:       false,
		ConnectorName: internal.MoneycorpConnector,
	}
}

func NewMoneycorpConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useMoneycorpConnector, flag.ExitOnError)
	flags.String(EndpointFlag, defaultEndpointMoneyCorp, "API endpoint")
	flags.String(PollingPeriodFlag, DefaultPollingPeriod, "Polling duration")
	return fctl.NewControllerConfig(
		useMoneycorpConnector,
		shortMoneycorpConnector,
		shortMoneycorpConnector,
		[]string{},
		flags,
		fctl.Organization, fctl.Stack,
	)

}

type MoneycorpController struct {
	store  *MoneycorpStore
	config *fctl.ControllerConfig
}

func NewMoneycorpController(config *fctl.ControllerConfig) *MoneycorpController {
	return &MoneycorpController{
		store:  NewMoneycorpStore(),
		config: config,
	}
}

func (c *MoneycorpController) GetStore() *MoneycorpStore {
	return c.store
}
func (c *MoneycorpController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *MoneycorpController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()
	out := c.config.GetOut()
	if len(args) < 2 {
		return nil, errors.New("missing required arguments")
	}

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

	if !fctl.CheckStackApprobation(flags, stack, "You are about to install connector '%s'", internal.MoneycorpConnector) {
		return nil, fctl.ErrMissingApproval
	}

	paymentsClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	request := operations.InstallConnectorRequest{
		Connector: shared.ConnectorMoneycorp,
		RequestBody: shared.MoneycorpConfig{
			ClientID:      args[0],
			APIKey:        args[1],
			Endpoint:      fctl.GetString(flags, EndpointFlag),
			PollingPeriod: fctl.Ptr(fctl.GetString(flags, PollingPeriodFlag)),
		},
	}
	response, err := paymentsClient.Payments.InstallConnector(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "installing connector")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Connector installed!")

	c.store.Success = true

	return c, nil
}

func (c *MoneycorpController) Render() error {

	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Connector %s installed!", c.store.ConnectorName)

	return nil
}

func NewMoneycorpCommand() *cobra.Command {
	config := NewMoneycorpConfig()
	c := NewMoneycorpController(config)

	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*MoneycorpStore](c),
	)
}
