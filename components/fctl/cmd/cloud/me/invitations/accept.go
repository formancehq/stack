package invitations

import (
	"flag"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useAccept   = "accept <invitation-id>"
	shortAccept = "Accept an invitation"
)

type AcceptStore struct {
	Success      bool   `json:"success"`
	InvitationId string `json:"invitationId"`
}

func NewAcceptStore() *AcceptStore {
	return &AcceptStore{}
}

func NewAcceptConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useAccept, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)
	return fctl.NewControllerConfig(
		useAccept,
		shortAccept,
		shortAccept,
		[]string{
			"a",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

var _ fctl.Controller[*AcceptStore] = (*AcceptController)(nil)

type AcceptController struct {
	store  *AcceptStore
	config *fctl.ControllerConfig
}

func NewAcceptController(config *fctl.ControllerConfig) *AcceptController {
	return &AcceptController{
		store:  NewAcceptStore(),
		config: config,
	}
}

func (c *AcceptController) GetStore() *AcceptStore {
	return c.store
}

func (c *AcceptController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *AcceptController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewMembershipClient(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	if !fctl.CheckOrganizationApprobation(flags, "You are about to accept an invitation") {
		return nil, fctl.ErrMissingApproval
	}

	_, err = client.DefaultApi.AcceptInvitation(ctx, args[0]).Execute()
	if err != nil {
		return nil, err
	}

	c.store.InvitationId = args[0]
	c.store.Success = true

	return c, nil
}

func (c *AcceptController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Invitation %s accepted!", c.store.InvitationId)
	return nil

}
func NewAcceptCommand() *cobra.Command {
	config := NewAcceptConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*AcceptStore](NewAcceptController(config)),
	)
}
