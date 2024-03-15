package ledger

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type StatsStore struct {
	Stats shared.Stats `json:"stats"`
}
type StatsController struct {
	store *StatsStore
}

var _ fctl.Controller[*StatsStore] = (*StatsController)(nil)

func NewDefaultStatsStore() *StatsStore {
	return &StatsStore{}
}

func NewStatsController() *StatsController {
	return &StatsController{
		store: NewDefaultStatsStore(),
	}
}

func NewStatsCommand() *cobra.Command {
	return fctl.NewCommand("stats",
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithAliases("st"),
		fctl.WithShortDescription("Read ledger stats"),
		fctl.WithController[*StatsStore](NewStatsController()),
	)
}

func (c *StatsController) GetStore() *StatsStore {
	return c.store
}

func (c *StatsController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	request := operations.ReadStatsRequest{
		Ledger: fctl.GetString(cmd, internal.LedgerFlag),
	}
	response, err := store.Client().Ledger.ReadStats(cmd.Context(), request)
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

// TODO: This need to use the ui.NewListModel
func (c *StatsController) Render(cmd *cobra.Command, args []string) error {

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Transactions"), fmt.Sprint(c.store.Stats.Transactions)})
	tableData = append(tableData, []string{pterm.LightCyan("Accounts"), fmt.Sprint(c.store.Stats.Accounts)})

	return pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}
