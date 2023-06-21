package payments

import (
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type PaymentsListStore struct {
	Cursor *shared.PaymentsCursorCursor `json:"cursor"`
}
type PaymentsListController struct {
	store *PaymentsListStore
}

var _ fctl.Controller[*PaymentsListStore] = (*PaymentsListController)(nil)

func NewDefaultPaymentsListStore() *PaymentsListStore {
	return &PaymentsListStore{
		Cursor: &shared.PaymentsCursorCursor{},
	}
}

func NewPaymentsListController() *PaymentsListController {
	return &PaymentsListController{
		store: NewDefaultPaymentsListStore(),
	}
}

func (c *PaymentsListController) GetStore() *PaymentsListStore {
	return c.store
}

func (c *PaymentsListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

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

	response, err := client.Payments.ListPayments(
		cmd.Context(),
		operations.ListPaymentsRequest{},
	)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Cursor = &response.PaymentsCursor.Cursor

	return c, nil
}

func (c *PaymentsListController) Render(cmd *cobra.Command, args []string) error {
	tableData := fctl.Map(c.store.Cursor.Data, func(payment shared.Payment) []string {
		return []string{
			payment.ID,
			string(payment.Type),
			fmt.Sprint(payment.InitialAmount),
			payment.Asset,
			string(payment.Status),
			string(payment.Scheme),
			payment.Reference,
			payment.AccountID,
			string(payment.Provider),
			payment.CreatedAt.Format(time.RFC3339),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"ID", "Type", "Amount", "Asset", "Status",
		"Scheme", "Reference", "Account ID", "Provider", "Created at"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}

func NewListPaymentsCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithAliases("ls"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithShortDescription("List payments"),
		fctl.WithController[*PaymentsListStore](NewPaymentsListController()),
	)
}
