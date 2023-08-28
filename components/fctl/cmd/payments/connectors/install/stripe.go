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
	stripeApiKeyFLag     = "api-key"
	useStripeConnector   = internal.StripeConnector + " <api-key>"
	shortStripeConnector = "Install Stripe connector"
)

type StripeStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
}
type StripeController struct {
	store  *StripeStore
	config *fctl.ControllerConfig
}

func NewStripeConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useStripeConnector, flag.ExitOnError)
	flags.String(stripeApiKeyFLag, "", "Stripe API key")
	fctl.WithConfirmFlag(flags)
	return fctl.NewControllerConfig(
		useStripeConnector,
		shortStripeConnector,
		shortStripeConnector,
		[]string{},
		flags,
		fctl.Organization, fctl.Stack,
	)

}

var _ fctl.Controller[*StripeStore] = (*StripeController)(nil)

func NewDefaultStripeStore() *StripeStore {
	return &StripeStore{
		Success: false,
	}
}

func NewStripeController(config *fctl.ControllerConfig) *StripeController {
	return &StripeController{
		store:  NewDefaultStripeStore(),
		config: config,
	}
}

func (c *StripeController) GetStore() *StripeStore {
	return c.store
}

func (c *StripeController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *StripeController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()

	soc, err := fctl.GetStackOrganizationConfigApprobation(flags, ctx, fmt.Sprintf("You are about to install connector '%s'", internal.StripeConnector), c.config.GetOut())
	if err != nil {
		return nil, fctl.ErrMissingApproval
	}

	paymentsClient, err := fctl.NewStackClient(flags, ctx, soc.Config, soc.Stack, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	response, err := paymentsClient.Payments.InstallConnector(ctx, operations.InstallConnectorRequest{
		RequestBody: shared.StripeConfig{
			APIKey: args[0],
		},
		Connector: shared.ConnectorStripe,
	})
	if err != nil {
		return nil, errors.Wrap(err, "installing connector")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.ConnectorName = internal.StripeConnector

	return c, nil
}

func (c *StripeController) Render() error {

	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Connector '%s' installed!", c.store.ConnectorName)

	return nil
}
func NewStripeCommand() *cobra.Command {
	config := NewStripeConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*StripeStore](NewStripeController(config)),
	)
}
