package transactions

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/spf13/cobra"
)

type RevertStore struct {
	Transaction *shared.Transaction `json:"transaction"`
}
type RevertController struct {
	store *RevertStore
}

var _ fctl.Controller[*RevertStore] = (*RevertController)(nil)

func NewDefaultRevertStore() *RevertStore {
	return &RevertStore{}
}

func NewRevertController() *RevertController {
	return &RevertController{
		store: NewDefaultRevertStore(),
	}
}

func NewRevertCommand() *cobra.Command {
	return fctl.NewCommand("revert <transaction-id>",
		fctl.WithConfirmFlag(),
		fctl.WithShortDescription("Revert a transaction"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithValidArgs("last"),
		fctl.WithController[*RevertStore](NewRevertController()),
	)
}

func (c *RevertController) GetStore() *RevertStore {
	return c.store
}

func (c *RevertController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to revert transaction %s", args[0]) {
		return nil, fctl.ErrMissingApproval
	}

	ledger := fctl.GetString(cmd, internal.LedgerFlag)
	txId, err := internal.TransactionIDOrLastN(cmd.Context(), store.Client(), ledger, args[0])
	if err != nil {
		return nil, err
	}

	request := operations.RevertTransactionRequest{
		Ledger: ledger,
		Txid:   txId,
	}

	response, err := store.Client().Ledger.RevertTransaction(cmd.Context(), request)
	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Transaction = &response.TransactionResponse.Data

	return c, nil
}

func (c *RevertController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintTransaction(cmd.OutOrStdout(), *c.store.Transaction)
}
