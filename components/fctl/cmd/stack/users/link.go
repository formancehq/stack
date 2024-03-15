package users

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/stack/store"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type LinkStore struct {
	OrganizationID string `json:"organizationId"`
	UserID         string `json:"userId"`
	StackID        string `json:"stackId"`
}
type LinkController struct {
	store *LinkStore
}

var _ fctl.Controller[*LinkStore] = (*LinkController)(nil)

func NewDefaultLinkStore() *LinkStore {
	return &LinkStore{}
}

func NewLinkController() *LinkController {
	return &LinkController{
		store: NewDefaultLinkStore(),
	}
}

func NewLinkCommand() *cobra.Command {
	return fctl.NewCommand("link <stack-id> <user-id>",
		fctl.WithStringFlag("role", "", "Roles: (ADMIN, GUEST, NONE)"),
		fctl.WithShortDescription("Link stack user with properties"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*LinkStore](NewLinkController()),
	)
}

func (c *LinkController) GetStore() *LinkStore {
	return c.store
}

func (c *LinkController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := store.GetStore(cmd.Context())

	role := membershipclient.Role(fctl.GetString(cmd, "role"))
	req := membershipclient.UpdateStackUserRequest{}
	if role != "" {
		req.Role = role
	} else {
		return nil, fmt.Errorf("role is required")
	}

	_, err := store.Client().
		UpsertStackUserAccess(cmd.Context(), store.OrganizationId(), args[0], args[1]).
		UpdateStackUserRequest(req).Execute()
	if err != nil {
		return nil, err
	}

	c.store.OrganizationID = store.OrganizationId()
	c.store.StackID = args[0]
	c.store.UserID = args[1]

	return c, nil
}

func (c *LinkController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Organization '%s': stack %s, access roles updated for user %s", c.store.OrganizationID, c.store.StackID, c.store.UserID)

	return nil

}
