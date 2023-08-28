package transactions

import (
	"flag"
	"fmt"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/spf13/cobra"
)

const (
	useRevert   = "revert <transaction-id>"
	shortRevert = "Revert transaction"
)

type RevertStore struct {
	Transaction *internal.ExportTransaction `json:"transaction"`
}

func NewRevertStore() *RevertStore {
	return &RevertStore{}
}
func NewRevertConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useRevert, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)
	return fctl.NewControllerConfig(
		useRevert,
		shortRevert,
		shortRevert,
		[]string{
			"rev",
		},
		flags,
		fctl.Organization, fctl.Stack, fctl.Ledger,
	)
}

type RevertController struct {
	store  *RevertStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*RevertStore] = (*RevertController)(nil)

func NewRevertController(config *fctl.ControllerConfig) *RevertController {
	return &RevertController{
		store:  NewRevertStore(),
		config: config,
	}
}

func (c *RevertController) GetStore() *RevertStore {
	return c.store
}

func (c *RevertController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *RevertController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()
	out := c.config.GetOut()

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

	if !fctl.CheckStackApprobation(flags, stack, "You are about to revert transaction %s", args[0]) {
		return nil, fctl.ErrMissingApproval
	}

	ledgerClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	ledger := fctl.GetString(flags, internal.LedgerFlag)
	txId, err := internal.TransactionIDOrLastN(ctx, ledgerClient, ledger, args[0])
	if err != nil {
		return nil, err
	}

	request := operations.RevertTransactionRequest{
		Ledger: ledger,
		Txid:   txId,
	}

	profile := fctl.GetCurrentProfile(flags, cfg)
	baseUrl := profile.ServicesBaseUrl(stack).String()

	response, err := internal.RevertTransaction(
		ledgerClient,
		ctx,
		baseUrl,
		request,
	)
	if err != nil {
		return nil, err
	}

	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Transaction = internal.NewExportTransaction(&response.RevertTransactionResponse.Data)
	return c, nil
}

func (c *RevertController) Render() error {
	return internal.PrintTransaction(c.config.GetOut(), c.store.Transaction)
}
func NewRevertCommand() *cobra.Command {

	config := NewRevertConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithValidArgs("last"),
		fctl.WithController[*RevertStore](NewRevertController(config)),
	)
}
