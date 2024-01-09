package role

import (
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type UpsertStore struct {
	OrganizationID string `json:"organizationId"`
	UserID         string `json:"userId"`
	StackID        string `json:"stackId"`
}
type UpsertController struct {
	store *UpsertStore
}

var _ fctl.Controller[*UpsertStore] = (*UpsertController)(nil)

func NewDefaultUpsertStore() *UpsertStore {
	return &UpsertStore{}
}

func NewUpsertController() *UpsertController {
	return &UpsertController{
		store: NewDefaultUpsertStore(),
	}
}

func NewUpsertCommand() *cobra.Command {
	return fctl.NewCommand("upsert <stack-id> <user-id> <stack-role>",
		fctl.WithAliases("usar"),
		fctl.WithShortDescription("Update Stack Access Roles within an organization (ADMIN, GUEST, NONE)"),
		fctl.WithArgs(cobra.MinimumNArgs(3)),
		fctl.WithController[*UpsertStore](NewUpsertController()),
	)
}

func (c *UpsertController) GetStore() *UpsertStore {
	return c.store
}

func (c *UpsertController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

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

	_, err = apiClient.DefaultApi.
		UpsertStackUserAccess(cmd.Context(), organizationID, args[0], args[1]).Body(`"` + string(membershipclient.Role(args[2])) + `"`).Execute()
	if err != nil {
		return nil, err
	}

	c.store.OrganizationID = organizationID
	c.store.StackID = args[0]
	c.store.UserID = args[1]

	return c, nil
}

func (c *UpsertController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Organization '%s': stack %s, access roles updated for user %s", c.store.OrganizationID, c.store.StackID, c.store.UserID)

	return nil

}
