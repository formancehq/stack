package instances

import (
	"fmt"
	"time"

	"github.com/formancehq/fctl/cmd/orchestration/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type InstancesShowStore struct {
	WorkflowInstance shared.WorkflowInstance `json:"workflow_instance"`
	Workflow         shared.Workflow         `json:"workflow"`
}
type InstancesShowController struct {
	store *InstancesShowStore
}

var _ fctl.Controller[*InstancesShowStore] = (*InstancesShowController)(nil)

func NewDefaultInstancesShowStore() *InstancesShowStore {
	return &InstancesShowStore{}
}

func NewInstancesShowController() *InstancesShowController {
	return &InstancesShowController{
		store: NewDefaultInstancesShowStore(),
	}
}

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("show <instance-id>",
		fctl.WithShortDescription("Show a specific workflow instance"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*InstancesShowStore](NewInstancesShowController()),
	)
}

func (c *InstancesShowController) GetStore() *InstancesShowStore {
	return c.store
}

func (c *InstancesShowController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	soc, err := fctl.GetStackOrganizationConfig(cmd)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	res, err := client.Orchestration.GetInstance(cmd.Context(), operations.GetInstanceRequest{
		InstanceID: args[0],
	})

	if err != nil {
		return nil, errors.Wrap(err, "reading instance")
	}

	c.store.WorkflowInstance = res.GetWorkflowInstanceResponse.Data
	response, err := client.Orchestration.GetWorkflow(cmd.Context(), operations.GetWorkflowRequest{
		FlowID: res.GetWorkflowInstanceResponse.Data.WorkflowID,
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

	c.store.Workflow = response.GetWorkflowResponse.Data

	return c, nil
}

func (c *InstancesShowController) Render(cmd *cobra.Command, args []string) error {
	// Print the instance information
	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.WorkflowInstance.ID})
	tableData = append(tableData, []string{pterm.LightCyan("Created at"), c.store.WorkflowInstance.CreatedAt.Format(time.RFC3339)})
	tableData = append(tableData, []string{pterm.LightCyan("Updated at"), c.store.WorkflowInstance.UpdatedAt.Format(time.RFC3339)})
	if c.store.WorkflowInstance.Terminated {
		tableData = append(tableData, []string{pterm.LightCyan("Terminated at"), c.store.WorkflowInstance.TerminatedAt.Format(time.RFC3339)})
	}
	if c.store.WorkflowInstance.Error != nil && *c.store.WorkflowInstance.Error != "" {
		tableData = append(tableData, []string{pterm.LightCyan("Error"), *c.store.WorkflowInstance.Error})
	}

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	if err := internal.PrintWorkflowInstance(cmd.OutOrStdout(), c.store.Workflow, c.store.WorkflowInstance); err != nil {
		return err
	}

	return nil
}
