package invitations

import (
	"encoding/json"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type InvitationSend struct {
	Email       string                              `json:"email"`
	StackClaims []membershipclient.StackClaimsInner `json:"stackClaims"`
	OrgClaim    string                              `json:"orgClaim"`
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
		fctl.WithStringFlag("org-claim", "", "Pre assign organization role e.g. 'ADMIN'"),
		fctl.WithStringFlag("stack-claims", "", "Pre assign stack roles e.g. '[{stackId: <stackId>, roles:[]},...]'"),
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

	invitationClaim := membershipclient.InvitationClaim{}
	orgClaimString := fctl.GetString(cmd, "org-claim")
	if orgClaimString != "" {
		invitationClaim.Roles = []string{orgClaimString}
	}

	stackClaimsStrings := fctl.GetString(cmd, "stack-claims")
	if stackClaimsStrings != "" {
		stackClaims := make([]membershipclient.StackClaimsInner, 0)
		err := json.Unmarshal([]byte(stackClaimsStrings), &stackClaims)
		if err != nil {
			return nil, err
		}
		invitationClaim.StackClaims = stackClaims
	}

	invitations, _, err := apiClient.DefaultApi.
		CreateInvitation(cmd.Context(), organizationID).
		Email(args[0]).InvitationClaim(invitationClaim).Execute()
	if err != nil {
		return nil, err
	}

	c.store.Invitation.Email = invitations.Data.UserEmail
	c.store.Invitation.StackClaims = invitations.Data.StackClaims
	c.store.Invitation.OrgClaim = func() string {
		if len(invitations.Data.Roles) != 0 {
			return ""
		}
		return invitations.Data.Roles[0]
	}()

	return c, nil
}

func (c *SendController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Invitation sent to %s", c.store.Invitation.Email)
	return nil

}
