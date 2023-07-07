package transactions

import (
	"fmt"
	"time"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStore struct {
	Transactions []shared.WalletsTransaction `json:"transactions"`
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
		fctl.WithShortDescription("List transactions"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController()),
	)
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "retriecing config")
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
		return nil, errors.Wrap(err, "creating stack client")
	}

	walletID, err := internal.RetrieveWalletID(cmd, client)
	if err != nil {
		return nil, err
	}

	request := operations.GetTransactionsRequest{
		WalletID: &walletID,
	}
	response, err := client.Wallets.GetTransactions(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "listing transactions")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Transactions = response.GetTransactionsResponse.Cursor.Data

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	if len(c.store.Transactions) == 0 {
		fctl.Println("No transactions found.")
		return nil
	}

	tableData := fctl.Map(c.store.Transactions, func(tx shared.WalletsTransaction) []string {
		return []string{
			fmt.Sprintf("%d", tx.Txid),
			tx.Timestamp.Format(time.RFC3339),
			fctl.MetadataAsShortString(tx.Metadata),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"ID", "Date", "Metadata"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()

}
