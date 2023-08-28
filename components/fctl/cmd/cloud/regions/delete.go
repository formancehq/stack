package regions

import (
	"flag"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useDelete   = "delete [region-id]"
	shortDelete = "Delete a private region with id"
)

type DeleteStore struct {
	Success bool `json:"success"`
}

func NewDeleteStore() *DeleteStore {
	return &DeleteStore{}
}
func NewDeleteConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useDelete, flag.ExitOnError)
	return fctl.NewControllerConfig(
		useDelete,
		shortDelete,
		shortDelete,
		[]string{
			"del", "d",
		},
		flags,
		fctl.Organization,
	)
}

var _ fctl.Controller[*DeleteStore] = (*DeleteController)(nil)

type DeleteController struct {
	store  *DeleteStore
	config *fctl.ControllerConfig
}

func NewDeleteController(config *fctl.ControllerConfig) *DeleteController {
	return &DeleteController{
		store:  NewDeleteStore(),
		config: config,
	}
}

func (c *DeleteController) GetStore() *DeleteStore {
	return c.store
}

func (c *DeleteController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *DeleteController) Run() (fctl.Renderable, error) {

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

	_, err = apiClient.DefaultApi.DeleteRegion(ctx, organizationID, args[0]).Execute()
	if err != nil {
		return nil, err
	}

	c.store.Success = true

	return c, nil
}

func (c *DeleteController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Region deleted successfully!")

	return nil

}

func NewDeleteCommand() *cobra.Command {

	config := NewDeleteConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*DeleteStore](NewDeleteController(config)),
	)
}
