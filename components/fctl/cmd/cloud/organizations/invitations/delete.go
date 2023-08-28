package invitations

import (
	"flag"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useDelete   = "delete <invitation-id>"
	shortDelete = "Delete an invitation by id"
)

type DeleteStore struct {
	Success        bool   `json:"success"`
	OrganizationID string `json:"organizationID"`
}

func NewDeleteStore() *DeleteStore {
	return &DeleteStore{
		Success:        false,
		OrganizationID: "",
	}
}

func NewDeleteConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useDelete, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)
	return fctl.NewControllerConfig(
		useDelete,
		shortDelete,
		shortDelete,
		[]string{
			"del",
		},
		flags,
		fctl.Organization, fctl.Stack,
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

	if !fctl.CheckOrganizationApprobation(flags, "You are about to delete an invitation") {
		return nil, fctl.ErrMissingApproval
	}

	_, err = apiClient.DefaultApi.
		DeleteInvitation(ctx, organizationID, args[0]).
		Execute()
	if err != nil {
		return nil, err
	}

	c.store.Success = true
	c.store.OrganizationID = organizationID

	return c, nil
}

func (c *DeleteController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Invitation %s deleted", c.config.GetArgs()[0])
	return nil
}

func NewDeleteCommand() *cobra.Command {
	config := NewDeleteConfig()

	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*DeleteStore](NewDeleteController(config)),
	)
}
