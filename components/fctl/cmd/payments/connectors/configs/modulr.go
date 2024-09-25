package configs

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/connectors/internal"
	"github.com/formancehq/fctl/cmd/payments/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type UpdateModulrConnectorConfigStore struct {
	Success     bool   `json:"success"`
	ConnectorID string `json:"connectorId"`
}

type UpdateModulrConnectorConfigController struct {
	PaymentsVersion versions.Version

	store *UpdateModulrConnectorConfigStore

	connectorIDFlag string
}

func (c *UpdateModulrConnectorConfigController) SetVersion(version versions.Version) {
	c.PaymentsVersion = version
}

var _ fctl.Controller[*UpdateModulrConnectorConfigStore] = (*UpdateModulrConnectorConfigController)(nil)

func NewUpdateModulrConnectorConfigStore() *UpdateModulrConnectorConfigStore {
	return &UpdateModulrConnectorConfigStore{
		Success: false,
	}
}

func NewUpdateModulrConnectorConfigController() *UpdateModulrConnectorConfigController {
	return &UpdateModulrConnectorConfigController{
		store:           NewUpdateModulrConnectorConfigStore(),
		connectorIDFlag: "connector-id",
	}
}

func newUpdateModulrCommand() *cobra.Command {
	c := NewUpdateModulrConnectorConfigController()
	return fctl.NewCommand(internal.ModulrConnector+" <file>|-",
		fctl.WithShortDescription("Update the config of a Modulr connector"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithStringFlag("connector-id", "", "Connector ID"),
		fctl.WithController[*UpdateModulrConnectorConfigStore](c),
	)
}

func (c *UpdateModulrConnectorConfigController) GetStore() *UpdateModulrConnectorConfigStore {
	return c.store
}

func (c *UpdateModulrConnectorConfigController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	config := &shared.ModulrConfig{}
	if err := json.Unmarshal([]byte(script), config); err != nil {
		return nil, err
	}

	response, err := store.Client().Payments.V1.UpdateConnectorConfigV1(cmd.Context(), operations.UpdateConnectorConfigV1Request{
		ConnectorConfig: shared.ConnectorConfig{
			ModulrConfig: config,
		},
		Connector:   shared.ConnectorModulr,
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

func (c *UpdateModulrConnectorConfigController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector '%s' updated!", c.store.ConnectorID)

	return nil
}
