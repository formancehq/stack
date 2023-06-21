package workflows

import (
	"fmt"
	"strings"

	"github.com/formancehq/fctl/cmd/orchestration/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type OrchestrationWorkflowsRunStore struct {
	WorkflowInstance shared.WorkflowInstance `json:"workflow_instance"`
}
type OrchestrationWorkflowsRunController struct {
	store        *OrchestrationWorkflowsRunStore
	variableFlag string
	waitFlag     string
	client       *formance.Formance
	wait         bool
}

var _ fctl.Controller[*OrchestrationWorkflowsRunStore] = (*OrchestrationWorkflowsRunController)(nil)

func NewDefaultOrchestrationWorkflowsRunStore() *OrchestrationWorkflowsRunStore {
	return &OrchestrationWorkflowsRunStore{}
}

func NewOrchestrationWorkflowsRunController() *OrchestrationWorkflowsRunController {
	return &OrchestrationWorkflowsRunController{
		store:        NewDefaultOrchestrationWorkflowsRunStore(),
		variableFlag: "variable",
		waitFlag:     "wait",
		wait:         false,
	}
}

func NewRunCommand() *cobra.Command {
	c := NewOrchestrationWorkflowsRunController()
	return fctl.NewCommand("run <id>",
		fctl.WithShortDescription("Run a workflow"),
		fctl.WithAliases("r"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithBoolFlag(c.waitFlag, false, "Wait end of the run"),
		fctl.WithStringSliceFlag(c.variableFlag, []string{}, "Variable to pass to the workflow"),
		fctl.WithController[*OrchestrationWorkflowsRunStore](c),
	)
}

func (c *OrchestrationWorkflowsRunController) GetStore() *OrchestrationWorkflowsRunStore {
	return c.store
}

func (c *OrchestrationWorkflowsRunController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	soc, err := fctl.GetStackOrganizationConfig(cmd)
	if err != nil {
		return nil, err
	}
	client, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	wait := fctl.GetBool(cmd, c.waitFlag)
	variables := make(map[string]string)
	for _, variable := range fctl.GetStringSlice(cmd, c.variableFlag) {
		parts := strings.SplitN(variable, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("malformed flag: " + variable)
		}
		variables[parts[0]] = parts[1]
	}

	response, err := client.Orchestration.
		RunWorkflow(cmd.Context(), operations.RunWorkflowRequest{
			RequestBody: variables,
			Wait:        &wait,
			WorkflowID:  args[0],
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

	c.wait = wait
	c.store.WorkflowInstance = response.RunWorkflowResponse.Data
	c.client = client
	return c, nil
}

func (c *OrchestrationWorkflowsRunController) Render(cmd *cobra.Command, args []string) error {

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Workflow instance created with ID: %s", c.store.WorkflowInstance.ID)
	if c.wait {
		w, err := c.client.Orchestration.GetWorkflow(cmd.Context(), operations.GetWorkflowRequest{
			FlowID: args[0],
		})
		if err != nil {
			panic(err)
		}

		return internal.PrintWorkflowInstance(cmd.OutOrStdout(), w.GetWorkflowResponse.Data, c.store.WorkflowInstance)
	}
	return nil
}
