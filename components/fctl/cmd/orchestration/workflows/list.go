package workflows

import (
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type WorkflowsListStore struct {
	Workflows []shared.Workflow `json:"workflows"`
}
type WorkflowsListController struct {
	store *WorkflowsListStore
}

var _ fctl.Controller[*WorkflowsListStore] = (*WorkflowsListController)(nil)

func NewDefaultWorkflowsListStore() *WorkflowsListStore {
	return &WorkflowsListStore{}
}

func NewWorkflowsListController() *WorkflowsListController {
	return &WorkflowsListController{
		store: NewDefaultWorkflowsListStore(),
	}
}

func NewListCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithShortDescription("List all workflows"),
		fctl.WithAliases("ls", "l"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*WorkflowsListStore](NewWorkflowsListController()),
	)
}

func (c *WorkflowsListController) GetStore() *WorkflowsListStore {
	return c.store
}

func (c *WorkflowsListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	soc, err := fctl.GetStackOrganizationConfig(cmd)
	if err != nil {
		return nil, err
	}
	client, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	response, err := client.Orchestration.ListWorkflows(cmd.Context())
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("%s: %s", response.Error.ErrorCode, response.Error.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Workflows = response.ListWorkflowsResponse.Data

	return c, nil
}

func (c *WorkflowsListController) Render(cmd *cobra.Command, args []string) error {
	if len(c.store.Workflows) == 0 {
		fctl.Println("No workflows found.")
		return nil
	}

	if err := pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(cmd.OutOrStdout()).
		WithData(
			fctl.Prepend(
				fctl.Map(c.store.Workflows,
					func(src shared.Workflow) []string {
						return []string{
							src.ID,
							func() string {
								if src.Config.Name != nil {
									return *src.Config.Name
								}
								return ""
							}(),
							src.CreatedAt.Format(time.RFC3339),
							src.UpdatedAt.Format(time.RFC3339),
						}
					}),
				[]string{"ID", "Name", "Created at", "Updated at"},
			),
		).Render(); err != nil {
		return errors.Wrap(err, "rendering table")
	}

	return nil
}
