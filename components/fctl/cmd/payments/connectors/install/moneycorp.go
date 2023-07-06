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

type PaymentsConnectorsMoneycorpStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
}
type PaymentsConnectorsMoneycorpController struct {
	store           *PaymentsConnectorsMoneycorpStore
	endpointFlag    string
	defaultEndpoint string
}

func NewDefaultPaymentsConnectorsMoneycorpStore() *PaymentsConnectorsMoneycorpStore {
	return &PaymentsConnectorsMoneycorpStore{
		Success:       false,
		ConnectorName: internal.MoneycorpConnector,
	}
}
func NewPaymentsConnectorsMoneycorpController() *PaymentsConnectorsMoneycorpController {
	return &PaymentsConnectorsMoneycorpController{
		store:           NewDefaultPaymentsConnectorsMoneycorpStore(),
		endpointFlag:    "endpoint",
		defaultEndpoint: "https://sandbox-corpapi.moneycorp.com",
	}
}
func NewMoneycorpCommand() *cobra.Command {
	c := NewPaymentsConnectorsMoneycorpController()

	return fctl.NewCommand(internal.MoneycorpConnector+" <clientID> <apiKey>",
		fctl.WithShortDescription("Install a Moneycorp connector"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithStringFlag(c.endpointFlag, c.defaultEndpoint, "API endpoint"),
		fctl.WithController[*PaymentsConnectorsMoneycorpStore](c),
	)
}

func (c *PaymentsConnectorsMoneycorpController) GetStore() *PaymentsConnectorsMoneycorpStore {
	return c.store
}

func (c *PaymentsConnectorsMoneycorpController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to install connector '%s'", internal.MoneycorpConnector) {
		return nil, fctl.ErrMissingApproval
	}

	paymentsClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	request := operations.InstallConnectorRequest{
		Connector: shared.ConnectorMoneycorp,
		RequestBody: shared.MoneycorpConfig{
			ClientID: args[0],
			APIKey:   args[1],
			Endpoint: fctl.GetString(cmd, c.endpointFlag),
		},
	}
	response, err := paymentsClient.Payments.InstallConnector(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "installing connector")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector installed!")

	c.store.Success = true

	return c, nil
}

func (c *PaymentsConnectorsMoneycorpController) Render(cmd *cobra.Command, args []string) error {

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector %s installed!", c.store.ConnectorName)

	return nil
}
