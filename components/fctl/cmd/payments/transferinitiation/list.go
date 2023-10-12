package transferinitiation

import (
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStore struct {
	Cursor *shared.TransferInitiationsCursorCursor `json:"cursor"`
}

type ListController struct {
	store *ListStore
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewListStore() *ListStore {
	return &ListStore{
		Cursor: &shared.TransferInitiationsCursorCursor{},
	}
}

func NewListController() *ListController {
	return &ListController{
		store: NewListStore(),
	}
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

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

	response, err := client.Payments.ListTransferInitiations(
		cmd.Context(),
		operations.ListTransferInitiationsRequest{},
	)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Cursor = &response.TransferInitiationsCursor.Cursor

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	tableData := fctl.Map(c.store.Cursor.Data, func(tf shared.TransferInitiation) []string {
		return []string{
			tf.ID,
			tf.CreatedAt.Format(time.RFC3339),
			tf.UpdatedAt.Format(time.RFC3339),
			tf.ScheduledAt.Format(time.RFC3339),
			tf.Description,
			tf.SourceAccountID,
			tf.DestinationAccountID,
			string(tf.Provider),
			string(tf.Type),
			fmt.Sprint(tf.Amount),
			tf.Asset,
			string(tf.Status),
			tf.Error,
		}
	})
	tableData = fctl.Prepend(tableData, []string{"ID", "CreatedAt", "UpdatedAt", "ScheduledAt", "Description", "Source Account ID",
		"Destination Account ID", "Provider", "Type", "Amount", "Asset", "Status", "Error"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}

func NewListCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithShortDescription("List transfer initiation"),
		fctl.WithController[*ListStore](NewListController()),
	)
}
