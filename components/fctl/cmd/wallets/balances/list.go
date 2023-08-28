package balances

import (
	"flag"
	"fmt"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useList   = "list"
	shortList = "List balances"
)

type ListStore struct {
	BalancesNames [][]string `json:"balancesNames"`
}

func NewListStore() *ListStore {
	return &ListStore{}
}
func NewListConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	internal.WithTargetingWalletByName(flags)
	internal.WithTargetingWalletByID(flags)
	return fctl.NewControllerConfig(
		useList,
		shortList,
		shortList,
		[]string{
			"ls", "l",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

type ListController struct {
	store  *ListStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewListController(config *fctl.ControllerConfig) *ListController {
	return &ListController{
		store:  NewListStore(),
		config: config,
	}
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *ListController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	out := c.config.GetOut()
	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config")
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	walletID, err := internal.RequireWalletID(flags, ctx, client)
	if err != nil {
		return nil, err
	}

	request := operations.ListBalancesRequest{
		ID: walletID,
	}
	response, err := client.Wallets.ListBalances(ctx, request)
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

func (c *ListController) Render() error {
	if len(c.store.BalancesNames) == 0 {
		fmt.Fprintln(c.config.GetOut(), "No balances found.")
		return nil
	}

	tableData := fctl.Prepend(c.store.BalancesNames, []string{"Name"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()

}

func NewListCommand() *cobra.Command {
	c := NewListConfig()

	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController(c)),
	)
}
