package modules

import (
	"github.com/formancehq/fctl/cmd/stack/store"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/stack/libs/go-libs/time"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStore struct {
	*membershipclient.ListModulesResponse
}
type ListController struct {
	store *ListStore
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewDefaultListStore() *ListStore {
	return &ListStore{}
}

func NewListController() *ListController {
	return &ListController{
		store: NewDefaultListStore(),
	}
}

func NewListCommand() *cobra.Command {
	return fctl.NewMembershipCommand("list --stack=<stack-id>",
		fctl.WithConfirmFlag(),
		fctl.WithShortDescription("List modules in a stack"),
		fctl.WithAliases("ls"),
		fctl.WithController(NewListController()),
	)
}
func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := store.GetStore(cmd.Context())

	modules, _, err := store.Client().ListModules(cmd.Context(), store.OrganizationId(), fctl.GetString(cmd, stackFlag)).Execute()
	if err != nil {
		return nil, err
	}

	c.store.ListModulesResponse = modules

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	header := []string{"Name", "State", "Cluster status", "Last state update", "Last cluster state update"}

	tableData := fctl.Map(c.store.ListModulesResponse.Data, func(module membershipclient.Module) []string {
		return []string{
			module.Name,
			module.State,
			module.Status,
			time.Time{Time: module.LastStateUpdate}.String(),
			time.Time{Time: module.LastStatusUpdate}.String(),
		}
	})

	tableData = fctl.Prepend(tableData, header)

	return pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithHasHeader().
		WithData(tableData).
		Render()
}
