package transferinitiation

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type RetryStore struct {
	TransferID string `json:"transferId"`
	Success    bool   `json:"success"`
}
type RetryController struct {
	PaymentsVersion versions.Version

	store *RetryStore
}

func (c *RetryController) SetVersion(version versions.Version) {
	c.PaymentsVersion = version
}

var _ fctl.Controller[*RetryStore] = (*RetryController)(nil)

func NewRetryStore() *RetryStore {
	return &RetryStore{}
}

func NewRetryController() *RetryController {
	return &RetryController{
		store: NewRetryStore(),
	}
}

func NewRetryCommand() *cobra.Command {
	c := NewRetryController()
	return fctl.NewCommand("retry <transferID>",
		fctl.WithConfirmFlag(),
		fctl.WithShortDescription("Retry a failed transfer initiation"),
		fctl.WithAliases("r"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*RetryStore](c),
	)
}

func (c *RetryController) GetStore() *RetryStore {
	return c.store
}

func (c *RetryController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())
	if err := versions.GetPaymentsVersion(cmd, args, c); err != nil {
		return nil, err
	}

	if c.PaymentsVersion < versions.V1 {
		return nil, fmt.Errorf("transfer initiation are only supported in >= v1.0.0")
	}

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to retry the transfer initiation '%s'", args[0]) {
		return nil, fctl.ErrMissingApproval
	}

	//nolint:gosimple
	response, err := store.Client().Payments.RetryTransferInitiation(cmd.Context(), operations.RetryTransferInitiationRequest{
		TransferID: args[0],
	})
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

func (c *RetryController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Retry Transfer Initiation with ID: %s", c.store.TransferID)

	return nil
}
