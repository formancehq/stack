package install

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/connectors/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type PaymentsConnectorsBankingCircleStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
}
type PaymentsConnectorsBankingCircleController struct {
	store                        *PaymentsConnectorsBankingCircleStore
	endpointFlag                 string
	authorizationEndpointFlag    string
	defaultEndpoint              string
	defaultAuthorizationEndpoint string
}

var _ fctl.Controller[*PaymentsConnectorsBankingCircleStore] = (*PaymentsConnectorsBankingCircleController)(nil)

func NewDefaultPaymentsConnectorsBankingCircleStore() *PaymentsConnectorsBankingCircleStore {
	return &PaymentsConnectorsBankingCircleStore{
		Success: false,
	}
}

func NewPaymentsConnectorsBankingCircleController() *PaymentsConnectorsBankingCircleController {
	return &PaymentsConnectorsBankingCircleController{
		store:                        NewDefaultPaymentsConnectorsBankingCircleStore(),
		endpointFlag:                 "endpoint",
		authorizationEndpointFlag:    "authorization-endpoint",
		defaultEndpoint:              "https://sandbox.bankingcircle.com",
		defaultAuthorizationEndpoint: "https://authorizationsandbox.bankingcircleconnect.com",
	}
}

func NewBankingCircleCommand() *cobra.Command {
	c := NewPaymentsConnectorsBankingCircleController()
	return fctl.NewCommand(internal.BankingCircleConnector+" <username> <password>",
		fctl.WithShortDescription("Install a Banking Circle connector"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithStringFlag(c.endpointFlag, c.defaultEndpoint, "API endpoint"),
		fctl.WithStringFlag(c.authorizationEndpointFlag, c.defaultAuthorizationEndpoint, "Authorization endpoint"),
		fctl.WithController[*PaymentsConnectorsBankingCircleStore](c),
	)
}

func (c *PaymentsConnectorsBankingCircleController) GetStore() *PaymentsConnectorsBankingCircleStore {
	return c.store
}

func (c *PaymentsConnectorsBankingCircleController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	soc, err := fctl.GetStackOrganizationConfigApprobation(cmd, "You are about to install connector '%s'", internal.BankingCircleConnector)
	if err != nil {
		return nil, fctl.ErrMissingApproval
	}

	paymentsClient, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, err
	}

	request := operations.InstallConnectorRequest{
		Connector: shared.ConnectorBankingCircle,
		RequestBody: shared.BankingCircleConfig{
			Username:              args[0],
			Password:              args[1],
			Endpoint:              fctl.GetString(cmd, c.endpointFlag),
			AuthorizationEndpoint: fctl.GetString(cmd, c.authorizationEndpointFlag),
		},
	}
	response, err := paymentsClient.Payments.InstallConnector(cmd.Context(), request)
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

func (c *PaymentsConnectorsBankingCircleController) Render(cmd *cobra.Command, args []string) error {

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector '%s' installed!", c.store.ConnectorName)

	return nil
}
