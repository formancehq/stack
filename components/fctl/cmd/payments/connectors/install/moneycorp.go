package install

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/connectors/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type PaymentsConnectorsMoneycorpStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
	ConnectorID   string `json:"connectorId"`
}
type PaymentsConnectorsMoneycorpController struct {
	store *PaymentsConnectorsMoneycorpStore
}

func NewDefaultPaymentsConnectorsMoneycorpStore() *PaymentsConnectorsMoneycorpStore {
	return &PaymentsConnectorsMoneycorpStore{
		Success:       false,
		ConnectorName: internal.MoneycorpConnector,
	}
}
func NewPaymentsConnectorsMoneycorpController() *PaymentsConnectorsMoneycorpController {
	return &PaymentsConnectorsMoneycorpController{
		store: NewDefaultPaymentsConnectorsMoneycorpStore(),
	}
}
func NewMoneycorpCommand() *cobra.Command {
	c := NewPaymentsConnectorsMoneycorpController()

	return fctl.NewCommand(internal.MoneycorpConnector+" <file>|-",
		fctl.WithShortDescription("Install a Moneycorp connector"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithConfirmFlag(),
		fctl.WithController[*PaymentsConnectorsMoneycorpStore](c),
	)
}

func (c *PaymentsConnectorsMoneycorpController) GetStore() *PaymentsConnectorsMoneycorpStore {
	return c.store
}

func (c *PaymentsConnectorsMoneycorpController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to install connector '%s'", internal.MoneycorpConnector) {
		return nil, fctl.ErrMissingApproval
	}

	script, err := fctl.ReadFile(cmd, store.Stack(), args[0])
	if err != nil {
		return nil, err
	}

	var config shared.MoneycorpConfig
	if err := json.Unmarshal([]byte(script), &config); err != nil {
		return nil, err
	}

	request := operations.InstallConnectorRequest{
		Connector: shared.ConnectorMoneycorp,
		ConnectorConfig: shared.ConnectorConfig{
			MoneycorpConfig: &config,
		},
	}
	response, err := store.Client().Payments.V1.InstallConnector(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "installing connector")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector installed!")

	c.store.Success = true
	c.store.ConnectorName = internal.MoneycorpConnector

	if response.ConnectorResponse != nil {
		c.store.ConnectorID = response.ConnectorResponse.Data.ConnectorID
	}

	return c, nil
}

func (c *PaymentsConnectorsMoneycorpController) Render(cmd *cobra.Command, args []string) error {
	if c.store.ConnectorID == "" {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector installed!", c.store.ConnectorName)
	} else {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector '%s' installed!", c.store.ConnectorName, c.store.ConnectorID)
	}

	return nil
}
