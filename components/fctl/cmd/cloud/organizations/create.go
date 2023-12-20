package organizations

import (
	"github.com/formancehq/fctl/cmd/cloud/organizations/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

type CreateStore struct {
	Organization *membershipclient.Organization `json:"organization"`
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
	return fctl.NewCommand("create <name> --default-stack-role <defaultStackRole...> --default-organization-role <defaultOrganizationRole...>",
		fctl.WithAliases("cr", "c"),
		fctl.WithShortDescription("Create organization"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithStringSliceFlag("default-stack-role", []string{}, "Default Stack Role"),
		fctl.WithStringSliceFlag("default-organization-role", []string{}, "Default Organization Role"),
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

	response, _, err := apiClient.DefaultApi.
		CreateOrganization(cmd.Context()).
		Body(membershipclient.OrganizationData{
			Name:                      args[0],
			DefaultOrganizationAccess: fctl.GetStringSlice(cmd, "default-organization-role"),
			DefaultStackAccess:        fctl.GetStringSlice(cmd, "default-stack-role"),
		}).Execute()
	if err != nil {
		return nil, err
	}

	c.store.Organization = response.Data

	return c, nil
}

func (c *CreateController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintOrganization(c.store.Organization)
}
