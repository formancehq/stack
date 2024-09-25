package transferinitiation

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DeleteStore struct {
	TransferID string `json:"transferId"`
	Success    bool   `json:"success"`
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
	return fctl.NewCommand("delete <transferID>",
		fctl.WithConfirmFlag(),
		fctl.WithAliases("d"),
		fctl.WithShortDescription("Delete a transfer Initiation"),
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
		return nil, fmt.Errorf("transfer initiation are only supported in >= v1.0.0")
	}

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to delete '%s'", args[0]) {
		return nil, fctl.ErrMissingApproval
	}

	response, err := store.Client().Payments.V1.DeleteTransferInitiation(
		cmd.Context(),
		operations.DeleteTransferInitiationRequest{
			TransferID: args[0],
		},
	)

	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.TransferID = args[0]
	c.store.Success = true

	return c, nil
}

func (c *DeleteController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Transfer Initiation %s Deleted!", c.store.TransferID)
	return nil
}
