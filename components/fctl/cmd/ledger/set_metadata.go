package ledger

import (
	"fmt"

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
	return fctl.NewCommand("set-metadata <ledger-name> [<key>=<value>...]",
		fctl.WithShortDescription("Set metadata on a ledger (Start from ledger v2 api)"),
		fctl.WithAliases("sm", "set-meta"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.MinimumNArgs(2)),
		fctl.WithController[*SetMetadataStore](NewSetMetadataController()),
	)
}

func (c *SetMetadataController) GetStore() *SetMetadataStore {
	return c.store
}

func (c *SetMetadataController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := fctl.GetStackStore(cmd.Context())

	metadata, err := fctl.ParseMetadata(args[1:])
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to set a metadata on ledger %s", args[0]) {
		return nil, fctl.ErrMissingApproval
	}

	request := operations.V2UpdateLedgerMetadataRequest{
		Ledger:      args[0],
		RequestBody: metadata,
	}
	response, err := store.Client().Ledger.V2UpdateLedgerMetadata(cmd.Context(), request)
	if err != nil {
		return nil, err
	}

	if response.V2ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.V2ErrorResponse.ErrorCode, response.V2ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = response.StatusCode == 204
	return c, nil
}

func (c *SetMetadataController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Metadata added!")
	return nil
}
