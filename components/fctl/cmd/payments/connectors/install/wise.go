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

type PaymentsConnectorsWiseStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
}
type PaymentsConnectorsWiseController struct {
	store                *PaymentsConnectorsWiseStore
	pollingPeriodFlag    string
	defaultpollingPeriod string
}

var _ fctl.Controller[*PaymentsConnectorsWiseStore] = (*PaymentsConnectorsWiseController)(nil)

func NewDefaultPaymentsConnectorsWiseStore() *PaymentsConnectorsWiseStore {
	return &PaymentsConnectorsWiseStore{
		Success: false,
	}
}

func NewPaymentsConnectorsWiseController() *PaymentsConnectorsWiseController {
	return &PaymentsConnectorsWiseController{
		store:                NewDefaultPaymentsConnectorsWiseStore(),
		pollingPeriodFlag:    "polling-period",
		defaultpollingPeriod: "2m",
	}
}

func NewWiseCommand() *cobra.Command {
	c := NewPaymentsConnectorsWiseController()
	return fctl.NewCommand(internal.WiseConnector+" <api-key>",
		fctl.WithShortDescription("Install a Wise connector"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithStringFlag(c.pollingPeriodFlag, c.defaultpollingPeriod, "Polling duration"),
		fctl.WithController[*PaymentsConnectorsWiseStore](c),
	)
}

func (c *PaymentsConnectorsWiseController) GetStore() *PaymentsConnectorsWiseStore {
	return c.store
}

func (c *PaymentsConnectorsWiseController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	soc, err := fctl.GetStackOrganizationConfigApprobation(cmd, "You are about to install connector '%s'", internal.WiseConnector)
	if err != nil {
		return nil, fctl.ErrMissingApproval
	}

	paymentsClient, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, err
	}

	response, err := paymentsClient.Payments.InstallConnector(cmd.Context(), operations.InstallConnectorRequest{
		RequestBody: shared.WiseConfig{
			APIKey:        args[1],
			PollingPeriod: fctl.Ptr(fctl.GetString(cmd, c.pollingPeriodFlag)),
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

	return c, nil
}

func (c *PaymentsConnectorsWiseController) Render(cmd *cobra.Command, args []string) error {

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector '%s' installed!", c.store.ConnectorName)

	return nil
}
