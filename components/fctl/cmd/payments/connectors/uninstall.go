package connectors

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/connectors/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	PaymentsConnectorsUninstall = "develop"
)

type PaymentsConnectorsUninstallStore struct {
	Success       bool   `json:"success"`
	ConnectorName string `json:"connectorName"`
}
type PaymentsConnectorsUninstallController struct {
	store *PaymentsConnectorsUninstallStore
}

var _ fctl.Controller[*PaymentsConnectorsUninstallStore] = (*PaymentsConnectorsUninstallController)(nil)

func NewDefaultPaymentsConnectorsUninstallStore() *PaymentsConnectorsUninstallStore {
	return &PaymentsConnectorsUninstallStore{
		Success:       false,
		ConnectorName: "",
	}
}

func NewPaymentsConnectorsUninstallController() *PaymentsConnectorsUninstallController {
	return &PaymentsConnectorsUninstallController{
		store: NewDefaultPaymentsConnectorsUninstallStore(),
	}
}

func NewUninstallCommand() *cobra.Command {
	return fctl.NewCommand("uninstall <connector-name>",
		fctl.WithAliases("uninstall", "u", "un"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithValidArgs(internal.AllConnectors...),
		fctl.WithShortDescription("Uninstall a connector"),
		fctl.WithController[*PaymentsConnectorsUninstallStore](NewPaymentsConnectorsUninstallController()),
	)
}

func (c *PaymentsConnectorsUninstallController) GetStore() *PaymentsConnectorsUninstallStore {
	return c.store
}

func (c *PaymentsConnectorsUninstallController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to uninstall connector '%s'", args[0]) {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	response, err := client.Payments.UninstallConnector(cmd.Context(), operations.UninstallConnectorRequest{
		Connector: shared.Connector(args[0]),
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	c.store.ConnectorName = args[0]

	return c, nil
}

// TODO: This need to use the ui.NewListModel
func (c *PaymentsConnectorsUninstallController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector '%s' uninstalled!", c.store.ConnectorName)
	return nil
}
