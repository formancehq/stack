package instances

import (
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type InstancesListStore struct {
	WorkflowInstance []shared.WorkflowInstance `json:"workflow_instance"`
}
type InstancesListController struct {
	store        *InstancesListStore
	workflowFlag string
	runningFlag  string
}

var _ fctl.Controller[*InstancesListStore] = (*InstancesListController)(nil)

func NewDefaultInstancesListStore() *InstancesListStore {
	return &InstancesListStore{}
}

func NewInstancesListController() *InstancesListController {
	return &InstancesListController{
		store:        NewDefaultInstancesListStore(),
		workflowFlag: "workflow",
		runningFlag:  "running",
	}
}

func NewListCommand() *cobra.Command {
	c := NewInstancesListController()
	return fctl.NewCommand("list",
		fctl.WithShortDescription("List all workflows instances"),
		fctl.WithAliases("ls", "l"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithStringFlag(c.workflowFlag, "", "Filter on workflow id"),
		fctl.WithBoolFlag(c.runningFlag, false, "Filter on running instances"),
		fctl.WithController[*InstancesListStore](c),
	)
}

func (c *InstancesListController) GetStore() *InstancesListStore {
	return c.store
}

func (c *InstancesListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	soc, err := fctl.GetStackOrganizationConfig(cmd)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	response, err := client.Orchestration.ListInstances(cmd.Context(), operations.ListInstancesRequest{
		Running:    fctl.Ptr(fctl.GetBool(cmd, c.runningFlag)),
		WorkflowID: fctl.Ptr(fctl.GetString(cmd, c.workflowFlag)),
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

	c.store.WorkflowInstance = response.ListRunsResponse.Data

	return c, nil
}

func (c *InstancesListController) Render(cmd *cobra.Command, args []string) error {

	if len(c.store.WorkflowInstance) == 0 {
		fctl.Println("No workflows found.")
		return nil
	}
	if err := pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(cmd.OutOrStdout()).
		WithData(
			fctl.Prepend(
				fctl.Map(c.store.WorkflowInstance,
					func(src shared.WorkflowInstance) []string {
						return []string{
							src.ID,
							src.WorkflowID,
							src.CreatedAt.Format(time.RFC3339),
							src.UpdatedAt.Format(time.RFC3339),
							func() string {
								if src.Terminated {
									return src.TerminatedAt.Format(time.RFC3339)
								}
								return ""
							}(),
						}
					}),
				[]string{"ID", "Workflow ID", "Created at", "Updated at", "Terminated at"},
			),
		).Render(); err != nil {
		return errors.Wrap(err, "rendering table")
	}

	return nil
}
