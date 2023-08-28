package regions

import (
	"flag"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useCreate   = "create [name]"
	shortCreate = "Create a new region with name"
)

type CreateStore struct {
	RegionId string `json:"regionId"`
	Secret   string `json:"secret"`
}

func NewCreateStore() *CreateStore {
	return &CreateStore{}
}
func NewCreateConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useCreate, flag.ExitOnError)
	return fctl.NewControllerConfig(
		useCreate,
		shortCreate,
		shortCreate,
		[]string{
			"sh", "s",
		},
		flags,
		fctl.Organization,
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

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	name := ""
	if len(args) > 0 {
		name = args[0]
	} else {
		name, err = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter a name")
		if err != nil {
			return nil, err
		}
	}

	regionResponse, _, err := apiClient.DefaultApi.CreatePrivateRegion(ctx, organizationID).
		CreatePrivateRegionRequest(membershipclient.CreatePrivateRegionRequest{
			Name: name,
		}).
		Execute()
	if err != nil {
		return nil, err
	}

	c.store.RegionId = regionResponse.Data.Id
	c.store.Secret = *regionResponse.Data.Secret.Clear

	return c, nil
}

func (c *CreateController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln(
		"Region created successfully with ID: %s", c.store.RegionId)
	pterm.Success.WithWriter(c.config.GetOut()).Printfln(
		"Your secret is (keep it safe, we will not be able to give it to you again): %s", c.store.Secret)

	return nil
}

func NewCreateCommand() *cobra.Command {

	config := NewCreateConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.RangeArgs(0, 1)),
		fctl.WithController[*CreateStore](NewCreateController(config)),
	)
}
