package payments

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

func NewSetMetadataStore() *SetMetadataStore {
	return &SetMetadataStore{}
}

func NewSetMetadataController() *SetMetadataController {
	return &SetMetadataController{
		store: NewSetMetadataStore(),
	}
}

func NewSetMetadataCommand() *cobra.Command {
	c := NewSetMetadataController()
	return fctl.NewCommand("set-metadata <paymentID> [<key>=<value>...]",
		fctl.WithConfirmFlag(),
		fctl.WithShortDescription("Set metadata on paymentID"),
		fctl.WithAliases("sm", "set-meta"),
		fctl.WithArgs(cobra.MinimumNArgs(2)),
		fctl.WithController[*SetMetadataStore](c),
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

	paymentID := args[0]

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to set a metadata on paymentID '%s'", paymentID) {
		return nil, fctl.ErrMissingApproval
	}

	request := operations.UpdateMetadataRequest{
		RequestBody: metadata,
		PaymentID:   paymentID,
	}

	response, err := store.Client().Payments.UpdateMetadata(cmd.Context(), request)
	if err != nil {
		return nil, err
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
