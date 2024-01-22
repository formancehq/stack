package organizations

import (
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

type UpdateController struct {
	store *DescribeStore
}

var _ fctl.Controller[*DescribeStore] = (*UpdateController)(nil)

func NewDefaultUpdateStore() *DescribeStore {
	return &DescribeStore{}
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
		fctl.WithStringFlag("default-stack-role", "", "Default Stack Role"),
		fctl.WithStringFlag("domain", "", "Organization Domain"),
		fctl.WithStringFlag("default-organization-role", "", "Default Organization Role"),
		fctl.WithController[*DescribeStore](NewUpdateController()),
	)
}

func (c *UpdateController) GetStore() *DescribeStore {
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
		DefaultOrganizationAccess: func() *membershipclient.Role {
			if cmd.Flags().Changed("default-organization-role") {
				s := fctl.GetString(cmd, "default-organization-role")
				return membershipclient.Role(s).Ptr()
			}
			return org.Data.DefaultOrganizationAccess
		}(),
		DefaultStackAccess: func() *membershipclient.Role {
			if cmd.Flags().Changed("default-stack-role") {
				s := fctl.GetString(cmd, "default-stack-role")
				return membershipclient.Role(s).Ptr()
			}
			return org.Data.DefaultStackAccess
		}(),
		Domain: func() *string {
			str := fctl.GetString(cmd, "domain")
			if str != "" {
				return &str
			}
			return org.Data.Domain
		}(),
	}

	response, _, err := apiClient.DefaultApi.
		UpdateOrganization(cmd.Context(), args[0]).OrganizationData(preparedData).Execute()

	if err != nil {
		return nil, err
	}

	c.store.OrganizationExpanded = response.Data

	return c, nil
}

func (c *UpdateController) Render(cmd *cobra.Command, args []string) error {
	return PrintOrganization(c.store.OrganizationExpanded)
}
