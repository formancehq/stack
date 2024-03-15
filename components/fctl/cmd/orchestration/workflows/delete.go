package workflows

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type WorkflowsDeleteStore struct {
	WorkflowId string `json:"workflowId"`
	Success    bool   `json:"success"`
}
type WorkflowsDeleteController struct {
	store *WorkflowsDeleteStore
}

var _ fctl.Controller[*WorkflowsDeleteStore] = (*WorkflowsDeleteController)(nil)

func NewDefaultWorkflowsDeleteStore() *WorkflowsDeleteStore {
	return &WorkflowsDeleteStore{}
}

func NewWorkflowsDeleteController() *WorkflowsDeleteController {
	return &WorkflowsDeleteController{
		store: NewDefaultWorkflowsDeleteStore(),
	}
}
func NewDeleteCommand() *cobra.Command {
	return fctl.NewCommand("delete <workflow-id>",
		fctl.WithAliases("del", "d"),
		fctl.WithShortDescription("Soft delete a workflow"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*WorkflowsDeleteStore](NewWorkflowsDeleteController()),
	)
}

func (c *WorkflowsDeleteController) GetStore() *WorkflowsDeleteStore {
	return c.store
}

func (c *WorkflowsDeleteController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	response, err := store.Client().Orchestration.DeleteWorkflow(
		cmd.Context(),
		operations.DeleteWorkflowRequest{
			FlowID: args[0],
		},
	)

	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("%s: %s", response.Error.ErrorCode, response.Error.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.WorkflowId = args[0]
	c.store.Success = true

	return c, nil
}

func (c *WorkflowsDeleteController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithShowLineNumber().Printfln("Workflow %s Deleted!", c.store.WorkflowId)
	return nil
}
