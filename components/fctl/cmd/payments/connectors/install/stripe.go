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

type PaymentsConnectorsStripeStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
}
type PaymentsConnectorsStripeController struct {
	store            *PaymentsConnectorsStripeStore
	stripeApiKeyFlag string
}

var _ fctl.Controller[*PaymentsConnectorsStripeStore] = (*PaymentsConnectorsStripeController)(nil)

func NewDefaultPaymentsConnectorsStripeStore() *PaymentsConnectorsStripeStore {
	return &PaymentsConnectorsStripeStore{
		Success: false,
	}
}

func NewPaymentsConnectorsStripeController() *PaymentsConnectorsStripeController {
	return &PaymentsConnectorsStripeController{
		store:            NewDefaultPaymentsConnectorsStripeStore(),
		stripeApiKeyFlag: "api-key",
	}
}

func NewStripeCommand() *cobra.Command {
	c := NewPaymentsConnectorsStripeController()
	return fctl.NewCommand(internal.StripeConnector+" <api-key>",
		fctl.WithShortDescription("Install a stripe connector"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithStringFlag(c.stripeApiKeyFlag, "", "Stripe API key"),
		fctl.WithController[*PaymentsConnectorsStripeStore](c),
	)
}

func (c *PaymentsConnectorsStripeController) GetStore() *PaymentsConnectorsStripeStore {
	return c.store
}

func (c *PaymentsConnectorsStripeController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	soc, err := fctl.GetStackOrganizationConfigApprobation(cmd, "You are about to install connector '%s'", internal.StripeConnector)
	if err != nil {
		return nil, fctl.ErrMissingApproval
	}

	paymentsClient, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, err
	}

	response, err := paymentsClient.Payments.InstallConnector(cmd.Context(), operations.InstallConnectorRequest{
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

func (c *PaymentsConnectorsStripeController) Render(cmd *cobra.Command, args []string) error {

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector '%s' installed!", c.store.ConnectorName)

	return nil
}
