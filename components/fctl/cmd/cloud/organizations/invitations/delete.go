package invitations

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DeleteStore struct {
	Success        bool   `json:"success"`
	OrganizationID string `json:"organizationID"`
}
type DeleteController struct {
	store           *DeleteStore
	endpointFlag    string
	defaultEndpoint string
}

func NewDefaultDeleteStore() *DeleteStore {
	return &DeleteStore{
		Success:        false,
		OrganizationID: "",
	}
}
func NewDeleteController() *DeleteController {
	return &DeleteController{
		store:           NewDefaultDeleteStore(),
		endpointFlag:    "endpoint",
		defaultEndpoint: "https://api.sandbox.mangopay.com",
	}
}

func NewDeleteCommand() *cobra.Command {
	return fctl.NewCommand("delete <id>",
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithShortDescription("Delete an invitation"),
		fctl.WithAliases("del"),
		fctl.WithConfirmFlag(),
		fctl.WithController[*DeleteStore](NewDeleteController()),
	)
}

func (c *DeleteController) GetStore() *DeleteStore {
	return c.store
}

func (c *DeleteController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetMembershipStore(cmd.Context())

	organizationID, err := fctl.ResolveOrganizationID(cmd, store.Config, store.Client())
	if err != nil {
		return nil, err
	}

	if !fctl.CheckOrganizationApprobation(cmd, "You are about to delete an invitation") {
		return nil, fctl.ErrMissingApproval
	}

	_, err = store.Client().
		DeleteInvitation(cmd.Context(), organizationID, args[0]).
		Execute()
	if err != nil {
		return nil, err
	}

	c.store.Success = true
	c.store.OrganizationID = organizationID

	return c, nil
}

func (c *DeleteController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Invitation %s deleted", args[0])
	return nil
}
