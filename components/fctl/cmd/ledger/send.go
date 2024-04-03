package ledger

import (
	"fmt"
	"math/big"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/spf13/cobra"
)

type SendStore struct {
	Transaction *shared.Transaction `json:"transaction"`
}
type SendController struct {
	store         *SendStore
	metadataFlag  string
	referenceFlag string
}

var _ fctl.Controller[*SendStore] = (*SendController)(nil)

func NewDefaultSendStore() *SendStore {
	return &SendStore{}
}

func NewSendController() *SendController {
	return &SendController{
		store:         NewDefaultSendStore(),
		metadataFlag:  "metadata",
		referenceFlag: "reference",
	}
}

func NewSendCommand() *cobra.Command {
	c := NewSendController()
	return fctl.NewCommand("send [<source>] <destination> <amount> <asset>",
		fctl.WithAliases("s", "se"),
		fctl.WithShortDescription("Send from one account to another"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.RangeArgs(3, 4)),
		fctl.WithStringSliceFlag(c.metadataFlag, []string{""}, "Metadata to use"),
		fctl.WithStringFlag(c.referenceFlag, "", "Reference to add to the generated transaction"),
		fctl.WithController[*SendStore](c),
	)
}

func (c *SendController) GetStore() *SendStore {
	return c.store
}

func (c *SendController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to create a new transaction") {
		return nil, fctl.ErrMissingApproval
	}

	var source, destination, asset, amountStr string
	if len(args) == 3 {
		source = "world"
		destination = args[0]
		amountStr = args[1]
		asset = args[2]
	} else {
		source = args[0]
		destination = args[1]
		amountStr = args[2]
		asset = args[3]
	}

	amount, ok := big.NewInt(0).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse '%s' as big int", amountStr)
	}

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, c.metadataFlag))
	if err != nil {
		return nil, err
	}

	reference := fctl.GetString(cmd, c.referenceFlag)

	response, err := store.Client().Ledger.CreateTransaction(cmd.Context(), operations.CreateTransactionRequest{
		PostTransaction: shared.PostTransaction{
			Metadata: collectionutils.ConvertMap(metadata, collectionutils.ToAny[string]),
			Postings: []shared.Posting{
				{
					Amount:      amount,
					Asset:       asset,
					Destination: destination,
					Source:      source,
				},
			},
			Reference: &reference,
		},
		Ledger: fctl.GetString(cmd, internal.LedgerFlag),
	})
	if err != nil {
		return nil, err
	}

	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code %d when creating transaction", response.StatusCode)
	}

	c.store.Transaction = &response.TransactionsResponse.Data[0]
	return c, nil
}

// TODO: This need to use the ui.NewListModel
func (c *SendController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintTransaction(cmd.OutOrStdout(), internal.WrapV1Transaction(*c.store.Transaction))
}
