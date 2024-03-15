package accounts

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
	return fctl.NewCommand("delete-metadata <address> [<key>...]",
		fctl.WithShortDescription("Delete metadata on account (Start from ledger v2 api)"),
		fctl.WithAliases("dm", "del-meta"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.MinimumNArgs(2)),
		fctl.WithController[*DeleteMetadataStore](NewDeleteMetadataController()),
	)
}

func (c *DeleteMetadataController) GetStore() *DeleteMetadataStore {
	return c.store
}

func (c *DeleteMetadataController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := fctl.GetStackStore(cmd.Context())

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to set a metadata on account %s", args[0]) {
		return nil, fctl.ErrMissingApproval
	}

	response, err := store.Client().Ledger.V2DeleteAccountMetadata(cmd.Context(), operations.V2DeleteAccountMetadataRequest{
		Address: args[0],
		Key:     args[1],
		Ledger:  fctl.GetString(cmd, internal.LedgerFlag),
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
