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
	authorizationEndpointFlag    = "authorization-endpoint"
	defaultEndpointBankingCircle = "https://sandbox.bankingcircle.com"
	defaultAuthorizationEndpoint = "https://authorizationsandbox.bankingcircleconnect.com"
	useBankingCircle             = internal.BankingCircleConnector + " <username> <password>"
	descriptionBankingCircle     = "Install Banking Circle connector"
)

type BankingCircleStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
}

func NewBankingCircleStore() *BankingCircleStore {
	return &BankingCircleStore{
		Success: false,
	}
}

func NewBankingCircleConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useBankingCircle, flag.ExitOnError)
	flags.String(EndpointFlag, defaultEndpointBankingCircle, "API endpoint")
	flags.String(authorizationEndpointFlag, defaultAuthorizationEndpoint, "Authorization endpoint")
	flags.String(PollingPeriodFlag, DefaultPollingPeriod, "Polling duration")
	return fctl.NewControllerConfig(
		useBankingCircle,
		descriptionBankingCircle,
		descriptionBankingCircle,
		[]string{},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

type BankingCircleController struct {
	store  *BankingCircleStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*BankingCircleStore] = (*BankingCircleController)(nil)

func NewBankingCircleController(config *fctl.ControllerConfig) *BankingCircleController {
	return &BankingCircleController{
		store:  NewBankingCircleStore(),
		config: config,
	}
}

func (c *BankingCircleController) GetStore() *BankingCircleStore {
	return c.store
}

func (c *BankingCircleController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *BankingCircleController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()

	soc, err := fctl.GetStackOrganizationConfigApprobation(flags, ctx, fmt.Sprintf("You are about to install connector '%s'", internal.BankingCircleConnector), c.config.GetOut())
	if err != nil {
		return nil, fctl.ErrMissingApproval
	}

	paymentsClient, err := fctl.NewStackClient(flags, ctx, soc.Config, soc.Stack, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	request := operations.InstallConnectorRequest{
		Connector: shared.ConnectorBankingCircle,
		RequestBody: shared.BankingCircleConfig{
			Username:              args[0],
			Password:              args[1],
			Endpoint:              fctl.GetString(flags, EndpointFlag),
			AuthorizationEndpoint: fctl.GetString(flags, authorizationEndpointFlag),
			PollingPeriod:         fctl.Ptr(fctl.GetString(flags, PollingPeriodFlag)),
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
	c.store.ConnectorName = internal.BankingCircleConnector

	return c, nil
}

func (c *BankingCircleController) Render() error {

	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Connector '%s' installed!", c.store.ConnectorName)

	return nil
}
func NewBankingCircleCommand() *cobra.Command {
	c := NewBankingCircleConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*BankingCircleStore](NewBankingCircleController(c)),
	)
}
