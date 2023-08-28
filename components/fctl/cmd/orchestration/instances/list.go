package instances

import (
	"flag"
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	workflowFlag    = "workflow"
	runningFlag     = "running"
	useList         = "list"
	descriptionList = "List all workflow instances"
	shortList       = "List all workflow instances"
)

type WorkflowInstance struct {
	InstanceID   string `json:"instanceId"`
	WorkflowID   string `json:"workflowId"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
	TerminatedAt string `json:"terminatedAt"`
}

type ListStore struct {
	WorkflowInstance []WorkflowInstance `json:"workflowInstances"`
}

func NewDefaultListStore() *ListStore {
	return &ListStore{}
}
func NewListConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	flags.String(workflowFlag, "", "Filter on workflow id")
	flags.Bool(runningFlag, false, "Filter on running instances")
	c := fctl.NewControllerConfig(
		useList,
		descriptionList,
		shortList,
		[]string{
			"list",
			"ls",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)

	return c
}

type ListController struct {
	store  *ListStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewListController(config *fctl.ControllerConfig) *ListController {
	return &ListController{
		store:  NewDefaultListStore(),
		config: config,
	}
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *ListController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	soc, err := fctl.GetStackOrganizationConfig(flags, ctx, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(flags, ctx, soc.Config, soc.Stack, c.config.GetOut())
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	response, err := client.Orchestration.ListInstances(ctx, operations.ListInstancesRequest{
		Running:    fctl.Ptr(fctl.GetBool(flags, runningFlag)),
		WorkflowID: fctl.Ptr(fctl.GetString(flags, workflowFlag)),
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

	c.store.WorkflowInstance = fctl.Map(response.ListRunsResponse.Data, func(src shared.WorkflowInstance) WorkflowInstance {
		return WorkflowInstance{
			InstanceID: src.ID,
			WorkflowID: src.WorkflowID,
			CreatedAt:  src.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  src.UpdatedAt.Format(time.RFC3339),
			TerminatedAt: func() string {
				if src.TerminatedAt == nil {
					return ""
				}
				return src.TerminatedAt.Format(time.RFC3339)
			}(),
		}
	})

	return c, nil
}

func (c *ListController) Render() error {

	if len(c.store.WorkflowInstance) == 0 {
		fmt.Fprintln(c.config.GetOut(), "No workflows found.")
		return nil
	}
	if err := pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(c.config.GetOut()).
		WithData(
			fctl.Prepend(
				fctl.Map(c.store.WorkflowInstance,
					func(src WorkflowInstance) []string {
						return []string{
							src.InstanceID,
							src.WorkflowID,
							src.CreatedAt,
							src.UpdatedAt,
							src.TerminatedAt,
						}
					}),
				[]string{"ID", "Workflow ID", "Created at", "Updated at", "Terminated at"},
			),
		).Render(); err != nil {
		return errors.Wrap(err, "rendering table")
	}

	return nil
}

func NewListCommand() *cobra.Command {
	c := NewListConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController(c)),
	)
}
