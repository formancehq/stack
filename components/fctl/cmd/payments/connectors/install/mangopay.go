package install

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/connectors/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type PaymentsConnectorsMangoPayStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
	ConnectorID   string `json:"connectorId"`
}
type PaymentsConnectorsMangoPayController struct {
	store *PaymentsConnectorsMangoPayStore
}

func NewDefaultPaymentsConnectorsMangoPayStore() *PaymentsConnectorsMangoPayStore {
	return &PaymentsConnectorsMangoPayStore{
		Success:       false,
		ConnectorName: internal.MangoPayConnector,
	}
}
func NewPaymentsConnectorsMangoPayController() *PaymentsConnectorsMangoPayController {
	return &PaymentsConnectorsMangoPayController{
		store: NewDefaultPaymentsConnectorsMangoPayStore(),
	}
}

func NewMangoPayCommand() *cobra.Command {
	c := NewPaymentsConnectorsMangoPayController()
	return fctl.NewCommand(internal.MangoPayConnector+" <file>|-",
		fctl.WithShortDescription("Install a MangoPay connector"),
		fctl.WithArgs(cobra.ExactArgs(1)),
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

	script, err := fctl.ReadFile(cmd, stack, args[0])
	if err != nil {
		return nil, err
	}

	var config shared.MangoPayConfig
	if err := json.Unmarshal([]byte(script), &config); err != nil {
		return nil, err
	}

	request := operations.InstallConnectorRequest{
		Connector: shared.ConnectorMangopay,
		ConnectorConfig: shared.ConnectorConfig{
			MangoPayConfig: &config,
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
	c.store.ConnectorName = internal.MangoPayConnector

	if response.ConnectorResponse != nil {
		c.store.ConnectorID = response.ConnectorResponse.Data.ConnectorID
	}

	return c, nil
}

func (c *PaymentsConnectorsMangoPayController) Render(cmd *cobra.Command, args []string) error {
	if c.store.ConnectorID == "" {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector installed!", c.store.ConnectorName)
	} else {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector '%s' installed!", c.store.ConnectorName, c.store.ConnectorID)
	}

	return nil
}
