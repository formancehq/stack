package transactions

import (
	"github.com/formancehq/stack/libs/go-libs/pointer"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/spf13/cobra"
)

type RevertStore struct {
	Transaction internal.Transaction `json:"transaction"`
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
		fctl.WithBoolFlag("at-effective-date", false, "set the timestamp to the original transaction timestamp"),
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

	if fctl.GetBool(cmd, "at-effective-date") {
		request := operations.V2RevertTransactionRequest{
			Ledger:          ledger,
			ID:              txId,
			AtEffectiveDate: pointer.For(true),
		}

		response, err := store.Client().Ledger.V2RevertTransaction(cmd.Context(), request)
		if err != nil {
			return nil, err
		}

		c.store.Transaction = internal.WrapV2Transaction(response.V2RevertTransactionResponse.Data)
	} else {
		request := operations.RevertTransactionRequest{
			Ledger: ledger,
			Txid:   txId,
		}

		response, err := store.Client().Ledger.RevertTransaction(cmd.Context(), request)
		if err != nil {
			return nil, err
		}

		c.store.Transaction = internal.WrapV1Transaction(response.TransactionResponse.Data)
	}

	return c, nil
}

func (c *RevertController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintTransaction(cmd.OutOrStdout(), c.store.Transaction)
}
