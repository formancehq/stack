package workflows

import (
	"fmt"
	"strings"

	"github.com/formancehq/fctl/cmd/orchestration/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type WorkflowsRunStore struct {
	WorkflowInstance shared.WorkflowInstance `json:"workflowInstance"`
}
type WorkflowsRunController struct {
	store        *WorkflowsRunStore
	variableFlag string
	waitFlag     string
	wait         bool
}

var _ fctl.Controller[*WorkflowsRunStore] = (*WorkflowsRunController)(nil)

func NewDefaultWorkflowsRunStore() *WorkflowsRunStore {
	return &WorkflowsRunStore{}
}

func NewWorkflowsRunController() *WorkflowsRunController {
	return &WorkflowsRunController{
		store:        NewDefaultWorkflowsRunStore(),
		variableFlag: "variable",
		waitFlag:     "wait",
		wait:         false,
	}
}

func NewRunCommand() *cobra.Command {
	c := NewWorkflowsRunController()
	return fctl.NewCommand("run <id>",
		fctl.WithShortDescription("Run a workflow"),
		fctl.WithAliases("r"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithBoolFlag(c.waitFlag, false, "Wait end of the run"),
		fctl.WithStringSliceFlag(c.variableFlag, []string{}, "Variable to pass to the workflow"),
		fctl.WithController[*WorkflowsRunStore](c),
	)
}

func (c *WorkflowsRunController) GetStore() *WorkflowsRunStore {
	return c.store
}

func (c *WorkflowsRunController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := fctl.GetStackStore(cmd.Context())

	wait := fctl.GetBool(cmd, c.waitFlag)
	variables := make(map[string]string)
	for _, variable := range fctl.GetStringSlice(cmd, c.variableFlag) {
		parts := strings.SplitN(variable, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("malformed flag: " + variable)
		}
		variables[parts[0]] = parts[1]
	}

	response, err := store.Client().Orchestration.
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
	return c, nil
}

func (c *WorkflowsRunController) Render(cmd *cobra.Command, args []string) error {
	store := fctl.GetStackStore(cmd.Context())
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Workflow instance created with ID: %s", c.store.WorkflowInstance.ID)
	if c.wait {
		w, err := store.Client().Orchestration.GetWorkflow(cmd.Context(), operations.GetWorkflowRequest{
			FlowID: args[0],
		})
		if err != nil {
			panic(err)
		}

		if w.Error != nil {
			panic(fmt.Sprintf("%s: %s", w.Error.ErrorCode, w.Error.ErrorMessage))
		}

		if w.StatusCode >= 300 {
			panic(fmt.Sprintf("unexpected status code: %d", w.StatusCode))
		}

		return internal.PrintWorkflowInstance(cmd.OutOrStdout(), w.GetWorkflowResponse.Data, c.store.WorkflowInstance)
	}
	return nil
}
