package pools

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DeleteStore struct {
	PoolID  string `json:"poolID"`
	Success bool   `json:"success"`
}

type DeleteController struct {
	PaymentsVersion versions.Version

	store *DeleteStore
}

func (c *DeleteController) SetVersion(version versions.Version) {
	c.PaymentsVersion = version
}

var _ fctl.Controller[*DeleteStore] = (*DeleteController)(nil)

func NewDeleteStore() *DeleteStore {
	return &DeleteStore{}
}

func NewDeleteController() *DeleteController {
	return &DeleteController{
		store: NewDeleteStore(),
	}
}
func NewDeleteCommand() *cobra.Command {
	c := NewDeleteController()
	return fctl.NewCommand("delete <poolID>",
		fctl.WithConfirmFlag(),
		fctl.WithAliases("d"),
		fctl.WithShortDescription("Delete a pool"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*DeleteStore](c),
	)
}

func (c *DeleteController) GetStore() *DeleteStore {
	return c.store
}

func (c *DeleteController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	if err := versions.GetPaymentsVersion(cmd, args, c); err != nil {
		return nil, err
	}

	if c.PaymentsVersion < versions.V1 {
		return nil, fmt.Errorf("pools are only supported in >= v1.0.0")
	}

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to delete '%s'", args[0]) {
		return nil, fctl.ErrMissingApproval
	}

	response, err := store.Client().Payments.DeletePool(
		cmd.Context(),
		operations.DeletePoolRequest{
			PoolID: args[0],
		},
	)

	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.PoolID = args[0]
	c.store.Success = true

	return c, nil
}

func (c *DeleteController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Pool %s Deleted!", c.store.PoolID)
	return nil
}
