package triggers

import (
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type TriggersListStore struct {
	WorkflowTrigger []shared.Trigger `json:"workflowTriggers"`
}
type TriggersListController struct {
	store *TriggersListStore
}

var _ fctl.Controller[*TriggersListStore] = (*TriggersListController)(nil)

func NewDefaultTriggersListStore() *TriggersListStore {
	return &TriggersListStore{}
}

func NewTriggersListController() *TriggersListController {
	return &TriggersListController{
		store: NewDefaultTriggersListStore(),
	}
}

func NewListCommand() *cobra.Command {
	c := NewTriggersListController()
	return fctl.NewCommand("list",
		fctl.WithShortDescription("List all workflows triggers"),
		fctl.WithAliases("ls", "l"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*TriggersListStore](c),
	)
}

func (c *TriggersListController) GetStore() *TriggersListStore {
	return c.store
}

func (c *TriggersListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	response, err := store.Client().Orchestration.ListTriggers(cmd.Context())
	if err != nil {
		return nil, err
	}

	c.store.WorkflowTrigger = response.ListTriggersResponse.Data

	return c, nil
}

func (c *TriggersListController) Render(cmd *cobra.Command, args []string) error {

	if len(c.store.WorkflowTrigger) == 0 {
		fctl.Println("No triggers found.")
		return nil
	}
	if err := pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(cmd.OutOrStdout()).
		WithData(
			fctl.Prepend(
				fctl.Map(c.store.WorkflowTrigger,
					func(src shared.Trigger) []string {
						return []string{
							src.ID,
							src.WorkflowID,
							src.CreatedAt.Format(time.RFC3339),
							src.Event,
							func() string {
								if src.Filter == nil {
									return ""
								}
								return *src.Filter
							}(),
						}
					}),
				[]string{"ID", "Workflow ID", "Created at", "Event", "Filter"},
			),
		).Render(); err != nil {
		return errors.Wrap(err, "rendering table")
	}

	return nil
}
