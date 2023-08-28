package payments

import (
	"flag"
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useList         = "list"
	descriptionList = "List all payments"
	shortList       = "List all payments"
)

type ListStore struct {
	Cursor *shared.PaymentsCursorCursor `json:"cursor"`
}

func NewListStore() *ListStore {
	return &ListStore{
		Cursor: &shared.PaymentsCursorCursor{},
	}
}

func NewListConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)

	return fctl.NewControllerConfig(
		useList,
		descriptionList,
		shortList,
		[]string{
			"list",
			"ls",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)

}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

type ListController struct {
	store  *ListStore
	config *fctl.ControllerConfig
}

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

	response, err := client.Payments.ListPayments(
		ctx,
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

func (c *ListController) Render() error {
	tableData := fctl.Map(c.store.Cursor.Data, func(payment shared.Payment) []string {
		return []string{
			payment.ID,
			string(payment.Type),
			fmt.Sprint(payment.InitialAmount),
			payment.Asset,
			string(payment.Status),
			string(payment.Scheme),
			payment.Reference,
			payment.SourceAccountID,
			payment.DestinationAccountID,
			string(payment.Provider),
			payment.CreatedAt.Format(time.RFC3339),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"ID", "Type", "Amount", "Asset", "Status",
		"Scheme", "Reference", "Source Account ID", "Destination Account ID", "Provider", "Created at"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()
}

func NewListPaymentsCommand() *cobra.Command {

	config := NewListConfig()

	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController(config)),
	)
}
