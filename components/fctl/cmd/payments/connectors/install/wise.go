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

type PaymentsConnectorsWiseStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
	ConnectorID   string `json:"connectorId"`
}
type PaymentsConnectorsWiseController struct {
	store *PaymentsConnectorsWiseStore
}

var _ fctl.Controller[*PaymentsConnectorsWiseStore] = (*PaymentsConnectorsWiseController)(nil)

func NewDefaultPaymentsConnectorsWiseStore() *PaymentsConnectorsWiseStore {
	return &PaymentsConnectorsWiseStore{
		Success: false,
	}
}

func NewPaymentsConnectorsWiseController() *PaymentsConnectorsWiseController {
	return &PaymentsConnectorsWiseController{
		store: NewDefaultPaymentsConnectorsWiseStore(),
	}
}

func NewWiseCommand() *cobra.Command {
	c := NewPaymentsConnectorsWiseController()
	return fctl.NewCommand(internal.WiseConnector+" <file>|-",
		fctl.WithShortDescription("Install a Wise connector"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*PaymentsConnectorsWiseStore](c),
	)
}

func (c *PaymentsConnectorsWiseController) GetStore() *PaymentsConnectorsWiseStore {
	return c.store
}

func (c *PaymentsConnectorsWiseController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())
	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to install connector '%s'", internal.WiseConnector) {
		return nil, fctl.ErrMissingApproval
	}
	script, err := fctl.ReadFile(cmd, store.Stack(), args[0])
	if err != nil {
		return nil, err
	}

	var config shared.WiseConfig
	if err := json.Unmarshal([]byte(script), &config); err != nil {
		return nil, err
	}

	response, err := store.Client().Payments.InstallConnector(cmd.Context(), operations.InstallConnectorRequest{
		ConnectorConfig: shared.ConnectorConfig{
			WiseConfig: &config,
		},
		Connector: shared.ConnectorWise,
	})
	if err != nil {
		return nil, errors.Wrap(err, "installing connector")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.ConnectorName = internal.WiseConnector

	if response.ConnectorResponse != nil {
		c.store.ConnectorID = response.ConnectorResponse.Data.ConnectorID
	}

	return c, nil
}

func (c *PaymentsConnectorsWiseController) Render(cmd *cobra.Command, args []string) error {
	if c.store.ConnectorID == "" {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector installed!", c.store.ConnectorName)
	} else {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector '%s' installed!", c.store.ConnectorName, c.store.ConnectorID)
	}
	return nil
}
