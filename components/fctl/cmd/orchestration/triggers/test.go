package triggers

import (
	"encoding/json"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type TriggersTestStore struct {
	Trigger shared.V2TriggerTest `json:"trigger"`
}
type TriggersTestController struct {
	store *TriggersTestStore
}

var _ fctl.Controller[*TriggersTestStore] = (*TriggersTestController)(nil)

func NewDefaultTriggersTestStore() *TriggersTestStore {
	return &TriggersTestStore{}
}

func NewTriggersTestController() *TriggersTestController {
	return &TriggersTestController{
		store: NewDefaultTriggersTestStore(),
	}
}

func NewTestCommand() *cobra.Command {
	return fctl.NewCommand("test <trigger-id> <event>",
		fctl.WithShortDescription("Test a specific workflow trigger"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*TriggersTestStore](NewTriggersTestController()),
	)
}

func (c *TriggersTestController) GetStore() *TriggersTestStore {
	return c.store
}

func (c *TriggersTestController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	data := make(map[string]any)
	if err := json.Unmarshal([]byte(args[1]), &data); err != nil {
		return nil, err
	}

	res, err := store.Client().Orchestration.V2.TestTrigger(cmd.Context(), operations.TestTriggerRequest{
		TriggerID:   args[0],
		RequestBody: data,
	})
	if err != nil {
		return nil, errors.Wrap(err, "testing trigger")
	}

	c.store.Trigger = res.V2TestTriggerResponse.Data

	return c, nil
}

func (c *TriggersTestController) Render(cmd *cobra.Command, args []string) error {
	// Print the trigger information
	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
	tableData := pterm.TableData{}
	if c.store.Trigger.Filter != nil {
		if c.store.Trigger.Filter.Match != nil {
			tableData = append(tableData, []string{pterm.LightCyan("Filter match"), fctl.BoolToString(*c.store.Trigger.Filter.Match)})
		}
		if c.store.Trigger.Filter.Error != nil {
			tableData = append(tableData, []string{pterm.LightCyan("Filter error"), *c.store.Trigger.Filter.Error})
		}
	}

	if err := pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	if c.store.Trigger.Variables != nil && len(c.store.Trigger.Variables) > 0 {
		fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Variables")
		tableData := make([][]string, 0)
		for key, variable := range c.store.Trigger.Variables {
			tableData = append(tableData, []string{
				key,
				func() string {
					if variable.Value != nil {
						return *variable.Value
					}
					return ""
				}(),
				func() string {
					if variable.Error != nil {
						return *variable.Error
					}
					return ""
				}(),
			})
		}
		tableData = fctl.Prepend(tableData, []string{"Key", "Value", "Error"})

		if err := pterm.DefaultTable.
			WithHasHeader().
			WithWriter(cmd.OutOrStdout()).
			WithData(tableData).
			Render(); err != nil {
			return err
		}
	}

	return nil
}
