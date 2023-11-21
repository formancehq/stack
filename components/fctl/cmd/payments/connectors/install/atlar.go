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

type PaymentsConnectorsAtlarStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
}

type PaymentsConnectorsAtlarController struct {
	store              *PaymentsConnectorsAtlarStore
	atlarAccessKeyFlag string
	atlarSecretFlag    string
}

var _ fctl.Controller[*PaymentsConnectorsAtlarStore] = (*PaymentsConnectorsAtlarController)(nil)

func NewDefaultPaymentsConnectorsAtlarStore() *PaymentsConnectorsAtlarStore {
	return &PaymentsConnectorsAtlarStore{
		Success: false,
	}
}

func NewPaymentsConnectorsAtlarController() *PaymentsConnectorsAtlarController {
	return &PaymentsConnectorsAtlarController{
		store:              NewDefaultPaymentsConnectorsAtlarStore(),
		atlarAccessKeyFlag: "access-key",
		atlarSecretFlag:    "secret",
	}
}

func NewAtlarCommand() *cobra.Command {
	c := NewPaymentsConnectorsAtlarController()
	return fctl.NewCommand(internal.AtlarConnector+" <access-key> <secret>",
		fctl.WithShortDescription("Install an atlar connector"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithStringFlag(c.atlarAccessKeyFlag, "", "Atlar API access key"),
		fctl.WithStringFlag(c.atlarSecretFlag, "", "Atlar API secret"),
		fctl.WithController[*PaymentsConnectorsAtlarStore](c),
	)
}

func (c *PaymentsConnectorsAtlarController) GetStore() *PaymentsConnectorsAtlarStore {
	return c.store
}

func (c *PaymentsConnectorsAtlarController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	soc, err := fctl.GetStackOrganizationConfigApprobation(cmd, "You are about to install connector '%s'", internal.AtlarConnector)
	if err != nil {
		return nil, fctl.ErrMissingApproval
	}

	paymentsClient, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, err
	}

	response, err := paymentsClient.Payments.InstallConnector(cmd.Context(), operations.InstallConnectorRequest{
		ConnectorConfig: shared.ConnectorConfig{
			AtlarConfig: &shared.AtlarConfig{
				AccessKey: args[0],
				Secret:    args[1],
			},
		},
		Connector: shared.ConnectorAtlar,
	})
	if err != nil {
		return nil, errors.Wrap(err, "installing connector")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.ConnectorName = internal.AtlarConnector

	return c, nil
}

func (c *PaymentsConnectorsAtlarController) Render(cmd *cobra.Command, args []string) error {

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector '%s' installed!", c.store.ConnectorName)

	return nil
}
