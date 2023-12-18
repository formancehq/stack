package users

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type UpdateStore struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}
type UpdateController struct {
	store *UpdateStore
}

var _ fctl.Controller[*UpdateStore] = (*UpdateController)(nil)

func NewDefaultUpdateStore() *UpdateStore {
	return &UpdateStore{}
}

func NewUpdateController() *UpdateController {
	return &UpdateController{
		store: NewDefaultUpdateStore(),
	}
}

func NewUpdateCommand() *cobra.Command {
	return fctl.NewCommand("update <user-id> <roles...>",
		fctl.WithAliases("s"),
		fctl.WithShortDescription("Update user roles by by id within an organization"),
		fctl.WithArgs(cobra.MinimumNArgs(1)),
		fctl.WithController[*UpdateStore](NewUpdateController()),
	)
}

func (c *UpdateController) GetStore() *UpdateStore {
	return c.store
}

func (c *UpdateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	roles := []string{}
	roles = append(roles, args[1:]...)

	response, err := apiClient.DefaultApi.UpsertOrganizationUser(
		cmd.Context(),
		organizationID,
		args[0]).RequestBody(roles).Execute()
	if err != nil {
		return nil, err
	}

	if response.StatusCode > 300 {
		return nil, fmt.Errorf("error updating user: %s", response.Status)
	}

	return c, nil
}

func (c *UpdateController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("User updated.")
	return nil
}
