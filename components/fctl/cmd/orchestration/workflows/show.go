package workflows

import (
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type WorkflowsShowStore struct {
	Workflow shared.Workflow `json:"workflow"`
}
type WorkflowsShowController struct {
	store *WorkflowsShowStore
}

var _ fctl.Controller[*WorkflowsShowStore] = (*WorkflowsShowController)(nil)

func NewDefaultWorkflowsShowStore() *WorkflowsShowStore {
	return &WorkflowsShowStore{}
}

func NewWorkflowsShowController() *WorkflowsShowController {
	return &WorkflowsShowController{
		store: NewDefaultWorkflowsShowStore(),
	}
}

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("show <id>",
		fctl.WithShortDescription("Show a workflow"),
		fctl.WithAliases("s"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*WorkflowsShowStore](NewWorkflowsShowController()),
	)
}

func (c *WorkflowsShowController) GetStore() *WorkflowsShowStore {
	return c.store
}

func (c *WorkflowsShowController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	soc, err := fctl.GetStackOrganizationConfig(cmd)
	if err != nil {
		return nil, err
	}
	client, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	response, err := client.Orchestration.
		GetWorkflow(cmd.Context(), operations.GetWorkflowRequest{
			FlowID: args[0],
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

func (c *WorkflowsShowController) Render(cmd *cobra.Command, args []string) error {
	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.Workflow.ID})
	tableData = append(tableData, []string{pterm.LightCyan("Name"), func() string {
		if c.store.Workflow.Config.Name != nil {
			return *c.store.Workflow.Config.Name
		}
		return ""
	}()})
	tableData = append(tableData, []string{pterm.LightCyan("Created at"), c.store.Workflow.CreatedAt.Format(time.RFC3339)})
	tableData = append(tableData, []string{pterm.LightCyan("Updated at"), c.store.Workflow.UpdatedAt.Format(time.RFC3339)})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout())

	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Configuration")
	configAsBytes, err := yaml.Marshal(c.store.Workflow.Config)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(cmd.OutOrStdout(), string(configAsBytes))

	return nil
}
