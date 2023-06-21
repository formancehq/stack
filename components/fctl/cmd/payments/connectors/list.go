package connectors

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	PaymentsConnectorsList = "develop"
)

type PaymentsConnectorsListStore struct {
	Connectors []shared.ConnectorsResponseData `json:"connectors"`
}
type PaymentsConnectorsListController struct {
	store *PaymentsConnectorsListStore
}

var _ fctl.Controller[*PaymentsConnectorsListStore] = (*PaymentsConnectorsListController)(nil)

func NewDefaultPaymentsConnectorsListStore() *PaymentsConnectorsListStore {
	return &PaymentsConnectorsListStore{
		Connectors: []shared.ConnectorsResponseData{},
	}
}

func NewPaymentsConnectorsListController() *PaymentsConnectorsListController {
	return &PaymentsConnectorsListController{
		store: NewDefaultPaymentsConnectorsListStore(),
	}
}

func NewListCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List all enabled connectors"),
		fctl.WithController[*PaymentsConnectorsListStore](NewPaymentsConnectorsListController()),
	)
}

func (c *PaymentsConnectorsListController) GetStore() *PaymentsConnectorsListStore {
	return c.store
}

func (c *PaymentsConnectorsListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

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

	response, err := client.Payments.ListAllConnectors(cmd.Context())
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	if response.ConnectorsResponse == nil {
		return nil, fmt.Errorf("unexpected response: %v", response)
	}

	c.store.Connectors = response.ConnectorsResponse.Data

	return c, nil
}

func (c *PaymentsConnectorsListController) Render(cmd *cobra.Command, args []string) error {
	tableData := fctl.Map(c.store.Connectors, func(connector shared.ConnectorsResponseData) []string {
		return []string{
			string(*connector.Provider),
			fctl.BoolToString(*connector.Enabled),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"Provider", "Enabled"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}
