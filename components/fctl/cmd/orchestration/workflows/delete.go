package workflows

import (
	"flag"
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useDelete   = "delete <workflow-id>"
	shortDelete = "Delete a workflow"
)

type DeleteStore struct {
	WorkflowId string `json:"workflowId"`
	Success    bool   `json:"success"`
}

func NewDeleteStore() *DeleteStore {
	return &DeleteStore{}
}

func NewDeleteConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useDelete, flag.ExitOnError)

	c := fctl.NewControllerConfig(
		useDelete,
		shortDelete,
		shortDelete,
		[]string{
			"del", "d",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)

	return c
}

type DeleteController struct {
	store  *DeleteStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*DeleteStore] = (*DeleteController)(nil)

func NewDeleteController(config *fctl.ControllerConfig) *DeleteController {
	return &DeleteController{
		store:  NewDeleteStore(),
		config: config,
	}
}

func (c *DeleteController) GetStore() *DeleteStore {
	return c.store
}

func (c *DeleteController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *DeleteController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()
	out := c.config.GetOut()
	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config")
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
		return nil, errors.Wrap(err, "creating stack client")
	}

	response, err := client.Orchestration.DeleteWorkflow(
		ctx,
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

func (c *DeleteController) Render() error {
	pterm.Success.WithShowLineNumber().Printfln("Workflow %s Deleted!", c.store.WorkflowId)
	return nil
}

func NewDeleteCommand() *cobra.Command {
	config := NewDeleteConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*DeleteStore](NewDeleteController(config)),
	)
}
