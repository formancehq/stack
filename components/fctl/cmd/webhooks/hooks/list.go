package hooks

import (
	
	"strings"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/printer"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var cursorFlag string = "cursor"

type ListWebhookStore struct {
	Cursor shared.V2HookCursorResponseCursor `json:"cursor"`
	ErrorResponse  error `json:"error"`
}
type ListWebhookController struct {
	store *ListWebhookStore
}

var _ fctl.Controller[*ListWebhookStore] = (*ListWebhookController)(nil)

func NewDefaultListWebhookStore() *ListWebhookStore {
	return &ListWebhookStore{
		Cursor: shared.V2HookCursorResponseCursor{},
	}
}

func NewListWebhookController() *ListWebhookController {
	return &ListWebhookController{
		store: NewDefaultListWebhookStore(),
	}
}
func (c *ListWebhookController) GetStore() *ListWebhookStore {
	return c.store
}

func (c *ListWebhookController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())
	cursor := fctl.GetString(cmd, cursorFlag)

	request := operations.GetManyHooksRequest{
		Cursor: &cursor,
	}

	response, err := store.Client().Webhooks.GetManyHooks(cmd.Context(), request)
	
	
	if err!= nil {
		c.store.ErrorResponse = err
	} else {
		c.store.Cursor = response.V2HookCursorResponse.Cursor
	}
	

	return c, nil
}

func (c *ListWebhookController) Render(cmd *cobra.Command, args []string) error {
	if c.store.ErrorResponse != nil {
		pterm.Warning.WithShowLineNumber(false).Printfln(c.store.ErrorResponse.Error())
		return nil
	}

	tableData := fctl.Map(c.store.Cursor.Data, func(hook shared.V2Hook) []string {
		return []string{
			hook.ID,
			hook.Name,
			hook.Endpoint,
			strings.Join(hook.Events, ","),
			hook.Secret,
			string(hook.Status),
			fctl.BoolToString(hook.Retry),
		}
	})

	tableData = fctl.Prepend(tableData, []string{"ID", "Name", "Endpoint", "Events", "Secret", "Status", "Retry"})

	tableData = printer.AddCursorRowsToTable(tableData, printer.CursorArgs{
		HasMore : c.store.Cursor.HasMore,
		Next: &c.store.Cursor.Next,
		PageSize: c.store.Cursor.PageSize,
		Previous: &c.store.Cursor.Previous,
	}) 
		
	writer := cmd.OutOrStdout()

	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(writer).
		WithData(tableData).
		Render()
}

func NewListCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithShortDescription("List all Hooks"),
		fctl.WithAliases("ls", "l"),
		fctl.WithStringFlag(cursorFlag, "", "Cursor pagination"),
		fctl.WithController[*ListWebhookStore](NewListWebhookController()),
	)
}
