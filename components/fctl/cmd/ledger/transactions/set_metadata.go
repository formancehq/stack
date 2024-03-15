package transactions

import (
	"fmt"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	"github.com/formancehq/fctl/cmd/ledger/store"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type SetMetadataStore struct {
	Success bool `json:"success"`
}
type SetMetadataController struct {
	store *SetMetadataStore
}

var _ fctl.Controller[*SetMetadataStore] = (*SetMetadataController)(nil)

func NewDefaultSetMetadataStore() *SetMetadataStore {
	return &SetMetadataStore{}
}

func NewSetMetadataController() *SetMetadataController {
	return &SetMetadataController{
		store: NewDefaultSetMetadataStore(),
	}
}

func NewSetMetadataCommand() *cobra.Command {
	return fctl.NewCommand("set-metadata <transaction-id> [<key>=<value>...]",
		fctl.WithShortDescription("Set metadata on transaction"),
		fctl.WithAliases("sm", "set-meta"),
		fctl.WithConfirmFlag(),
		fctl.WithValidArgs("last"),
		fctl.WithArgs(cobra.MinimumNArgs(2)),
		fctl.WithController[*SetMetadataStore](NewSetMetadataController()),
	)
}

func (c *SetMetadataController) GetStore() *SetMetadataStore {
	return c.store
}

func (c *SetMetadataController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := store.GetStore(cmd.Context())

	metadata, err := fctl.ParseMetadata(args[1:])
	if err != nil {
		return nil, err
	}

	transactionID, err := internal.TransactionIDOrLastN(cmd.Context(), store.Client(),
		fctl.GetString(cmd, internal.LedgerFlag), args[0])
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to set a metadata on transaction %d", transactionID) {
		return nil, fctl.ErrMissingApproval
	}

	request := operations.AddMetadataOnTransactionRequest{
		Ledger:      fctl.GetString(cmd, internal.LedgerFlag),
		Txid:        transactionID,
		RequestBody: collectionutils.ConvertMap(metadata, collectionutils.ToAny[string]),
	}
	response, err := store.Client().Ledger.AddMetadataOnTransaction(cmd.Context(), request)
	if err != nil {
		return nil, err
	}

	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = response.StatusCode == 204
	return c, nil
}

// TODO: This need to use the ui.NewListModel
func (c *SetMetadataController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Metadata added!")
	return nil
}
