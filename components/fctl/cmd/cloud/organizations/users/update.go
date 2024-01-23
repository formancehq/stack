package users

import (
	"fmt"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
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
	return fctl.NewCommand("update <user-id>",
		fctl.WithAliases("s"),
		fctl.WithStringFlag("role", "", "Roles: (ADMIN, GUEST, NONE)"),
		fctl.WithShortDescription("Update user roles by by id within an organization"),
		fctl.WithArgs(cobra.MinimumNArgs(1)),
		fctl.WithPreRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}

			apiClient, err := fctl.NewMembershipClient(cmd, cfg)
			if err != nil {
				return err
			}

			version := fctl.MembershipServerInfo(cmd.Context(), apiClient)
			if !semver.IsValid(version) {
				return nil
			}

			if semver.Compare(version, "v0.26.1") >= 0 {
				return nil
			}

			return fmt.Errorf("unsupported membership server version: %s", version)
		}),
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

	role := fctl.GetString(cmd, "role")
	req := membershipclient.UpdateOrganizationUserRequest{}
	if role != "" {
		req.Role = membershipclient.Role(role)
	} else {
		access, res, err := apiClient.DefaultApi.ReadUserOfOrganization(cmd.Context(), organizationID, args[0]).Execute()
		if err != nil {
			return nil, err
		}
		if res.StatusCode > 300 {
			return nil, fmt.Errorf("error reading user: %s", res.Status)
		}
		req.Role = access.Data.Role
	}
	response, err := apiClient.DefaultApi.UpsertOrganizationUser(
		cmd.Context(),
		organizationID,
		args[0]).
		UpdateOrganizationUserRequest(req).Execute()
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
