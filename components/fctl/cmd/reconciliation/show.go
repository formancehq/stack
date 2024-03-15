package reconciliation

import (
	"fmt"
	"math/big"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ShowStore struct {
	Reconciliation shared.Reconciliation `json:"policy"`
}
type ShowController struct {
	store *ShowStore
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

func NewShowStore() *ShowStore {
	return &ShowStore{}
}

func NewShowController() *ShowController {
	return &ShowController{
		store: NewShowStore(),
	}
}

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("get <reconciliationID>",
		fctl.WithShortDescription("Get reconciliation"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithAliases("sh", "s"),
		fctl.WithController[*ShowStore](NewShowController()),
	)
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	response, err := store.Client().Reconciliation.GetReconciliation(cmd.Context(), operations.GetReconciliationRequest{
		ReconciliationID: args[0],
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	if response.ReconciliationResponse == nil {
		return nil, fmt.Errorf("policy not found")
	}

	c.store.Reconciliation = response.ReconciliationResponse.Data

	return c, nil
}

func (c *ShowController) Render(cmd *cobra.Command, args []string) error {
	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.Reconciliation.ID})
	tableData = append(tableData, []string{pterm.LightCyan("PolicyID"), c.store.Reconciliation.PolicyID})
	tableData = append(tableData, []string{pterm.LightCyan("CreatedAt"), c.store.Reconciliation.CreatedAt.Format(time.RFC3339)})
	tableData = append(tableData, []string{pterm.LightCyan("ReconciledAtLedger"), c.store.Reconciliation.ReconciledAtLedger.Format(time.RFC3339)})
	tableData = append(tableData, []string{pterm.LightCyan("ReconciledAtPayments"), c.store.Reconciliation.ReconciledAtPayments.Format(time.RFC3339)})
	tableData = append(tableData, []string{pterm.LightCyan("Status"), c.store.Reconciliation.Status})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Ledger Balances")
	tableData = fctl.MapMap(c.store.Reconciliation.LedgerBalances, func(key string, value *big.Int) []string {
		return []string{
			key,
			value.String(),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"Asset", "Amount"})
	if err := pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Payments Balances")
	tableData = fctl.MapMap(c.store.Reconciliation.PaymentsBalances, func(key string, value *big.Int) []string {
		return []string{
			key,
			value.String(),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"Asset", "Amount"})
	if err := pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Drift Balances")
	tableData = fctl.MapMap(c.store.Reconciliation.DriftBalances, func(key string, value *big.Int) []string {
		return []string{
			key,
			value.String(),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"Asset", "Amount"})
	if err := pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	return nil
}
