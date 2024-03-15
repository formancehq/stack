package transactions

import (
	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DeleteMetadataStore struct {
	Success bool `json:"success"`
}
type DeleteMetadataController struct {
	store *DeleteMetadataStore
}

var _ fctl.Controller[*DeleteMetadataStore] = (*DeleteMetadataController)(nil)

func NewDefaultDeleteMetadataStore() *DeleteMetadataStore {
	return &DeleteMetadataStore{}
}

func NewDeleteMetadataController() *DeleteMetadataController {
	return &DeleteMetadataController{
		store: NewDefaultDeleteMetadataStore(),
	}
}

func NewDeleteMetadataCommand() *cobra.Command {
	return fctl.NewCommand("delete-metadata <transaction-id> [<key>...]",
		fctl.WithShortDescription("Delete metadata on transaction (Start from ledger v2 api)"),
		fctl.WithAliases("dm", "del-meta"),
		fctl.WithConfirmFlag(),
		fctl.WithValidArgs("last"),
		fctl.WithArgs(cobra.MinimumNArgs(2)),
		fctl.WithController[*DeleteMetadataStore](NewDeleteMetadataController()),
	)
}

func (c *DeleteMetadataController) GetStore() *DeleteMetadataStore {
	return c.store
}

func (c *DeleteMetadataController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := fctl.GetStackStore(cmd.Context())

	transactionID, err := internal.TransactionIDOrLastN(cmd.Context(), store.Client(),
		fctl.GetString(cmd, internal.LedgerFlag), args[0])
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to set a metadata on transaction %d", transactionID) {
		return nil, fctl.ErrMissingApproval
	}

	response, err := store.Client().Ledger.V2DeleteTransactionMetadata(cmd.Context(), operations.V2DeleteTransactionMetadataRequest{
		ID:     transactionID,
		Key:    args[1],
		Ledger: fctl.GetString(cmd, internal.LedgerFlag),
	})
	if err != nil {
		return nil, err
	}

	c.store.Success = (response.StatusCode % 200) < 100
	return c, nil
}

func (c *DeleteMetadataController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Metadata deleted!")
	return nil
}
