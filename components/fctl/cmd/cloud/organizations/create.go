package organizations

import (
	"flag"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useCreate   = "create <name>"
	shortCreate = "Create a new organization"
)

type CreateStore struct {
	OrganizationId   string `json:"organizationId"`
	OrganizationName string `json:"organizationName"`
}

func NewCreateStore() *CreateStore {
	return &CreateStore{}
}

func NewCreateConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useCreate, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)

	return fctl.NewControllerConfig(
		useCreate,
		shortCreate,
		shortCreate,
		[]string{
			"cr", "c",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

var _ fctl.Controller[*CreateStore] = (*CreateController)(nil)

type CreateController struct {
	store  *CreateStore
	config *fctl.ControllerConfig
}

func NewCreateController(config *fctl.ControllerConfig) *CreateController {
	return &CreateController{
		store:  NewCreateStore(),
		config: config,
	}
}

func (c *CreateController) GetStore() *CreateStore {
	return c.store
}

func (c *CreateController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *CreateController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	if !fctl.CheckOrganizationApprobation(flags, "You are about to create a new organization") {
		return nil, fctl.ErrMissingApproval
	}

	response, _, err := apiClient.DefaultApi.
		CreateOrganization(ctx).
		Body(membershipclient.OrganizationData{
			Name: args[0],
		}).Execute()
	if err != nil {
		return nil, err
	}

	c.store.OrganizationId = response.Data.Id
	c.store.OrganizationName = args[0]

	return c, nil
}

func (c *CreateController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Organization '%s' created with ID: %s", c.store.OrganizationName, c.store.OrganizationId)

	return nil

}

func NewCreateCommand() *cobra.Command {
	config := NewCreateConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*CreateStore](NewCreateController(config)),
	)
}
