package ledger

import (
	"fmt"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStore struct {
	Ledgers []shared.V2Ledger `json:"ledgers"`
}
type ListController struct {
	store *ListStore
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewDefaultListStore() *ListStore {
	return &ListStore{
		Ledgers: []shared.V2Ledger{},
	}
}

func NewListController() *ListController {
	return &ListController{
		store: NewDefaultListStore(),
	}
}

func NewListCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithAliases("l", "ls"),
		fctl.WithShortDescription("List ledgers (starting from ledger v2)"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController()),
	)
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	response, err := store.Client().Ledger.V2ListLedgers(cmd.Context(), operations.V2ListLedgersRequest{})
	if err != nil {
		return nil, err
	}

	if response.V2ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.V2ErrorResponse.ErrorCode, response.V2ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Ledgers = response.V2LedgerListResponse.Cursor.Data

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	tableData := fctl.Map(c.store.Ledgers, func(ledger shared.V2Ledger) []string {
		return []string{
			ledger.Name,
		}
	})
	tableData = fctl.Prepend(tableData, []string{"Name"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}
