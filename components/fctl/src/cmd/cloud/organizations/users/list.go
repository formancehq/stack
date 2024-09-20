package users

import (
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type User struct {
	ID    string                `json:"id"`
	Email string                `json:"email"`
	Role  membershipclient.Role `json:"role"`
}

type ListStore struct {
	Users []User `json:"users"`
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
	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List users"),
		fctl.WithController[*ListStore](NewListController()),
	)
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetMembershipStore(cmd.Context())

	organizationID, err := fctl.ResolveOrganizationID(cmd, store.Config, store.Client())
	if err != nil {
		return nil, err
	}

	usersResponse, _, err := store.Client().ListUsersOfOrganization(cmd.Context(), organizationID).Execute()
	if err != nil {
		return nil, err
	}

	c.store.Users = fctl.Map(usersResponse.Data, func(i membershipclient.OrganizationUserArrayInner) User {
		return User{
			i.Id,
			i.Email,
			i.Role,
		}
	})

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {

	usersRow := fctl.Map(c.store.Users, func(i User) []string {
		return []string{
			i.ID,
			i.Email,
			string(i.Role),
		}
	})

	tableData := fctl.Prepend(usersRow, []string{"ID", "Email", "Role"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()

}
