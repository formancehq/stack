package organizations

import (
	"github.com/formancehq/fctl/cmd/cloud/organizations/internal"
	"github.com/formancehq/fctl/membershipclient"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

type DescribeStore struct {
	*membershipclient.OrganizationExpanded
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
		fctl.WithBoolFlag("expand", false, "Expand the organization"),
		fctl.WithController[*DescribeStore](NewDescribeController()),
	)
}

func (c *DescribeController) GetStore() *DescribeStore {
	return c.store
}

func (c *DescribeController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := fctl.GetMembershipStore(cmd.Context())

	expand := fctl.GetBool(cmd, "expand")
	response, _, err := store.Client().
		ReadOrganization(cmd.Context(), args[0]).Expand(expand).Execute()
	if err != nil {
		return nil, err
	}

	c.store.OrganizationExpanded = response.Data
	return c, nil
}

func (c *DescribeController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintOrganization(c.store.OrganizationExpanded)
}
