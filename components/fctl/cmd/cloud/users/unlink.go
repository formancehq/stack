package users

import (
	"flag"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useUnlink   = "unlink <user-id>"
	shortUnlink = "Unlink user from organization"
)

type UnlinkStore struct {
	OrganizationID string `json:"organizationId"`
	UserID         string `json:"userId"`
}

func NewUnlinkStore() *UnlinkStore {
	return &UnlinkStore{}
}
func NewUnlinkConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useUnlink, flag.ExitOnError)
	return fctl.NewControllerConfig(
		useUnlink,
		shortUnlink,
		shortUnlink,
		[]string{
			"u", "un",
		},
		flags,
		fctl.Organization,
	)
}

var _ fctl.Controller[*UnlinkStore] = (*UnlinkController)(nil)

type UnlinkController struct {
	store  *UnlinkStore
	config *fctl.ControllerConfig
}

func NewUnlinkController(config *fctl.ControllerConfig) *UnlinkController {
	return &UnlinkController{
		store:  NewUnlinkStore(),
		config: config,
	}
}

func (c *UnlinkController) GetStore() *UnlinkStore {
	return c.store
}

func (c *UnlinkController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *UnlinkController) Run() (fctl.Renderable, error) {
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

	_, err = apiClient.DefaultApi.UnlinkUserFromOrganization(ctx, organizationID, args[0]).Execute()
	if err != nil {
		return nil, err
	}

	c.store.OrganizationID = organizationID
	c.store.UserID = args[0]

	return c, nil
}

func (c *UnlinkController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("User '%s' unlinked from organization '%s'", c.store.UserID, c.store.OrganizationID)

	return nil

}

func NewUnlinkCommand() *cobra.Command {
	config := NewUnlinkConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*UnlinkStore](NewUnlinkController(config)),
	)
}
