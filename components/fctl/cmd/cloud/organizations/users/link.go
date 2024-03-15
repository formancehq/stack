package users

import (
	"fmt"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type LinkStore struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
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
	return fctl.NewCommand("link <user-id>",
		fctl.WithStringFlag("role", "", "Roles: (ADMIN, GUEST, NONE)"),
		fctl.WithShortDescription("Link user to an organization with properties"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithPreRunE(func(cmd *cobra.Command, args []string) error {
			store := fctl.GetMembershipStore(cmd.Context())

			version := fctl.MembershipServerInfo(cmd.Context(), store.Client())
			if !semver.IsValid(version) {
				return nil
			}

			if semver.Compare(version, "v0.26.1") >= 0 {
				return nil
			}

			return fmt.Errorf("unsupported membership server version: %s", version)
		}),
		fctl.WithController[*LinkStore](NewLinkController()),
	)
}

func (c *LinkController) GetStore() *LinkStore {
	return c.store
}

func (c *LinkController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetMembershipStore(cmd.Context())
	organizationID, err := fctl.ResolveOrganizationID(cmd, store.Config, store.Client())
	if err != nil {
		return nil, err
	}

	role := fctl.GetString(cmd, "role")
	req := membershipclient.UpdateOrganizationUserRequest{}
	if role != "" {
		req.Role = membershipclient.Role(role)
	} else {
		return nil, fmt.Errorf("role is required")
	}
	response, err := store.Client().UpsertOrganizationUser(
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

func (c *LinkController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("User Addd.")
	return nil
}
