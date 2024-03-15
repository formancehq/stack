package triggers

import (
	"fmt"
	"strings"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type TriggersCreateStore struct {
	Trigger shared.Trigger `json:"trigger"`
}
type TriggersCreateController struct {
	store      *TriggersCreateStore
	filterFlag string
	varsFlag   string
}

var _ fctl.Controller[*TriggersCreateStore] = (*TriggersCreateController)(nil)

func NewDefaultTriggersCreateStore() *TriggersCreateStore {
	return &TriggersCreateStore{}
}

func NewTriggersCreateController() *TriggersCreateController {
	return &TriggersCreateController{
		store:      NewDefaultTriggersCreateStore(),
		filterFlag: "filter",
		varsFlag:   "vars",
	}
}

func NewCreateCommand() *cobra.Command {
	ctrl := NewTriggersCreateController()
	return fctl.NewCommand("create <event> <workflow-id>",
		fctl.WithShortDescription("Create a trigger"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*TriggersCreateStore](ctrl),
		fctl.WithStringFlag(ctrl.filterFlag, "", "Filter events"),
		fctl.WithStringSliceFlag(ctrl.varsFlag, []string{}, "Variables to pass to the workflow"),
	)
}

func (c *TriggersCreateController) GetStore() *TriggersCreateStore {
	return c.store
}

func (c *TriggersCreateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	var (
		event    = args[0]
		filter   = fctl.GetString(cmd, c.filterFlag)
		vars     = fctl.GetStringSlice(cmd, c.varsFlag)
		workflow = args[1]
	)

	data := &shared.TriggerData{
		Event:      event,
		WorkflowID: workflow,
		Vars:       map[string]interface{}{},
	}
	if filter != "" {
		data.Filter = pointer.For(filter)
	}
	if len(vars) > 0 {
		for _, v := range vars {
			parts := strings.SplitN(v, "=", 2)
			if len(parts) != 2 {
				return nil, errors.New("invalid 'vars' flag")
			}
			data.Vars[parts[0]] = parts[1]
		}
	}

	res, err := store.Client().Orchestration.CreateTrigger(cmd.Context(), data)
	if err != nil {
		return nil, errors.Wrap(err, "reading trigger")
	}

	if res.Error != nil {
		return nil, fmt.Errorf("%s: %s", res.Error.ErrorCode, res.Error.ErrorMessage)
	}

	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	c.store.Trigger = res.CreateTriggerResponse.Data

	return c, nil
}

func (c *TriggersCreateController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Trigger created with ID: %s", c.store.Trigger.ID)

	return nil
}
