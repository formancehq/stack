package transactions

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/spf13/cobra"
)

type RevertStore struct {
	Transaction *shared.V2Transaction `json:"transaction"`
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
	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to revert transaction %s", args[0]) {
		return nil, fctl.ErrMissingApproval
	}

	ledgerClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	ledger := fctl.GetString(cmd, internal.LedgerFlag)
	txId, err := internal.TransactionIDOrLastN(cmd.Context(), ledgerClient, ledger, args[0])
	if err != nil {
		return nil, err
	}

	request := operations.V2RevertTransactionRequest{
		Ledger: ledger,
		ID:     txId,
	}

	response, err := ledgerClient.Ledger.V2RevertTransaction(cmd.Context(), request)
	if response.V2ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.V2ErrorResponse.ErrorCode, response.V2ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Transaction = &response.V2RevertTransactionResponse.Data

	return c, nil
}

func (c *RevertController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintTransaction(cmd.OutOrStdout(), *c.store.Transaction)
}
