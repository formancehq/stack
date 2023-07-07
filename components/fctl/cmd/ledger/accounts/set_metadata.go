package accounts

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
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
	return &SetMetadataStore{
		Success: false,
	}
}

func NewSetMetadataController() *SetMetadataController {
	return &SetMetadataController{
		store: NewDefaultSetMetadataStore(),
	}
}

func NewSetMetadataCommand() *cobra.Command {
	return fctl.NewCommand("set-metadata <address> [<key>=<value>...]",
		fctl.WithConfirmFlag(),
		fctl.WithShortDescription("Set metadata on address"),
		fctl.WithAliases("sm", "set-meta"),
		fctl.WithArgs(cobra.MinimumNArgs(2)),
		fctl.WithController[*SetMetadataStore](NewSetMetadataController()),
	)
}

func (c *SetMetadataController) GetStore() *SetMetadataStore {
	return c.store
}

func (c *SetMetadataController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	metadata, err := fctl.ParseMetadata(args[1:])
	if err != nil {
		return nil, err
	}

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

	address := args[0]

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to set a metadata on address '%s'", address) {
		return nil, fctl.ErrMissingApproval
	}

	ledgerClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	request := operations.AddMetadataToAccountRequest{
		Ledger:      fctl.GetString(cmd, internal.LedgerFlag),
		Address:     address,
		RequestBody: metadata,
	}
	response, err := ledgerClient.Ledger.AddMetadataToAccount(cmd.Context(), request)
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

func (c *SetMetadataController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Metadata added!")
	return nil
}
