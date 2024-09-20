package users

import (
	"github.com/formancehq/fctl/cmd/stack/store"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStore struct {
	list []membershipclient.StackUserAccess
}
type ListController struct {
	store *ListStore
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewDefaultListStore() *ListStore {
	return &ListStore{
		list: []membershipclient.StackUserAccess{},
	}
}

func NewListController() *ListController {
	return &ListController{
		store: NewDefaultListStore(),
	}
}

func NewListCommand() *cobra.Command {
	return fctl.NewCommand("list <stack-id>",
		fctl.WithAliases("usar"),
		fctl.WithShortDescription("List Stack Access Role within an organization by stacks"),
		fctl.WithArgs(cobra.MinimumNArgs(1)),
		fctl.WithController[*ListStore](NewListController()),
	)
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := store.GetStore(cmd.Context())

	ListStackUsersAccesses, response, err := store.Client().ListStackUsersAccesses(cmd.Context(), store.OrganizationId(), args[0]).Execute()
	if err != nil {
		return nil, err
	}

	if response.StatusCode > 300 {
		return nil, err
	}

	c.store.list = append(c.store.list, ListStackUsersAccesses.Data...)

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	stackUserAccessMap := fctl.Map(c.store.list, func(o membershipclient.StackUserAccess) []string {
		return []string{
			o.StackId,
			o.UserId,
			string(o.Role),
		}
	})

	tableData := fctl.Prepend(stackUserAccessMap, []string{"Stack Id", "User Id", "Role"})

	return pterm.DefaultTable.WithHasHeader().WithWriter(cmd.OutOrStdout()).WithData(tableData).Render()

}
