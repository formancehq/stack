package transactions

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/spf13/cobra"
)

type ShowStore struct {
	Transaction internal.ExpandedTransaction `json:"transaction"`
}
type ShowController struct {
	store *ShowStore
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

func NewDefaultShowStore() *ShowStore {
	return &ShowStore{}
}

func NewShowController() *ShowController {
	return &ShowController{
		store: NewDefaultShowStore(),
	}
}

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("show <transaction-id>",
		fctl.WithShortDescription("Print a transaction"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithAliases("sh"),
		fctl.WithValidArgs("last"),
		fctl.WithController[*ShowStore](NewShowController()),
	)
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	ledgerClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	ledger := fctl.GetString(cmd, internal.LedgerFlag)
	txId, err := internal.TransactionIDOrLastN(cmd.Context(), ledgerClient, ledger, args[0])
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(cmd, cfg)
	baseUrl := profile.ServicesBaseUrl(stack).String()

	response, err := internal.GetTransaction(
		ledgerClient,
		cmd.Context(),
		baseUrl,
		operations.GetTransactionRequest{
			Ledger: ledger,
			Txid:   txId,
		},
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

	c.store.Transaction = response.GetTransactionResponse.Data

	return c, nil
}

func (c *ShowController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintExpandedTransaction(cmd.OutOrStdout(), c.store.Transaction)
}
