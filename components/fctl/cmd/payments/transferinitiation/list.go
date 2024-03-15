package transferinitiation

import (
	"fmt"
	"time"

	"github.com/formancehq/fctl/cmd/payments/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStore struct {
	Cursor *shared.TransferInitiationsCursorCursor `json:"cursor"`
}

type ListController struct {
	PaymentsVersion versions.Version

	store *ListStore

	cursorFlag   string
	pageSizeFlag string
}

func (c *ListController) SetVersion(version versions.Version) {
	c.PaymentsVersion = version
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

		cursorFlag:   "cursor",
		pageSizeFlag: "page-size",
	}
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())
	if err := versions.GetPaymentsVersion(cmd, args, c); err != nil {
		return nil, err
	}

	if c.PaymentsVersion < versions.V1 {
		return nil, fmt.Errorf("transfer initiation are only supported in >= v1.0.0")
	}

	var cursor *string
	if c := fctl.GetString(cmd, c.cursorFlag); c != "" {
		cursor = &c
	}

	var pageSize *int64
	if ps := fctl.GetInt(cmd, c.pageSizeFlag); ps > 0 {
		pageSize = fctl.Ptr(int64(ps))
	}

	response, err := store.Client().Payments.ListTransferInitiations(
		cmd.Context(),
		operations.ListTransferInitiationsRequest{
			Cursor:   cursor,
			PageSize: pageSize,
		},
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
			tf.Reference,
			tf.CreatedAt.Format(time.RFC3339),
			tf.ScheduledAt.Format(time.RFC3339),
			tf.Description,
			tf.SourceAccountID,
			tf.DestinationAccountID,
			tf.ConnectorID,
			string(tf.Type),
			fmt.Sprint(tf.Amount),
			fmt.Sprint(tf.InitialAmount),
			tf.Asset,
			string(tf.Status),
			tf.Error,
		}
	})
	tableData = fctl.Prepend(tableData, []string{"ID", "Reference", "CreatedAt", "ScheduledAt", "Description", "Source Account ID",
		"Destination Account ID", "ConnectorID", "Type", "Amount", "InitialAmount", "Asset", "Status", "Error"})
	if err := pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	tableData = pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("HasMore"), fmt.Sprintf("%v", c.store.Cursor.HasMore)})
	tableData = append(tableData, []string{pterm.LightCyan("PageSize"), fmt.Sprintf("%d", c.store.Cursor.PageSize)})
	tableData = append(tableData, []string{pterm.LightCyan("Next"), func() string {
		if c.store.Cursor.Next == nil {
			return ""
		}
		return *c.store.Cursor.Next
	}()})
	tableData = append(tableData, []string{pterm.LightCyan("Previous"), func() string {
		if c.store.Cursor.Previous == nil {
			return ""
		}
		return *c.store.Cursor.Previous
	}()})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	return nil
}

func NewListCommand() *cobra.Command {
	c := NewListController()
	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithShortDescription("List transfer initiation"),
		fctl.WithStringFlag(c.cursorFlag, "", "Cursor"),
		fctl.WithIntFlag(c.pageSizeFlag, 0, "PageSize"),
		fctl.WithController[*ListStore](c),
	)
}
