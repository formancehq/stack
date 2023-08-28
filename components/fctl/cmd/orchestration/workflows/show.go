package workflows

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
	"gopkg.in/yaml.v3"
)

const (
	useShow   = "show <id>"
	shortShow = "Show a workflow"
)

type ShowStore struct {
	Workflow shared.Workflow `json:"workflow"`
}

func NewShowStore() *ShowStore {
	return &ShowStore{}
}
func NewShowConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useShow, flag.ExitOnError)

	c := fctl.NewControllerConfig(
		useShow,
		shortShow,
		shortShow,
		[]string{
			"s",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)

	return c
}

type ShowController struct {
	store  *ShowStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

func NewShowController(config *fctl.ControllerConfig) *ShowController {
	return &ShowController{
		store:  NewShowStore(),
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

	response, err := client.Orchestration.
		GetWorkflow(ctx, operations.GetWorkflowRequest{
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

func (c *ShowController) Render() error {

	out := c.config.GetOut()

	fctl.Section.WithWriter(out).Println("Information")
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
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	fmt.Fprintln(out)

	fctl.Section.WithWriter(out).Println("Configuration")
	configAsBytes, err := yaml.Marshal(c.store.Workflow.Config)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(out, string(configAsBytes))

	return nil
}

func NewShowCommand() *cobra.Command {
	config := NewShowConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*ShowStore](NewShowController(config)),
	)
}
