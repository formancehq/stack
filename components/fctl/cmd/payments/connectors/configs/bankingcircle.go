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

type UpdateBankingCircleConnectorConfigStore struct {
	Success     bool   `json:"success"`
	ConnectorID string `json:"connectorId"`
}

type UpdateBankingCircleConnectorConfigController struct {
	PaymentsVersion versions.Version

	store *UpdateBankingCircleConnectorConfigStore

	connectorIDFlag string
}

func (c *UpdateBankingCircleConnectorConfigController) SetVersion(version versions.Version) {
	c.PaymentsVersion = version
}

var _ fctl.Controller[*UpdateBankingCircleConnectorConfigStore] = (*UpdateBankingCircleConnectorConfigController)(nil)

func NewUpdateBankingCircleConnectorConfigStore() *UpdateBankingCircleConnectorConfigStore {
	return &UpdateBankingCircleConnectorConfigStore{
		Success: false,
	}
}

func NewUpdateBankingCircleConnectorConfigController() *UpdateBankingCircleConnectorConfigController {
	return &UpdateBankingCircleConnectorConfigController{
		store:           NewUpdateBankingCircleConnectorConfigStore(),
		connectorIDFlag: "connector-id",
	}
}

func newUpdateBankingCircleCommand() *cobra.Command {
	c := NewUpdateBankingCircleConnectorConfigController()
	return fctl.NewCommand(internal.BankingCircleConnector+" <file>|-",
		fctl.WithShortDescription("Update the config of a Banking Circle connector"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithStringFlag("connector-id", "", "Connector ID"),
		fctl.WithController[*UpdateBankingCircleConnectorConfigStore](c),
	)
}

func (c *UpdateBankingCircleConnectorConfigController) GetStore() *UpdateBankingCircleConnectorConfigStore {
	return c.store
}

func (c *UpdateBankingCircleConnectorConfigController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	config := &shared.BankingCircleConfig{}
	if err := json.Unmarshal([]byte(script), config); err != nil {
		return nil, err
	}

	response, err := store.Client().Payments.V1.UpdateConnectorConfigV1(cmd.Context(), operations.UpdateConnectorConfigV1Request{
		ConnectorConfig: shared.ConnectorConfig{
			BankingCircleConfig: config,
		},
		Connector:   shared.ConnectorBankingCircle,
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

func (c *UpdateBankingCircleConnectorConfigController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector '%s' updated!", c.store.ConnectorID)

	return nil
}
