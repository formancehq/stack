package workflows

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
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
	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config")
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
		return nil, errors.Wrap(err, "creating stack client")
	}

	response, err := client.Orchestration.DeleteWorkflow(
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
