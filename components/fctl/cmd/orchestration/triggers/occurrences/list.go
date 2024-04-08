package occurrences

import (
	"time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type OccurrencesListStore struct {
	WorkflowOccurrence []shared.TriggerOccurrence `json:"occurrences"`
}
type OccurrencesListController struct {
	store *OccurrencesListStore
}

var _ fctl.Controller[*OccurrencesListStore] = (*OccurrencesListController)(nil)

func NewDefaultOccurrencesListStore() *OccurrencesListStore {
	return &OccurrencesListStore{}
}

func NewOccurrencesListController() *OccurrencesListController {
	return &OccurrencesListController{
		store: NewDefaultOccurrencesListStore(),
	}
}

func NewListCommand() *cobra.Command {
	c := NewOccurrencesListController()
	return fctl.NewCommand("list",
		fctl.WithShortDescription("List all workflows occurrences"),
		fctl.WithAliases("ls", "l"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*OccurrencesListStore](c),
	)
}

func (c *OccurrencesListController) GetStore() *OccurrencesListStore {
	return c.store
}

func (c *OccurrencesListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	response, err := store.Client().Orchestration.ListTriggersOccurrences(cmd.Context(), operations.ListTriggersOccurrencesRequest{
		TriggerID: args[0],
	})
	if err != nil {
		return nil, err
	}

	c.store.WorkflowOccurrence = response.ListTriggersOccurrencesResponse.Data

	return c, nil
}

func (c *OccurrencesListController) Render(cmd *cobra.Command, args []string) error {

	if len(c.store.WorkflowOccurrence) == 0 {
		fctl.Println("No occurrences found.")
		return nil
	}
	if err := pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(cmd.OutOrStdout()).
		WithData(
			fctl.Prepend(
				fctl.Map(c.store.WorkflowOccurrence,
					func(src shared.TriggerOccurrence) []string {
						return []string{
							func() string {
								if src.WorkflowInstanceID != nil {
									return ""
								}
								return *src.WorkflowInstanceID
							}(),
							src.Date.Format(time.RFC3339),
							fctl.BoolToString(src.WorkflowInstance.Terminated),
							src.WorkflowInstance.TerminatedAt.Format(time.RFC3339),
							func() string {
								if src.Error != nil && *src.Error != "" {
									return *src.Error
								}
								if src.WorkflowInstance != nil && src.WorkflowInstance.Error != nil && *src.WorkflowInstance.Error != "" {
									return *src.WorkflowInstance.Error
								}
								return ""
							}(),
						}
					}),
				[]string{"Workflow instance ID", "Date", "Terminated", "Terminated at", "Error"},
			),
		).Render(); err != nil {
		return errors.Wrap(err, "rendering table")
	}

	return nil
}
