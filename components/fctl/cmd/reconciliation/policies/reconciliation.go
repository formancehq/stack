package policies

import (
	"fmt"
	"math/big"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ReconciliationStore struct {
	Reconciliation shared.Reconciliation `json:"reconciliation"`
}
type ReconciliationController struct {
	store *ReconciliationStore
}

var _ fctl.Controller[*ReconciliationStore] = (*ReconciliationController)(nil)

func NewReconciliationStore() *ReconciliationStore {
	return &ReconciliationStore{}
}

func NewReconciliationController() *ReconciliationController {
	return &ReconciliationController{
		store: NewReconciliationStore(),
	}
}

func NewReconciliationCommand() *cobra.Command {
	return fctl.NewCommand("reconcile <policyID> <atLedger> <atPayments>",
		fctl.WithShortDescription("Launch a reconciliation from a policy"),
		fctl.WithArgs(cobra.ExactArgs(3)),
		fctl.WithAliases("r"),
		fctl.WithController[*ReconciliationStore](NewReconciliationController()),
	)
}

func (c *ReconciliationController) GetStore() *ReconciliationStore {
	return c.store
}

func (c *ReconciliationController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	ledgerClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	atLedger, err := time.Parse(time.RFC3339, args[1])
	if err != nil {
		return nil, err
	}

	atPayments, err := time.Parse(time.RFC3339, args[2])
	if err != nil {
		return nil, err
	}

	fmt.Println(args[0], atLedger, atPayments)
	response, err := ledgerClient.Reconciliation.Reconcile(cmd.Context(), operations.ReconcileRequest{
		PolicyID: args[0],
		ReconciliationRequest: shared.ReconciliationRequest{
			ReconciledAtLedger:   atLedger,
			ReconciledAtPayments: atPayments,
		},
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

func (c *ReconciliationController) Render(cmd *cobra.Command, args []string) error {
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
