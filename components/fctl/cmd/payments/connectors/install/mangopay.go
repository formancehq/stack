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

type PaymentsConnectorsMangoPayStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
}
type PaymentsConnectorsMangoPayController struct {
	store           *PaymentsConnectorsMangoPayStore
	endpointFlag    string
	defaultEndpoint string
}

func NewDefaultPaymentsConnectorsMangoPayStore() *PaymentsConnectorsMangoPayStore {
	return &PaymentsConnectorsMangoPayStore{
		Success:       false,
		ConnectorName: internal.MangoPayConnector,
	}
}
func NewPaymentsConnectorsMangoPayController() *PaymentsConnectorsMangoPayController {
	return &PaymentsConnectorsMangoPayController{
		store:           NewDefaultPaymentsConnectorsMangoPayStore(),
		endpointFlag:    "endpoint",
		defaultEndpoint: "https://api.sandbox.mangopay.com",
	}
}

func NewMangoPayCommand() *cobra.Command {
	c := NewPaymentsConnectorsMangoPayController()
	return fctl.NewCommand(internal.MangoPayConnector+" <clientID> <apiKey>",
		fctl.WithShortDescription("Install a MangoPay connector"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithStringFlag(c.endpointFlag, c.defaultEndpoint, "API endpoint"),
		fctl.WithController[*PaymentsConnectorsMangoPayStore](c),
	)
}

func (c *PaymentsConnectorsMangoPayController) GetStore() *PaymentsConnectorsMangoPayStore {
	return c.store
}

func (c *PaymentsConnectorsMangoPayController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to install connector '%s'", internal.MangoPayConnector) {
		return nil, fctl.ErrMissingApproval
	}

	paymentsClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	request := operations.InstallConnectorRequest{
		Connector: shared.ConnectorMangopay,
		RequestBody: shared.MangoPayConfig{
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

	c.store.Success = true

	return c, nil
}

func (c *PaymentsConnectorsMangoPayController) Render(cmd *cobra.Command, args []string) error {

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector %s installed!", c.store.ConnectorName)

	return nil
}
