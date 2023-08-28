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
	defaultEndpointMangoPay = "https://api.sandbox.mangopay.com"
	useMangoPay             = internal.MangoPayConnector + " <clientID> <apiKey>"
	shortMangoPay           = "Install MangoPay connector"
)

type MangoPayStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
}

func NewMangoPayStore() *MangoPayStore {
	return &MangoPayStore{
		Success:       false,
		ConnectorName: internal.MangoPayConnector,
	}
}
func NewMangoPayConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useMangoPay, flag.ExitOnError)
	flags.String(EndpointFlag, defaultEndpointMangoPay, "API endpoint")
	flags.String(PollingPeriodFlag, DefaultPollingPeriod, "Polling duration")
	return fctl.NewControllerConfig(
		useMangoPay,
		shortMangoPay,
		shortMangoPay,
		[]string{},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

type MangoPayController struct {
	store  *MangoPayStore
	config *fctl.ControllerConfig
}

func NewMangoPayController(config *fctl.ControllerConfig) *MangoPayController {
	return &MangoPayController{
		store:  NewMangoPayStore(),
		config: config,
	}
}

func (c *MangoPayController) GetStore() *MangoPayStore {
	return c.store
}

func (c *MangoPayController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *MangoPayController) Run() (fctl.Renderable, error) {

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

	if !fctl.CheckStackApprobation(flags, stack, "You are about to install connector '%s'", internal.MangoPayConnector) {
		return nil, fctl.ErrMissingApproval
	}

	paymentsClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	request := operations.InstallConnectorRequest{
		Connector: shared.ConnectorMangopay,
		RequestBody: shared.MangoPayConfig{
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

	c.store.Success = true

	return c, nil
}

func (c *MangoPayController) Render() error {

	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Connector %s installed!", c.store.ConnectorName)

	return nil
}
func NewMangoPayCommand() *cobra.Command {
	c := NewMangoPayConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*MangoPayStore](NewMangoPayController(c)),
	)
}
