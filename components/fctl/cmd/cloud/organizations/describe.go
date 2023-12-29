package organizations

import (
	"github.com/formancehq/fctl/cmd/cloud/organizations/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

type DescribeStore struct {
	Organization *membershipclient.Organization `json:"organization"`
}
type DescribeController struct {
	store *DescribeStore
}

var _ fctl.Controller[*DescribeStore] = (*DescribeController)(nil)

func NewDefaultDescribeStore() *DescribeStore {
	return &DescribeStore{}
}

func NewDescribeController() *DescribeController {
	return &DescribeController{
		store: NewDefaultDescribeStore(),
	}
}

func NewDescribeCommand() *cobra.Command {
	return fctl.NewCommand("describe <organizationId>",
		fctl.WithShortDescription("Describe organization"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithConfirmFlag(),
		fctl.WithController[*DescribeStore](NewDescribeController()),
	)
}

func (c *DescribeController) GetStore() *DescribeStore {
	return c.store
}

func (c *DescribeController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}

	response, _, err := apiClient.DefaultApi.
		ReadOrganization(cmd.Context(), args[0]).Execute()
	if err != nil {
		return nil, err
	}

	c.store.Organization = response.Data

	return c, nil
}

func (c *DescribeController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintOrganization(c.store.Organization)
}
