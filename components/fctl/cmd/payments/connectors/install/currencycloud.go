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
	defaultEndpointCurrencyCloud = "https://devapi.currencycloud.com"
	useCurrencyCloud             = internal.CurrencyCloudConnector + " <login-id> <api-key>"
	descriptionCurrencyCloud     = "Install CurrencyCloud connector"
	shortCurrencyCloud           = "Install CurrencyCloud connector"
)

type CurrencyCloudStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
}

func NewCurrencyCloudStore() *CurrencyCloudStore {
	return &CurrencyCloudStore{
		Success: false,
	}
}
func NewCurrencyCloudConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useCurrencyCloud, flag.ExitOnError)
	flags.String(EndpointFlag, defaultEndpointCurrencyCloud, "API endpoint")
	flags.String(PollingPeriodFlag, DefaultPollingPeriod, "Polling duration")

	return fctl.NewControllerConfig(
		useCurrencyCloud,
		descriptionCurrencyCloud,
		shortCurrencyCloud,
		[]string{},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

type CurrencyCloudController struct {
	store  *CurrencyCloudStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*CurrencyCloudStore] = (*CurrencyCloudController)(nil)

func NewCurrencyCloudController(config *fctl.ControllerConfig) *CurrencyCloudController {
	return &CurrencyCloudController{
		store:  NewCurrencyCloudStore(),
		config: config,
	}
}

func (c *CurrencyCloudController) GetStore() *CurrencyCloudStore {
	return c.store
}

func (c *CurrencyCloudController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *CurrencyCloudController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()

	soc, err := fctl.GetStackOrganizationConfigApprobation(flags, ctx, fmt.Sprintf("You are about to install connector '%s'", internal.CurrencyCloudConnector), c.config.GetOut())
	if err != nil {
		return nil, fctl.ErrMissingApproval
	}

	paymentsClient, err := fctl.NewStackClient(flags, ctx, soc.Config, soc.Stack, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	var endpoint *string
	if e := fctl.GetString(flags, EndpointFlag); e != "" {
		endpoint = &e
	}

	response, err := paymentsClient.Payments.InstallConnector(ctx, operations.InstallConnectorRequest{
		RequestBody: shared.CurrencyCloudConfig{
			APIKey:        args[1],
			LoginID:       args[0],
			Endpoint:      endpoint,
			PollingPeriod: fctl.Ptr(fctl.GetString(flags, PollingPeriodFlag)),
		},
		Connector: shared.ConnectorCurrencyCloud,
	})
	if err != nil {
		return nil, errors.Wrap(err, "installing connector")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.ConnectorName = internal.CurrencyCloudConnector

	return c, nil
}

func (c *CurrencyCloudController) Render() error {

	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Connector '%s' installed!", c.store.ConnectorName)

	return nil
}

func NewCurrencyCloudCommand() *cobra.Command {
	c := NewCurrencyCloudConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*CurrencyCloudStore](NewCurrencyCloudController(c)),
	)
}
