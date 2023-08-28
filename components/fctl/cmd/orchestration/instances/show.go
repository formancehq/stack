package instances

import (
	"flag"
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

const (
	useShow   = "show <instance-id>"
	shortShow = "Show a specific workflow instance"
)

type ShowStore struct {
	WorkflowInstance shared.WorkflowInstance `json:"workflowInstance"`
	Workflow         shared.Workflow         `json:"workflow"`
}

func NewDefaultShowStore() *ShowStore {
	return &ShowStore{}
}
func NewShowConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useShow, flag.ExitOnError)

	c := fctl.NewControllerConfig(
		useShow,
		shortShow,
		shortShow,
		[]string{
			"sh",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)

	return c
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

type ShowController struct {
	store  *ShowStore
	config *fctl.ControllerConfig
}

func NewShowController(config *fctl.ControllerConfig) *ShowController {
	return &ShowController{
		store:  NewDefaultShowStore(),
		config: config,
	}
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *ShowController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()

	soc, err := fctl.GetStackOrganizationConfig(flags, ctx, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(flags, ctx, soc.Config, soc.Stack, c.config.GetOut())
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	res, err := client.Orchestration.GetInstance(ctx, operations.GetInstanceRequest{
		InstanceID: args[0],
	})
	if err != nil {
		return nil, errors.Wrap(err, "reading instance")
	}

	if res.Error != nil {
		return nil, fmt.Errorf("%s: %s", res.Error.ErrorCode, res.Error.ErrorMessage)
	}

	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	c.store.WorkflowInstance = res.GetWorkflowInstanceResponse.Data
	response, err := client.Orchestration.GetWorkflow(ctx, operations.GetWorkflowRequest{
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

func (c *ShowController) Render() error {

	out := c.config.GetOut()

	// Print the instance information
	fctl.Section.WithWriter(out).Println("Information")
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
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	if err := internal.PrintWorkflowInstance(out, c.store.Workflow, c.store.WorkflowInstance); err != nil {
		return err
	}

	return nil
}

func NewShowCommand() *cobra.Command {

	config := NewShowConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*ShowStore](NewShowController(config)),
	)
}
