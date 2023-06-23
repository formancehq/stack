package invitations

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type AcceptStore struct {
	Success      bool   `json:"success"`
	InvitationId string `json:"invitation_id"`
}
type AcceptController struct {
	store *AcceptStore
}

var _ fctl.Controller[*AcceptStore] = (*AcceptController)(nil)

func NewDefaultAcceptStore() *AcceptStore {
	return &AcceptStore{}
}

func NewAcceptController() *AcceptController {
	return &AcceptController{
		store: NewDefaultAcceptStore(),
	}
}

func NewAcceptCommand() *cobra.Command {
	return fctl.NewCommand("accept <invitation-id>",
		fctl.WithAliases("a"),
		fctl.WithShortDescription("Accept invitation"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithConfirmFlag(),
		fctl.WithController[*AcceptStore](NewAcceptController()),
	)
}

func (c *AcceptController) GetStore() *AcceptStore {
	return c.store
}

func (c *AcceptController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckOrganizationApprobation(cmd, "You are about to accept an invitation") {
		return nil, fctl.ErrMissingApproval
	}

	_, err = client.DefaultApi.AcceptInvitation(cmd.Context(), args[0]).Execute()
	if err != nil {
		return nil, err
	}

	c.store.InvitationId = args[0]
	c.store.Success = true

	return c, nil
}

func (c *AcceptController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Invitation %s accepted!", c.store.InvitationId)
	return nil

}
