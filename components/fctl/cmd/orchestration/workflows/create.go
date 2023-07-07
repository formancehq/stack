package workflows

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type WorkflowsCreateStore struct {
	WorkflowId string `json:"workflowId"`
}
type WorkflowsCreateController struct {
	store *WorkflowsCreateStore
}

var _ fctl.Controller[*WorkflowsCreateStore] = (*WorkflowsCreateController)(nil)

func NewDefaultWorkflowsCreateStore() *WorkflowsCreateStore {
	return &WorkflowsCreateStore{}
}

func NewWorkflowsCreateController() *WorkflowsCreateController {
	return &WorkflowsCreateController{
		store: NewDefaultWorkflowsCreateStore(),
	}
}

func NewCreateCommand() *cobra.Command {
	return fctl.NewCommand("create <file>|-",
		fctl.WithShortDescription("Create a workflow"),
		fctl.WithAliases("cr", "c"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*WorkflowsCreateStore](NewWorkflowsCreateController()),
	)
}

func (c *WorkflowsCreateController) GetStore() *WorkflowsCreateStore {
	return c.store
}

func (c *WorkflowsCreateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	soc, err := fctl.GetStackOrganizationConfig(cmd)
	if err != nil {
		return nil, err
	}
	client, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	script, err := fctl.ReadFile(cmd, soc.Stack, args[0])
	if err != nil {
		return nil, err
	}

	config := shared.WorkflowConfig{}
	if err := yaml.Unmarshal([]byte(script), &config); err != nil {
		return nil, err
	}

	//nolint:gosimple
	response, err := client.Orchestration.
		CreateWorkflow(cmd.Context(), shared.CreateWorkflowRequest{
			Name:   config.Name,
			Stages: config.Stages,
		})
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("%s: %s", response.Error.ErrorCode, response.Error.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.WorkflowId = response.CreateWorkflowResponse.Data.ID

	return c, nil
}

func (c *WorkflowsCreateController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Workflow created with ID: %s", c.store.WorkflowId)

	return nil
}
