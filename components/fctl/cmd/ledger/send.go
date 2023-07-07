package ledger

import (
	"strconv"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/spf13/cobra"
)

type SendStore struct {
	Transaction *internal.Transaction `json:"transaction"`
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

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to create a new transaction") {
		return nil, fctl.ErrMissingApproval
	}

	ledgerClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
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

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		return nil, err
	}

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, c.metadataFlag))
	if err != nil {
		return nil, err
	}

	reference := fctl.GetString(cmd, c.referenceFlag)

	tx, err := internal.CreateTransaction(ledgerClient, cmd.Context(), operations.CreateTransactionRequest{
		PostTransaction: shared.PostTransaction{
			Metadata: metadata,
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
	c.store.Transaction = tx
	return c, nil
}

// TODO: This need to use the ui.NewListModel
func (c *SendController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintTransaction(cmd.OutOrStdout(), *c.store.Transaction)
}
