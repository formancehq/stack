package ledger

import (
	"flag"
	"fmt"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	addressFlag   = "address"
	useBalances   = "balances"
	shortBalances = "Read balances"
)

type BalancesStore struct {
	Balances shared.BalancesCursorResponseCursor `json:"balances"`
}

func NewBalancesStore() *BalancesStore {
	return &BalancesStore{}
}

func NewBalancesConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useBalances, flag.ExitOnError)
	flags.String(addressFlag, "", "Filter on specific address")
	return fctl.NewControllerConfig(
		useBalances,
		shortBalances,
		shortBalances,
		[]string{
			"balance", "bal", "b",
		},
		flags,
		fctl.Organization, fctl.Stack, fctl.Ledger,
	)
}

type BalancesController struct {
	store  *BalancesStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*BalancesStore] = (*BalancesController)(nil)

func NewBalancesController(config *fctl.ControllerConfig) *BalancesController {
	return &BalancesController{
		store:  NewBalancesStore(),
		config: config,
	}
}

func (c *BalancesController) GetStore() *BalancesStore {
	return c.store
}

func (c *BalancesController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *BalancesController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	out := c.config.GetOut()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	response, err := client.Ledger.GetBalances(
		ctx,
		operations.GetBalancesRequest{
			Address: fctl.Ptr(fctl.GetString(flags, addressFlag)),
			Ledger:  fctl.GetString(flags, internal.LedgerFlag),
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

func (c *BalancesController) Render() error {
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
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()
}

func NewBalancesCommand() *cobra.Command {
	c := NewBalancesConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*BalancesStore](NewBalancesController(c)),
	)
}
