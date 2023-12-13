package stacks

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type UpsertStackAccessRolesStore struct {
	OrganizationID string `json:"organizationId"`
	UserID         string `json:"userId"`
	StackID        string `json:"stackId"`
}
type UpsertStackAccessRolesController struct {
	store *UpsertStackAccessRolesStore
}

var _ fctl.Controller[*UpsertStackAccessRolesStore] = (*UpsertStackAccessRolesController)(nil)

func NewDefaultUpsertStackAccessRolesStore() *UpsertStackAccessRolesStore {
	return &UpsertStackAccessRolesStore{}
}

func NewUpsertStackAccessRolesController() *UpsertStackAccessRolesController {
	return &UpsertStackAccessRolesController{
		store: NewDefaultUpsertStackAccessRolesStore(),
	}
}

func NewUpsertStackAccessRolesCommand() *cobra.Command {
	return fctl.NewCommand("update-access-roles <user-id> <stack-id> <stack-roles...>",
		fctl.WithAliases("usar"),
		fctl.WithShortDescription("Update Stack Access Roles within an organization (ADMIN, GUEST, \"\")"),
		fctl.WithArgs(cobra.MinimumNArgs(2)),
		fctl.WithController[*UpsertStackAccessRolesStore](NewUpsertStackAccessRolesController()),
	)
}

func (c *UpsertStackAccessRolesController) GetStore() *UpsertStackAccessRolesStore {
	return c.store
}

func (c *UpsertStackAccessRolesController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

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

	_, err = apiClient.DefaultApi.UpsertStackUserAccess(cmd.Context(), organizationID, args[1], args[0]).Execute()
	if err != nil {
		return nil, err
	}

	c.store.OrganizationID = organizationID
	c.store.StackID = args[0]
	c.store.UserID = args[1]

	return c, nil
}

func (c *UpsertStackAccessRolesController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Organization '%s': stack %s, access roles updated for user %s", c.store.OrganizationID, c.store.StackID, c.store.UserID)

	return nil

}
