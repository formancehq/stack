package organizations

import (
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
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
	return fctl.NewCommand("create <name> --default-stack-role \"ADMIN,GUEST\" --default-organization-role \"ADMIN,GUEST\"",
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

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckOrganizationApprobation(cmd, "You are about to create a new organization") {
		return nil, fctl.ErrMissingApproval
	}

	defaultStackRole := fctl.GetString(cmd, "default-stack-role")
	defaultOrganizationRole := fctl.GetString(cmd, "default-organization-role")
	response, _, err := apiClient.DefaultApi.
		CreateOrganization(cmd.Context()).
		Body(membershipclient.OrganizationData{
			Name:                      args[0],
			DefaultOrganizationAccess: membershipclient.Role(defaultOrganizationRole).Ptr(),
			DefaultStackAccess:        membershipclient.Role(defaultStackRole).Ptr(),
			Domain: func() *string {
				str := fctl.GetString(cmd, "domain")
				if str == "" {
					return nil
				}
				return &str
			}(),
		}).Execute()
	if err != nil {
		return nil, err
	}

	c.store.Organization = response.Data

	return c, nil
}

func (c *CreateController) Render(cmd *cobra.Command, args []string) error {
	return PrintOrganization(c.store.Organization)
}
