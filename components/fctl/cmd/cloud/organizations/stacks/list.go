package stacks

import (
	"strings"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStackAccessRolesStore struct {
	list []membershipclient.StackUserAccess
}
type ListStackAccessRolesController struct {
	store *ListStackAccessRolesStore
}

var _ fctl.Controller[*ListStackAccessRolesStore] = (*ListStackAccessRolesController)(nil)

func NewDefaultListStackAccessRolesStore() *ListStackAccessRolesStore {
	return &ListStackAccessRolesStore{
		list: []membershipclient.StackUserAccess{},
	}
}

func NewListStackAccessRolesController() *ListStackAccessRolesController {
	return &ListStackAccessRolesController{
		store: NewDefaultListStackAccessRolesStore(),
	}
}

func NewListStackAccessRolesCommand() *cobra.Command {
	return fctl.NewCommand("list-access-roles <stack-id>",
		fctl.WithAliases("usar"),
		fctl.WithShortDescription("List Stack Access Roles within an organization by stacks"),
		fctl.WithArgs(cobra.MinimumNArgs(1)),
		fctl.WithController[*ListStackAccessRolesStore](NewListStackAccessRolesController()),
	)
}

func (c *ListStackAccessRolesController) GetStore() *ListStackAccessRolesStore {
	return c.store
}

func (c *ListStackAccessRolesController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

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

	c.store.list = append(c.store.list, ListStackUsersAccesses...)

	return c, nil
}

func (c *ListStackAccessRolesController) Render(cmd *cobra.Command, args []string) error {
	stackUserAccessMap := fctl.Map(c.store.list, func(o membershipclient.StackUserAccess) []string {
		return []string{
			o.StackId,
			o.UserId,
			strings.Join(o.Roles, "|"),
		}
	})

	tableData := fctl.Prepend(stackUserAccessMap, []string{"StackId", "UserId", "Roles"})

	return pterm.DefaultTable.WithHasHeader().WithWriter(cmd.OutOrStdout()).WithData(tableData).Render()

}
