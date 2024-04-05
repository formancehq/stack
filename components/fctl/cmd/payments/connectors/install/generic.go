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

type PaymentsConnectorsGenericStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
	ConnectorID   string `json:"connectorId"`
}
type PaymentsConnectorsGenericController struct {
	store *PaymentsConnectorsGenericStore
}

func NewDefaultPaymentsConnectorsGenericStore() *PaymentsConnectorsGenericStore {
	return &PaymentsConnectorsGenericStore{
		Success:       false,
		ConnectorName: internal.GenericConnector,
	}
}
func NewPaymentsConnectorsGenericController() *PaymentsConnectorsGenericController {
	return &PaymentsConnectorsGenericController{
		store: NewDefaultPaymentsConnectorsGenericStore(),
	}
}

func NewGenericCommand() *cobra.Command {
	c := NewPaymentsConnectorsGenericController()
	return fctl.NewCommand(internal.GenericConnector+" <file>|-",
		fctl.WithShortDescription("Install a Generic connector"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*PaymentsConnectorsGenericStore](c),
	)
}

func (c *PaymentsConnectorsGenericController) GetStore() *PaymentsConnectorsGenericStore {
	return c.store
}

func (c *PaymentsConnectorsGenericController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to install connector '%s'", internal.GenericConnector) {
		return nil, fctl.ErrMissingApproval
	}

	script, err := fctl.ReadFile(cmd, store.Stack(), args[0])
	if err != nil {
		return nil, err
	}

	var config shared.GenericConfig
	if err := json.Unmarshal([]byte(script), &config); err != nil {
		return nil, err
	}

	request := operations.InstallConnectorRequest{
		Connector: shared.ConnectorGeneric,
		ConnectorConfig: shared.ConnectorConfig{
			GenericConfig: &config,
		},
	}
	response, err := store.Client().Payments.InstallConnector(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "installing connector")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.ConnectorName = internal.GenericConnector

	if response.ConnectorResponse != nil {
		c.store.ConnectorID = response.ConnectorResponse.Data.ConnectorID
	}

	return c, nil
}

func (c *PaymentsConnectorsGenericController) Render(cmd *cobra.Command, args []string) error {
	if c.store.ConnectorID == "" {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector installed!", c.store.ConnectorName)
	} else {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector '%s' installed!", c.store.ConnectorName, c.store.ConnectorID)
	}

	return nil
}
