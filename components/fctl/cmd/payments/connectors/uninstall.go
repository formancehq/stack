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
	Success   bool   `json:"success"`
	Connector string `json:"connector"`
}
type PaymentsConnectorsUninstallController struct {
	store           *PaymentsConnectorsUninstallStore
	providerFlag    string
	connectorIDFlag string
}

var _ fctl.Controller[*PaymentsConnectorsUninstallStore] = (*PaymentsConnectorsUninstallController)(nil)

func NewDefaultPaymentsConnectorsUninstallStore() *PaymentsConnectorsUninstallStore {
	return &PaymentsConnectorsUninstallStore{
		Success:   false,
		Connector: "",
	}
}

func NewPaymentsConnectorsUninstallController() *PaymentsConnectorsUninstallController {
	return &PaymentsConnectorsUninstallController{
		store:           NewDefaultPaymentsConnectorsUninstallStore(),
		providerFlag:    "provider",
		connectorIDFlag: "connectorid",
	}
}

func NewUninstallCommand() *cobra.Command {
	c := NewPaymentsConnectorsUninstallController()
	return fctl.NewCommand("uninstall",
		fctl.WithAliases("uninstall", "u", "un"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithValidArgs(internal.AllConnectors...),
		fctl.WithStringFlag(c.providerFlag, "", "Provider name"),
		fctl.WithStringFlag(c.connectorIDFlag, "", "Connector ID"),
		fctl.WithShortDescription("Uninstall a connector"),
		fctl.WithController[*PaymentsConnectorsUninstallStore](c),
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

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	provider := fctl.GetString(cmd, c.providerFlag)
	connectorID := fctl.GetString(cmd, c.connectorIDFlag)
	switch {
	case provider != "" && connectorID != "":
		if !fctl.CheckStackApprobation(cmd, stack, "You are about to uninstall connector '%s' from provider '%s'", connectorID, provider) {
			return nil, fctl.ErrMissingApproval
		}

		response, err := client.Payments.UninstallConnectorV1(cmd.Context(), operations.UninstallConnectorV1Request{
			ConnectorID: connectorID,
			Connector:   shared.Connector(provider),
		})
		if err != nil {
			return nil, err
		}

		if response.StatusCode >= 300 {
			return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}

		c.store.Connector = connectorID
	case provider == "" && connectorID == "":
		return nil, fmt.Errorf("must use either --provider and --connector-id")
	case connectorID != "":
		return nil, fmt.Errorf("must use --provider when using --connector-id")
	case provider != "":
		if !fctl.CheckStackApprobation(cmd, stack, "You are about to uninstall connector '%s'", provider) {
			return nil, fctl.ErrMissingApproval
		}

		response, err := client.Payments.UninstallConnector(cmd.Context(), operations.UninstallConnectorRequest{
			Connector: shared.Connector(provider),
		})
		if err != nil {
			return nil, err
		}

		if response.StatusCode >= 300 {
			return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
		}

		c.store.Connector = provider
	}

	c.store.Success = true

	return c, nil
}

// TODO: This need to use the ui.NewListModel
func (c *PaymentsConnectorsUninstallController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector '%s' uninstalled!", c.store.Connector)
	return nil
}
