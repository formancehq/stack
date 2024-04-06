package holds

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type VoidStore struct {
	Success bool   `json:"success"`
	HoldId  string `json:"holdId"`
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

	store := fctl.GetStackStore(cmd.Context())
	request := operations.VoidHoldRequest{
		HoldID: args[0],
	}
	_, err := store.Client().Wallets.VoidHold(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "voiding hold")
	}

	c.store.Success = true
	c.store.HoldId = args[0]

	return c, nil
}

func (c *VoidController) Render(cmd *cobra.Command, args []string) error {

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Hold '%s' voided!", args[0])

	return nil
}
