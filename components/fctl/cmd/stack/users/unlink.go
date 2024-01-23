package users

import (
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type UnlinkStore struct {
	Stack  *membershipclient.Stack `json:"stack"`
	Status string                  `json:"status"`
}
type UnlinkController struct {
	store *UnlinkStore
}

var _ fctl.Controller[*UnlinkStore] = (*UnlinkController)(nil)

func NewDefaultUnlinkStore() *UnlinkStore {
	return &UnlinkStore{
		Stack:  &membershipclient.Stack{},
		Status: "",
	}
}

func NewUnlinkController() *UnlinkController {
	return &UnlinkController{
		store: NewDefaultUnlinkStore(),
	}
}

func NewUnlinkCommand() *cobra.Command {
	return fctl.NewMembershipCommand("unlink <stack-id> <user-id>",
		fctl.WithShortDescription("Unlink stack user within an organization"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*UnlinkStore](NewUnlinkController()),
	)
}
func (c *UnlinkController) GetStore() *UnlinkStore {
	return c.store
}

func (c *UnlinkController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

func (c *UnlinkController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Stack user access deleted.")
	return nil
}
