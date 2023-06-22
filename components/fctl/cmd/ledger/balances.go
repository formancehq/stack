package ledger

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type BalancesStore struct {
	Balances shared.BalancesCursorResponseCursor `json:"balances"`
}
type BalancesController struct {
	store       *BalancesStore
	addressFlag string
}

var _ fctl.Controller[*BalancesStore] = (*BalancesController)(nil)

func NewDefaultBalancesStore() *BalancesStore {
	return &BalancesStore{}
}

func NewBalancesController() *BalancesController {
	return &BalancesController{
		store:       NewDefaultBalancesStore(),
		addressFlag: "address",
	}
}

func NewBalancesCommand() *cobra.Command {
	c := NewBalancesController()
	return fctl.NewCommand("balances",
		fctl.WithAliases("balance", "bal", "b"),
		fctl.WithStringFlag(c.addressFlag, "", "Filter on specific address"),
		fctl.WithShortDescription("Read balances"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*BalancesStore](c),
	)
}

func (c *BalancesController) GetStore() *BalancesStore {
	return c.store
}

func (c *BalancesController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	response, err := client.Ledger.GetBalances(
		cmd.Context(),
		operations.GetBalancesRequest{
			Address: fctl.Ptr(fctl.GetString(cmd, c.addressFlag)),
			Ledger:  fctl.GetString(cmd, internal.LedgerFlag),
		},
	)
	if err != nil {
		return nil, err
	}

	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Balances = response.BalancesCursorResponse.Cursor

	return c, nil
}

func (c *BalancesController) Render(cmd *cobra.Command, args []string) error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{"Account", "Asset", "Balance"})
	for _, accountBalances := range c.store.Balances.Data {
		for account, volumes := range accountBalances {
			for asset, balance := range volumes {
				tableData = append(tableData, []string{
					account, asset, fmt.Sprint(balance),
				})
			}
		}
	}

	return pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}
