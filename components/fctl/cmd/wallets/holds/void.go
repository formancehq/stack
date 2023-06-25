package holds

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type VoidStore struct {
	Success bool   `json:"success"`
	HoldId  string `json:"hold_id"`
}
type VoidController struct {
	store *VoidStore
}

var _ fctl.Controller[*VoidStore] = (*VoidController)(nil)

func NewDefaultVoidStore() *VoidStore {
	return &VoidStore{}
}

func NewVoidController() *VoidController {
	return &VoidController{
		store: NewDefaultVoidStore(),
	}
}

func NewVoidCommand() *cobra.Command {
	return fctl.NewCommand("void <hold-id>",
		fctl.WithShortDescription("Void a hold"),
		fctl.WithAliases("v"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*VoidStore](NewVoidController()),
	)
}

func (c *VoidController) GetStore() *VoidStore {
	return c.store
}

func (c *VoidController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config")
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	stackClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	request := operations.VoidHoldRequest{
		HoldID: args[0],
	}
	response, err := stackClient.Wallets.VoidHold(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "voiding hold")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true //Todo: check status code 204/200
	c.store.HoldId = args[0]

	return c, nil
}

func (c *VoidController) Render(cmd *cobra.Command, args []string) error {

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Hold '%s' voided!", args[0])

	return nil
}
