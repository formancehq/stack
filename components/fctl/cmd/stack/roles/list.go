package roles

import (
	"strings"

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
		fctl.WithShortDescription("List Stack Access Roles within an organization by stacks"),
		fctl.WithArgs(cobra.MinimumNArgs(1)),
		fctl.WithController[*ListStore](NewListController()),
	)
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	ListStackUsersAccesses, response, err := apiClient.DefaultApi.ListStackUsersAccesses(cmd.Context(), organizationID, args[0]).Execute()
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
			func() string {
				roles := []string{}

				for _, role := range o.Roles {
					if role == "ADMIN" {
						roles = append(roles, pterm.LightRed(role))
					} else {
						roles = append(roles, pterm.LightGreen(role))
					}
				}

				return strings.Join(roles, " | ")
			}(),
		}
	})

	tableData := fctl.Prepend(stackUserAccessMap, []string{"Stack Id", "User Id", "Roles"})

	return pterm.DefaultTable.WithHasHeader().WithWriter(cmd.OutOrStdout()).WithData(tableData).Render()

}
