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

type PaymentsConnectorsBankingCircleStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
	ConnectorID   string `json:"connectorId"`
}
type PaymentsConnectorsBankingCircleController struct {
	store *PaymentsConnectorsBankingCircleStore
}

var _ fctl.Controller[*PaymentsConnectorsBankingCircleStore] = (*PaymentsConnectorsBankingCircleController)(nil)

func NewDefaultPaymentsConnectorsBankingCircleStore() *PaymentsConnectorsBankingCircleStore {
	return &PaymentsConnectorsBankingCircleStore{
		Success: false,
	}
}

func NewPaymentsConnectorsBankingCircleController() *PaymentsConnectorsBankingCircleController {
	return &PaymentsConnectorsBankingCircleController{
		store: NewDefaultPaymentsConnectorsBankingCircleStore(),
	}
}

func NewBankingCircleCommand() *cobra.Command {
	c := NewPaymentsConnectorsBankingCircleController()
	return fctl.NewCommand(internal.BankingCircleConnector+" <file>|-",
		fctl.WithShortDescription("Install a Banking Circle connector"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*PaymentsConnectorsBankingCircleStore](c),
	)
}

func (c *PaymentsConnectorsBankingCircleController) GetStore() *PaymentsConnectorsBankingCircleStore {
	return c.store
}

func (c *PaymentsConnectorsBankingCircleController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	soc, err := fctl.GetStackOrganizationConfigApprobation(cmd, "You are about to install connector '%s'", internal.BankingCircleConnector)
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

	var config shared.BankingCircleConfig
	if err := json.Unmarshal([]byte(script), &config); err != nil {
		return nil, err
	}

	request := operations.InstallConnectorRequest{
		Connector: shared.ConnectorBankingCircle,
		ConnectorConfig: shared.ConnectorConfig{
			BankingCircleConfig: &config,
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
	c.store.ConnectorName = internal.BankingCircleConnector

	if response.ConnectorResponse != nil {
		c.store.ConnectorID = response.ConnectorResponse.Data.ConnectorID
	}

	return c, nil
}

func (c *PaymentsConnectorsBankingCircleController) Render(cmd *cobra.Command, args []string) error {
	if c.store.ConnectorID == "" {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector installed!", c.store.ConnectorName)
	} else {
		pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("%s: connector '%s' installed!", c.store.ConnectorName, c.store.ConnectorID)
	}

	return nil
}
