package transactions

import (
	"flag"
	"fmt"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useSetMetadata   = "set-metadata <transaction-id> [<key>=<value>...]"
	shortSetMetadata = "Set metadata on transaction"
)

type SetMetadataStore struct {
	Success bool `json:"success"`
}

func NewSetMetadataStore() *SetMetadataStore {
	return &SetMetadataStore{}
}
func NewSetMetadataConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useSetMetadata, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)

	return fctl.NewControllerConfig(
		useSetMetadata,
		shortSetMetadata,
		shortSetMetadata,
		[]string{
			"sm", "set-meta",
		},
		flags, fctl.Organization, fctl.Stack, fctl.Ledger,
	)
}

var _ fctl.Controller[*SetMetadataStore] = (*SetMetadataController)(nil)

type SetMetadataController struct {
	store  *SetMetadataStore
	config *fctl.ControllerConfig
}

func NewSetMetadataController(config *fctl.ControllerConfig) *SetMetadataController {
	return &SetMetadataController{
		store:  NewSetMetadataStore(),
		config: config,
	}
}

func (c *SetMetadataController) GetStore() *SetMetadataStore {
	return c.store
}

func (c *SetMetadataController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *SetMetadataController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()
	out := c.config.GetOut()

	metadata, err := fctl.ParseMetadata(args[1:])
	if err != nil {
		return nil, err
	}

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	ledgerClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	transactionID, err := internal.TransactionIDOrLastN(ctx, ledgerClient,
		fctl.GetString(flags, internal.LedgerFlag), args[0])
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(flags, stack, "You are about to set a metadata on transaction %d", transactionID) {
		return nil, fctl.ErrMissingApproval
	}

	request := operations.AddMetadataOnTransactionRequest{
		Ledger:      fctl.GetString(flags, internal.LedgerFlag),
		Txid:        transactionID,
		RequestBody: metadata,
	}
	response, err := ledgerClient.Ledger.AddMetadataOnTransaction(ctx, request)
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

func (c *SetMetadataController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Metadata added!")
	return nil
}

func NewSetMetadataCommand() *cobra.Command {

	config := NewSetMetadataConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithValidArgs("last"),
		fctl.WithArgs(cobra.MinimumNArgs(2)),
		fctl.WithController[*SetMetadataStore](NewSetMetadataController(config)),
	)
}
