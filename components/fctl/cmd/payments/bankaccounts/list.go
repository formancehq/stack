package bankaccounts

import (
	"fmt"
	"time"

	"github.com/formancehq/fctl/cmd/payments/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStore struct {
	Cursor *shared.BankAccountsCursorCursor `json:"cursor"`
}

type ListController struct {
	PaymentsVersion versions.Version

	store *ListStore

	cursorFlag   string
	pageSizeFlag string
}

func (c *ListController) SetVersion(version versions.Version) {
	c.PaymentsVersion = version
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewListStore() *ListStore {
	return &ListStore{
		Cursor: &shared.BankAccountsCursorCursor{},
	}
}

func NewListController() *ListController {
	return &ListController{
		store: NewListStore(),

		cursorFlag:   "cursor",
		pageSizeFlag: "page-size",
	}
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	if err := versions.GetPaymentsVersion(cmd, args, c); err != nil {
		return nil, err
	}

	if c.PaymentsVersion < versions.V1 {
		return nil, fmt.Errorf("bank accounts are only supported in >= v1.0.0")
	}

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

	response, err := client.Payments.ListBankAccounts(
		cmd.Context(),
		operations.ListBankAccountsRequest{
			Cursor:   cursor,
			PageSize: pageSize,
		},
	)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Cursor = &response.BankAccountsCursor.Cursor

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	tableData := fctl.Map(c.store.Cursor.Data, func(bc shared.BankAccount) []string {
		return []string{
			bc.ID,
			bc.Name,
			bc.CreatedAt.Format(time.RFC3339),
			bc.Country,
		}
	})
	tableData = fctl.Prepend(tableData, []string{"ID", "Name", "CreatedAt", "Country"})
	if err := pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	return nil
}

func NewListCommand() *cobra.Command {
	c := NewListController()
	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithShortDescription("List bank accounts"),
		fctl.WithStringFlag(c.cursorFlag, "", "Cursor"),
		fctl.WithIntFlag(c.pageSizeFlag, 0, "PageSize"),
		fctl.WithController[*ListStore](c),
	)
}
