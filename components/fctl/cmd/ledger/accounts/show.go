package accounts

import (
	"flag"
	"fmt"

	internal "github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useShow   = "show <address>"
	shortShow = "Show account"
)

type ShowStore struct {
	Account *shared.AccountWithVolumesAndBalances `json:"account"`
}

func NewShowStore() *ShowStore {
	return &ShowStore{}
}

func NewShowConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useShow, flag.ExitOnError)
	return fctl.NewControllerConfig(
		useShow,
		shortShow,
		shortShow,
		[]string{
			"sh", "s",
		},
		flags, fctl.Organization, fctl.Stack, fctl.Ledger,
	)
}

type ShowController struct {
	store  *ShowStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

func NewShowController(config *fctl.ControllerConfig) *ShowController {
	return &ShowController{
		store:  NewShowStore(),
		config: config,
	}
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *ShowController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()
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

	ledgerClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	ledger := fctl.GetString(flags, internal.LedgerFlag)
	response, err := ledgerClient.Ledger.GetAccount(ctx, operations.GetAccountRequest{
		Address: args[0],
		Ledger:  ledger,
	})
	if err != nil {
		return nil, err
	}

	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Account = &response.AccountResponse.Data

	return c, nil
}

func (c *ShowController) Render() error {

	out := c.config.GetOut()

	fctl.Section.WithWriter(out).Println("Information")
	if c.store.Account.Volumes != nil && len(c.store.Account.Volumes) > 0 {
		tableData := pterm.TableData{}
		tableData = append(tableData, []string{"Asset", "Input", "Output"})
		for asset, volumes := range c.store.Account.Volumes {
			input := volumes["input"]
			output := volumes["output"]
			tableData = append(tableData, []string{pterm.LightCyan(asset), fmt.Sprint(input), fmt.Sprint(output)})
		}
		if err := pterm.DefaultTable.
			WithHasHeader(true).
			WithWriter(out).
			WithData(tableData).
			Render(); err != nil {
			return err
		}
	} else {
		fmt.Fprintln(out, "No balances.")
	}

	fmt.Fprintln(out)

	return fctl.PrintMetadata(out, c.store.Account.Metadata)
}

func NewShowCommand() *cobra.Command {

	config := NewShowConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*ShowStore](NewShowController(config)),
	)
}
