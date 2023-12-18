package roles

import (
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DeletedUserAccessStore struct {
	Stack  *membershipclient.Stack `json:"stack"`
	Status string                  `json:"status"`
}
type StackDeleteController struct {
	store *DeletedUserAccessStore
}

var _ fctl.Controller[*DeletedUserAccessStore] = (*StackDeleteController)(nil)

func NewDefaultDeletedUserAccessStore() *DeletedUserAccessStore {
	return &DeletedUserAccessStore{
		Stack:  &membershipclient.Stack{},
		Status: "",
	}
}

func NewStackDeleteController() *StackDeleteController {
	return &StackDeleteController{
		store: NewDefaultDeletedUserAccessStore(),
	}
}

func NewDeleteCommand() *cobra.Command {
	return fctl.NewMembershipCommand("delete <stack-id> <user-id>",
		fctl.WithConfirmFlag(),
		fctl.WithShortDescription("Delete user access from a stack within an organization"),
		fctl.WithAliases("del", "d"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*DeletedUserAccessStore](NewStackDeleteController()),
	)
}
func (c *StackDeleteController) GetStore() *DeletedUserAccessStore {
	return c.store
}

func (c *StackDeleteController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	res, err := apiClient.DefaultApi.DeleteStackUserAccess(cmd.Context(), organizationID, args[0], args[1]).Execute()
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 300 {
		return nil, err
	}

	return c, nil
}

func (c *StackDeleteController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Stack user access deleted.")
	return nil
}
