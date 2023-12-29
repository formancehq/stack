package organizations

import (
	"github.com/formancehq/fctl/cmd/cloud/organizations/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

type UpdateStore struct {
	Organization *membershipclient.Organization `json:"organization"`
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
	return fctl.NewCommand("update <organizationId> --name <name> --default-stack-role <defaultStackRole...> --default-organization-role <defaultOrganizationRole...>",
		fctl.WithAliases("update"),
		fctl.WithShortDescription("Update organization"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithConfirmFlag(),
		fctl.WithStringFlag("name", "", "Organization Name"),
		fctl.WithStringSliceFlag("default-stack-role", []string{}, "Default Stack Role"),
		fctl.WithStringSliceFlag("default-organization-role", []string{}, "Default Organization Role"),
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

	if !fctl.CheckOrganizationApprobation(cmd, "You are about to update an organization") {
		return nil, fctl.ErrMissingApproval
	}

	org, _, err := apiClient.DefaultApi.ReadOrganization(cmd.Context(), args[0]).Execute()
	if err != nil {
		return nil, err
	}

	preparedData := membershipclient.OrganizationData{
		Name: func() string {
			if cmd.Flags().Changed("name") {
				return cmd.Flag("name").Value.String()
			}
			return org.Data.Name
		}(),
		DefaultOrganizationAccess: func() []string {
			if cmd.Flags().Changed("default-organization-role") {
				return fctl.GetStringSlice(cmd, "default-organization-role")
			}
			return org.Data.DefaultOrganizationAccess
		}(),
		DefaultStackAccess: func() []string {
			if cmd.Flags().Changed("default-stack-role") {
				return fctl.GetStringSlice(cmd, "default-stack-role")
			}
			return org.Data.DefaultStackAccess
		}(),
	}

	response, _, err := apiClient.DefaultApi.
		UpdateOrganization(cmd.Context(), args[0]).OrganizationData(preparedData).Execute()

	if err != nil {
		return nil, err
	}

	c.store.Organization = response.Data

	return c, nil
}

func (c *UpdateController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintOrganization(c.store.Organization)
}
