package balances

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStore struct {
	BalancesNames [][]string `json:"balancesNames"`
}
type ListController struct {
	store *ListStore
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewDefaultListStore() *ListStore {
	return &ListStore{}
}

func NewListController() *ListController {
	return &ListController{
		store: NewDefaultListStore(),
	}
}

func NewListCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List balances"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		internal.WithTargetingWalletByName(),
		internal.WithTargetingWalletByID(),
		fctl.WithController[*ListStore](NewListController()),
	)
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	walletID, err := internal.RequireWalletID(cmd, store.Client())
	if err != nil {
		return nil, err
	}

	request := operations.ListBalancesRequest{
		ID: walletID,
	}
	response, err := store.Client().Wallets.ListBalances(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "listing balance")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.BalancesNames = fctl.Map(response.ListBalancesResponse.Cursor.Data, func(balance shared.Balance) []string {
		return []string{
			balance.Name,
		}
	})

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	if len(c.store.BalancesNames) == 0 {
		fctl.Println("No balances found.")
		return nil
	}

	tableData := fctl.Prepend(c.store.BalancesNames, []string{"Name"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()

}
