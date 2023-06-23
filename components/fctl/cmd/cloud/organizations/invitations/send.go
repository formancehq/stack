package invitations

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type InvitationSend struct {
	Email string `json:"email"`
}

type SendStore struct {
	Invitation InvitationSend `json:"invitation"`
}
type SendController struct {
	store *SendStore
}

var _ fctl.Controller[*SendStore] = (*SendController)(nil)

func NewDefaultSendStore() *SendStore {
	return &SendStore{
		Invitation: InvitationSend{},
	}
}

func NewSendController() *SendController {
	return &SendController{
		store: NewDefaultSendStore(),
	}
}

func NewSendCommand() *cobra.Command {
	return fctl.NewCommand("send <email>",
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithShortDescription("Invite a user by email"),
		fctl.WithAliases("s"),
		fctl.WithConfirmFlag(),
		fctl.WithController[*SendStore](NewSendController()),
	)
}

func (c *SendController) GetStore() *SendStore {
	return c.store
}

func (c *SendController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckOrganizationApprobation(cmd, "You are about to send an invitation") {
		return nil, fctl.ErrMissingApproval
	}

	_, _, err = apiClient.DefaultApi.
		CreateInvitation(cmd.Context(), organizationID).
		Email(args[0]).
		Execute()
	if err != nil {
		return nil, err
	}

	c.store.Invitation.Email = args[0]

	return c, nil
}

func (c *SendController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Invitation sent to %s", c.store.Invitation.Email)
	return nil

}
