package configs

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/connectors/internal"
	"github.com/formancehq/fctl/cmd/payments/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type UpdateCurrencyCloudConnectorConfigStore struct {
	Success     bool   `json:"success"`
	ConnectorID string `json:"connectorId"`
}

type UpdateCurrencyCloudConnectorConfigController struct {
	PaymentsVersion versions.Version

	store *UpdateCurrencyCloudConnectorConfigStore

	connectorIDFlag string
}

func (c *UpdateCurrencyCloudConnectorConfigController) SetVersion(version versions.Version) {
	c.PaymentsVersion = version
}

var _ fctl.Controller[*UpdateCurrencyCloudConnectorConfigStore] = (*UpdateCurrencyCloudConnectorConfigController)(nil)

func NewUpdateCurrencyCloudConnectorConfigStore() *UpdateCurrencyCloudConnectorConfigStore {
	return &UpdateCurrencyCloudConnectorConfigStore{
		Success: false,
	}
}

func NewUpdateCurrencyCloudConnectorConfigController() *UpdateCurrencyCloudConnectorConfigController {
	return &UpdateCurrencyCloudConnectorConfigController{
		store:           NewUpdateCurrencyCloudConnectorConfigStore(),
		connectorIDFlag: "connector-id",
	}
}

func newUpdateCurrencyCloudCommand() *cobra.Command {
	c := NewUpdateCurrencyCloudConnectorConfigController()
	return fctl.NewCommand(internal.CurrencyCloudConnector+" <file>|-",
		fctl.WithShortDescription("Update the config of a CurrencyCloud connector"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithStringFlag("connector-id", "", "Connector ID"),
		fctl.WithController[*UpdateCurrencyCloudConnectorConfigStore](c),
	)
}

func (c *UpdateCurrencyCloudConnectorConfigController) GetStore() *UpdateCurrencyCloudConnectorConfigStore {
	return c.store
}

func (c *UpdateCurrencyCloudConnectorConfigController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	soc, err := fctl.GetStackOrganizationConfigApprobation(cmd, "You are about to update the config of connector '%s'", connectorID)
	if err != nil {
		return nil, fctl.ErrMissingApproval
	}

	paymentsClient, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, err
	}

	script, err := fctl.ReadFile(cmd, soc.Stack, args[0])
	if err != nil {
		return nil, err
	}

	config := &shared.CurrencyCloudConfig{}
	if err := json.Unmarshal([]byte(script), config); err != nil {
		return nil, err
	}

	response, err := paymentsClient.Payments.UpdateConnectorConfigV1(cmd.Context(), operations.UpdateConnectorConfigV1Request{
		ConnectorConfig: shared.ConnectorConfig{
			CurrencyCloudConfig: config,
		},
		Connector:   shared.ConnectorCurrencyCloud,
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

func (c *UpdateCurrencyCloudConnectorConfigController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector '%s' updated!", c.store.ConnectorID)

	return nil
}
