package triggers

import (
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type TriggersShowStore struct {
	Trigger shared.Trigger `json:"trigger"`
}
type TriggersShowController struct {
	store *TriggersShowStore
}

var _ fctl.Controller[*TriggersShowStore] = (*TriggersShowController)(nil)

func NewDefaultTriggersShowStore() *TriggersShowStore {
	return &TriggersShowStore{}
}

func NewTriggersShowController() *TriggersShowController {
	return &TriggersShowController{
		store: NewDefaultTriggersShowStore(),
	}
}

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("show <trigger-id>",
		fctl.WithShortDescription("Show a specific workflow trigger"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*TriggersShowStore](NewTriggersShowController()),
	)
}

func (c *TriggersShowController) GetStore() *TriggersShowStore {
	return c.store
}

func (c *TriggersShowController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	soc, err := fctl.GetStackOrganizationConfig(cmd)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	res, err := client.Orchestration.ReadTrigger(cmd.Context(), operations.ReadTriggerRequest{
		TriggerID: args[0],
	})
	if err != nil {
		return nil, errors.Wrap(err, "reading trigger")
	}

	if res.Error != nil {
		return nil, fmt.Errorf("%s: %s", res.Error.ErrorCode, res.Error.ErrorMessage)
	}

	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	c.store.Trigger = res.ReadTriggerResponse.Data

	return c, nil
}

func (c *TriggersShowController) Render(cmd *cobra.Command, args []string) error {
	// Print the trigger information
	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.Trigger.ID})
	tableData = append(tableData, []string{pterm.LightCyan("Created at"), c.store.Trigger.CreatedAt.Format(time.RFC3339)})
	tableData = append(tableData, []string{pterm.LightCyan("Workflow ID"), c.store.Trigger.WorkflowID})
	tableData = append(tableData, []string{pterm.LightCyan("Event"), c.store.Trigger.Event})
	tableData = append(tableData, []string{pterm.LightCyan("Filter"), func() string {
		if c.store.Trigger.Filter == nil {
			return ""
		}
		return *c.store.Trigger.Filter
	}()})

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	return nil
}
