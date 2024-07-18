package attempts

import (
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/printer"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListAbortedStore struct {
	Cursor shared.V2AttemptCursorResponseCursor `json:"attempts"`
	ErrorResponse error `json:"error"`
}
type ListAbortedController struct {
	store *ListAbortedStore
	cursorFlag string
}

func (c *ListAbortedController) GetStore() *ListAbortedStore {
	return c.store
}

var _ fctl.Controller[*ListAbortedStore] = (*ListAbortedController)(nil)


func (c *ListAbortedController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())
	cursor := fctl.GetString(cmd, c.cursorFlag)


	request := operations.GetWaitingAttemptsRequest{
		Cursor: &cursor,
	}

	response, err := store.Client().Webhooks.GetWaitingAttempts(cmd.Context(), request)
	if err!= nil {
		c.store.ErrorResponse = err
	} else {
		c.store.Cursor = response.V2AttemptCursorResponse.Cursor
	}

	return c, nil
}

func (c *ListAbortedController) Render(cmd *cobra.Command, args []string) error {
	
	
	if c.store.ErrorResponse != nil {
		pterm.Warning.WithShowLineNumber(false).Printfln(c.store.ErrorResponse.Error())
		return nil
	}

	tableData := fctl.Map(c.store.Cursor.Data, func(attempt shared.V2Attempt) []string{

		return []string{
			attempt.ID,
			string(attempt.Status),
			attempt.HookName,
			attempt.HookID,
			attempt.HookEndpoint,
			attempt.Event,
			attempt.DateOccured.Format(time.RFC3339),
			attempt.NextRetryAfter.Format(time.RFC3339),
			attempt.Payload,
		}

	})
	
	tableData = fctl.Prepend(tableData, []string{"ID", "Status", "Hook Name", "Hook ID", "Hook Endpoint", "Event", "Created At", "Next Try", "Payload"})
	
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

func NewListAbortedController() *ListAbortedController {
	return &ListAbortedController{
		store: &ListAbortedStore{},
		cursorFlag: "cursor",
	}
}

func NewListAbortedCommand() *cobra.Command {
	
	c := NewListAbortedController()

	return fctl.NewCommand("list-aborted",
		fctl.WithShortDescription("List all aborted attempts"),
		fctl.WithAliases("lsw", "lw"),
		fctl.WithStringFlag(c.cursorFlag, "", "Cursor pagination"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListAbortedStore](NewListAbortedController()),
	)
}
