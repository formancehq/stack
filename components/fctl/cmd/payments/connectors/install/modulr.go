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
	defaultEndpointModulr = "https://api-sandbox.modulrfinance.com"
	useModulrConnector    = internal.ModulrConnector + " <api-key> <api-secret>"
	shortModulrConnector  = "Install Modulr connector"
)

type ModulrStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
}

func NewDefaultModulrStore() *ModulrStore {
	return &ModulrStore{
		Success: false,
	}
}
func NewModulrConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useModulrConnector, flag.ExitOnError)
	flags.String(EndpointFlag, defaultEndpointModulr, "API endpoint")
	flags.String(PollingPeriodFlag, DefaultPollingPeriod, "Polling duration")
	return fctl.NewControllerConfig(
		useModulrConnector,
		shortModulrConnector,
		shortModulrConnector,
		[]string{},
		flags,
		fctl.Organization, fctl.Stack,
	)

}

type ModulrController struct {
	store  *ModulrStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*ModulrStore] = (*ModulrController)(nil)

func NewModulrController(config *fctl.ControllerConfig) *ModulrController {
	return &ModulrController{
		store:  NewDefaultModulrStore(),
		config: config,
	}
}

func (c *ModulrController) GetStore() *ModulrStore {
	return c.store
}

func (c *ModulrController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *ModulrController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()

	soc, err := fctl.GetStackOrganizationConfigApprobation(flags, ctx, fmt.Sprintf("You are about to install connector '%s'", internal.ModulrConnector), c.config.GetOut())
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
		RequestBody: shared.ModulrConfig{
			APIKey:        args[0],
			APISecret:     args[1],
			Endpoint:      endpoint,
			PollingPeriod: fctl.Ptr(fctl.GetString(flags, PollingPeriodFlag)),
		},
		Connector: shared.ConnectorModulr,
	})
	if err != nil {
		return nil, errors.Wrap(err, "installing connector")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.ConnectorName = internal.ModulrConnector

	return c, nil
}

func (c *ModulrController) Render() error {

	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Connector '%s' installed!", c.store.ConnectorName)

	return nil
}

func NewModulrCommand() *cobra.Command {
	config := NewModulrConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*ModulrStore](NewModulrController(config)),
	)
}
