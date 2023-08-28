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
	useStats         = "stats"
	descriptionStats = "Read ledger stats"
)

type StatsStore struct {
	Stats shared.Stats `json:"stats"`
}

func NewStatsStore() *StatsStore {
	return &StatsStore{}
}

func NewStatsConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useStats, flag.ExitOnError)

	return fctl.NewControllerConfig(
		useStats,
		descriptionStats,
		descriptionStats,
		[]string{
			"st",
		},
		flags,
		fctl.Organization, fctl.Stack, fctl.Ledger,
	)
}

var _ fctl.Controller[*StatsStore] = (*StatsController)(nil)

type StatsController struct {
	store  *StatsStore
	config *fctl.ControllerConfig
}

func NewStatsController(config *fctl.ControllerConfig) *StatsController {
	return &StatsController{
		store:  NewStatsStore(),
		config: config,
	}
}

func (c *StatsController) GetStore() *StatsStore {
	return c.store
}

func (c *StatsController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *StatsController) Run() (fctl.Renderable, error) {
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

	ledgerClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	request := operations.ReadStatsRequest{
		Ledger: fctl.GetString(flags, internal.LedgerFlag),
	}
	response, err := ledgerClient.Ledger.ReadStats(ctx, request)
	if err != nil {
		return nil, err
	}

	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Stats = response.StatsResponse.Data

	return c, nil
}

func (c *StatsController) Render() error {

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Transactions"), fmt.Sprint(c.store.Stats.Transactions)})
	tableData = append(tableData, []string{pterm.LightCyan("Accounts"), fmt.Sprint(c.store.Stats.Accounts)})

	return pterm.DefaultTable.
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()
}

func NewStatsCommand() *cobra.Command {
	config := NewStatsConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*StatsStore](NewStatsController(config)),
	)
}
