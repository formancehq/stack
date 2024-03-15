package connectors

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
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
	PaymentsVersion versions.Version

	store *PaymentsConnectorsListStore
}

func (c *PaymentsConnectorsListController) SetVersion(version versions.Version) {
	c.PaymentsVersion = version
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
	c := NewPaymentsConnectorsListController()
	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List all enabled connectors"),
		fctl.WithController[*PaymentsConnectorsListStore](c),
	)
}

func (c *PaymentsConnectorsListController) GetStore() *PaymentsConnectorsListStore {
	return c.store
}

func (c *PaymentsConnectorsListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	if err := versions.GetPaymentsVersion(cmd, args, c); err != nil {
		return nil, err
	}

	response, err := store.Client().Payments.ListAllConnectors(cmd.Context())
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
		switch c.PaymentsVersion {
		case versions.V1:
			return []string{
				string(connector.Provider),
				connector.Name,
				connector.ConnectorID,
			}
		default:
			// V0
			return []string{
				string(connector.Provider),
			}
		}

	})
	switch c.PaymentsVersion {
	case versions.V1:
		tableData = fctl.Prepend(tableData, []string{"Provider", "Name", "ConnectorID"})
	default:
		tableData = fctl.Prepend(tableData, []string{"Provider"})
	}

	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}
