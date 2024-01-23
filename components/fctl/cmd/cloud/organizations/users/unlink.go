package users

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type UnlinkStore struct {
	OrganizationID string `json:"organizationId"`
	UserID         string `json:"userId"`
}
type DeleteController struct {
	store *UnlinkStore
}

var _ fctl.Controller[*UnlinkStore] = (*DeleteController)(nil)

func NewDefaultUnlinkStore() *UnlinkStore {
	return &UnlinkStore{}
}

func NewUnlinkController() *DeleteController {
	return &DeleteController{
		store: NewDefaultUnlinkStore(),
	}
}

func NewUnlinkCommand() *cobra.Command {
	return fctl.NewCommand("unlink <user-id>",
		fctl.WithAliases("u", "un"),
		fctl.WithShortDescription("Unlink user from organization"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*UnlinkStore](NewUnlinkController()),
	)
}

func (c *DeleteController) GetStore() *UnlinkStore {
	return c.store
}

func (c *DeleteController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

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

	_, err = apiClient.DefaultApi.DeleteUserFromOrganization(cmd.Context(), organizationID, args[0]).Execute()
	if err != nil {
		return nil, err
	}

	c.store.OrganizationID = organizationID
	c.store.UserID = args[0]

	return c, nil
}

func (c *DeleteController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("User '%s' Deleted from organization '%s'", c.store.UserID, c.store.OrganizationID)

	return nil

}
