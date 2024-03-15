package configs

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/connectors/internal"
	"github.com/formancehq/fctl/cmd/payments/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type UpdateAtlarConnectorConfigStore struct {
	Success     bool   `json:"success"`
	ConnectorID string `json:"connectorID"`
}

type UpdateAtlarConnectorConfigController struct {
	PaymentsVersion versions.Version

	store *UpdateAtlarConnectorConfigStore

	connectorIDFlag string
}

func (c *UpdateAtlarConnectorConfigController) SetVersion(version versions.Version) {
	c.PaymentsVersion = version
}

var _ fctl.Controller[*UpdateAtlarConnectorConfigStore] = (*UpdateAtlarConnectorConfigController)(nil)

func NewUpdateAtlarConnectorConfigStore() *UpdateAtlarConnectorConfigStore {
	return &UpdateAtlarConnectorConfigStore{
		Success: false,
	}
}

func NewUpdateAtlarConnectorConfigController() *UpdateAtlarConnectorConfigController {
	return &UpdateAtlarConnectorConfigController{
		store:           NewUpdateAtlarConnectorConfigStore(),
		connectorIDFlag: "connector-id",
	}
}

func newUpdateAtlarCommand() *cobra.Command {
	c := NewUpdateAtlarConnectorConfigController()
	return fctl.NewCommand(internal.AtlarConnector+" <file>|-",
		fctl.WithShortDescription("Update the config of an Atlar connector"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithStringFlag("connector-id", "", "Connector ID"),
		fctl.WithController[*UpdateAtlarConnectorConfigStore](c),
	)
}

func (c *UpdateAtlarConnectorConfigController) GetStore() *UpdateAtlarConnectorConfigStore {
	return c.store
}

func (c *UpdateAtlarConnectorConfigController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	if err := versions.GetPaymentsVersion(cmd, args, c); err != nil {
		return nil, err
	}

	if c.PaymentsVersion < versions.V1 {
		return nil, fmt.Errorf("update configs are only supported in >= v1.0.0")
	}

	connectorID := fctl.GetString(cmd, c.connectorIDFlag)
	if connectorID == "" {
		return nil, fmt.Errorf("missing connector ID")
	}

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to update the config of connector '%s'", connectorID) {
		return nil, fctl.ErrMissingApproval
	}

	script, err := fctl.ReadFile(cmd, store.Stack(), args[0])
	if err != nil {
		return nil, err
	}

	config := &shared.AtlarConfig{}
	if err := json.Unmarshal([]byte(script), config); err != nil {
		return nil, err
	}

	response, err := store.Client().Payments.UpdateConnectorConfigV1(cmd.Context(), operations.UpdateConnectorConfigV1Request{
		ConnectorConfig: shared.ConnectorConfig{
			AtlarConfig: config,
		},
		Connector:   shared.ConnectorAtlar,
		ConnectorID: connectorID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "updating config of connector")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.ConnectorID = connectorID

	return c, nil
}

func (c *UpdateAtlarConnectorConfigController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector '%s' updated!", c.store.ConnectorID)

	return nil
}
