package organizations

import (
	"github.com/formancehq/fctl/cmd/cloud/organizations/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/go-libs/pointer"
	"github.com/spf13/cobra"
)

type CreateStore struct {
	Organization *membershipclient.OrganizationExpanded `json:"organization"`
}
type CreateController struct {
	store *CreateStore
}

var _ fctl.Controller[*CreateStore] = (*CreateController)(nil)

func NewDefaultCreateStore() *CreateStore {
	return &CreateStore{}
}

func NewCreateController() *CreateController {
	return &CreateController{
		store: NewDefaultCreateStore(),
	}
}

func NewCreateCommand() *cobra.Command {
	return fctl.NewCommand(`create <name> --default-stack-role "ADMIN" --default-organization-role "ADMIN"`,
		fctl.WithAliases("cr", "c"),
		fctl.WithShortDescription("Create organization"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithStringFlag("default-stack-role", "", "Default Stack Role roles: (ADMIN,GUEST)"),
		fctl.WithStringFlag("domain", "", "Organization Domain"),
		fctl.WithStringFlag("default-organization-role", "", "Default Organization Role roles: (ADMIN,GUEST)"),
		fctl.WithConfirmFlag(),
		fctl.WithController[*CreateStore](NewCreateController()),
	)
}

func (c *CreateController) GetStore() *CreateStore {
	return c.store
}

func (c *CreateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := fctl.GetMembershipStore(cmd.Context())
	if !fctl.CheckOrganizationApprobation(cmd, "You are about to create a new organization") {
		return nil, fctl.ErrMissingApproval
	}

	defaultStackRole := fctl.GetString(cmd, "default-stack-role")
	defaultOrganizationRole := fctl.GetString(cmd, "default-organization-role")
	domain := fctl.GetString(cmd, "domain")

	orgData := membershipclient.OrganizationData{
		Name: args[0],
	}

	if defaultOrganizationRole != "" {
		orgData.DefaultOrganizationAccess = membershipclient.Role(defaultOrganizationRole).Ptr()
	}

	if defaultStackRole != "" {
		orgData.DefaultStackAccess = membershipclient.Role(defaultStackRole).Ptr()
	}

	if domain != "" {
		orgData.Domain = pointer.For(domain)
	}

	response, _, err := store.Client().
		CreateOrganization(cmd.Context()).
		Body(orgData).Execute()
	if err != nil {
		return nil, err
	}

	c.store.Organization = response.Data

	return c, nil
}

func (c *CreateController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintOrganization(c.store.Organization)
}
