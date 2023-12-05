package accounts

import (
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListBalancesStore struct {
	Cursor *shared.BalancesCursorCursor `json:"cursor"`
}

type ListBalancesController struct {
	store *ListBalancesStore
}

var _ fctl.Controller[*ListBalancesStore] = (*ListBalancesController)(nil)

func NewListBalanceStore() *ListBalancesStore {
	return &ListBalancesStore{
		Cursor: &shared.BalancesCursorCursor{},
	}
}

func NewListBalancesController() *ListBalancesController {
	return &ListBalancesController{
		store: NewListBalanceStore(),
	}
}

func (c *ListBalancesController) GetStore() *ListBalancesStore {
	return c.store
}

func (c *ListBalancesController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	response, err := client.Payments.GetAccountBalances(
		cmd.Context(),
		operations.GetAccountBalancesRequest{
			AccountID: args[0],
		},
	)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Cursor = &response.BalancesCursor.Cursor

	return c, nil
}

func (c *ListBalancesController) Render(cmd *cobra.Command, args []string) error {
	tableData := fctl.Map(c.store.Cursor.Data, func(balance shared.AccountBalance) []string {
		return []string{
			balance.AccountID,
			balance.Asset,
			balance.Balance.String(),
			balance.CreatedAt.Format(time.RFC3339),
			balance.LastUpdatedAt.Format(time.RFC3339),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"ID", "Asset", "Balance",
		"CreatedAt", "LastUpdatedAt"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}

func NewListBalanceCommand() *cobra.Command {
	return fctl.NewCommand("balances <accountID>",
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithShortDescription("List accounts balances"),
		fctl.WithController[*ListBalancesStore](NewListBalancesController()),
	)
}
