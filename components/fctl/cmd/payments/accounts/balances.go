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

	cursorFlag   string
	pageSizeFlag string
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

		cursorFlag:   "cursor",
		pageSizeFlag: "page-size",
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

	var cursor *string
	if c := fctl.GetString(cmd, c.cursorFlag); c != "" {
		cursor = &c
	}

	var pageSize *int64
	if ps := fctl.GetInt(cmd, c.pageSizeFlag); ps > 0 {
		pageSize = fctl.Ptr(int64(ps))
	}

	response, err := client.Payments.GetAccountBalances(
		cmd.Context(),
		operations.GetAccountBalancesRequest{
			Cursor:    cursor,
			PageSize:  pageSize,
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
	if err := pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	tableData = pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("HasMore"), fmt.Sprintf("%v", c.store.Cursor.HasMore)})
	tableData = append(tableData, []string{pterm.LightCyan("PageSize"), fmt.Sprintf("%d", c.store.Cursor.PageSize)})
	tableData = append(tableData, []string{pterm.LightCyan("Next"), func() string {
		if c.store.Cursor.Next == nil {
			return ""
		}
		return *c.store.Cursor.Next
	}()})
	tableData = append(tableData, []string{pterm.LightCyan("Previous"), func() string {
		if c.store.Cursor.Previous == nil {
			return ""
		}
		return *c.store.Cursor.Previous
	}()})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	return nil
}

func NewListBalanceCommand() *cobra.Command {
	c := NewListBalancesController()
	return fctl.NewCommand("balances <accountID>",
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithShortDescription("List accounts balances"),
		fctl.WithStringFlag(c.cursorFlag, "", "Cursor"),
		fctl.WithIntFlag(c.pageSizeFlag, 0, "PageSize"),
		fctl.WithController[*ListBalancesStore](c),
	)
}
