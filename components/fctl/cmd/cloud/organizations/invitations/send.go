package invitations

import (
	"flag"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useSend   = "send <email>"
	shortSend = "Invite a user by email"
)

type InvitationSend struct {
	Email string `json:"email"`
}

type SendStore struct {
	Invitation InvitationSend `json:"invitation"`
}

func NewSendStore() *SendStore {
	return &SendStore{
		Invitation: InvitationSend{},
	}
}
func NewSendConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useSend, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)

	return fctl.NewControllerConfig(
		useSend,
		shortSend,
		shortSend,
		[]string{},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

var _ fctl.Controller[*SendStore] = (*SendController)(nil)

type SendController struct {
	store  *SendStore
	config *fctl.ControllerConfig
}

func NewSendController(config *fctl.ControllerConfig) *SendController {
	return &SendController{
		store:  NewSendStore(),
		config: config,
	}
}

func (c *SendController) GetStore() *SendStore {
	return c.store
}

func (c *SendController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *SendController) Run() (fctl.Renderable, error) {

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

	if !fctl.CheckOrganizationApprobation(flags, "You are about to send an invitation") {
		return nil, fctl.ErrMissingApproval
	}

	_, _, err = apiClient.DefaultApi.
		CreateInvitation(ctx, organizationID).
		Email(args[0]).
		Execute()
	if err != nil {
		return nil, err
	}

	c.store.Invitation.Email = args[0]

	return c, nil
}

func (c *SendController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Invitation sent to %s", c.store.Invitation.Email)
	return nil

}

func NewSendCommand() *cobra.Command {

	config := NewSendConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*SendStore](NewSendController(config)),
	)
}
