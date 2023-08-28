package ledger

import (
	"flag"
	"fmt"
	"math/big"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/spf13/cobra"
)

const (
	referenceFlag = "reference"
	useSend       = "send [<source>] <destination> <amount> <asset>"
	shortSend     = "Send from one account to another"
	descriptionSend
)

type SendStore struct {
	Transaction *internal.ExportTransaction `json:"transaction"`
}

func NewSendStore() *SendStore {
	return &SendStore{}
}
func NewSendConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useSend, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)
	fctl.WithMetadataFlag(flags)
	flags.String(referenceFlag, "", "Reference to add to the generated transaction")

	return fctl.NewControllerConfig(
		useSend,
		descriptionSend,
		shortSend,
		[]string{
			"s", "se",
		},
		flags,
		fctl.Organization, fctl.Stack, fctl.Ledger,
	)
}

type SendController struct {
	store  *SendStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*SendStore] = (*SendController)(nil)

func NewSendController(config *fctl.ControllerConfig) *SendController {
	return &SendController{
		store:  NewSendStore(),
		config: config,
	}
}

func (c *SendController) GetStore() *SendStore {
	return c.store
}

func (c *SendController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *SendController) Run() (fctl.Renderable, error) {
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

	if !fctl.CheckStackApprobation(flags, stack, "You are about to create a new transaction") {
		return nil, fctl.ErrMissingApproval
	}

	ledgerClient, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
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

	amount, ok := big.NewInt(0).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse '%s' as big int", amountStr)
	}

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(flags, fctl.MetadataFlag))
	if err != nil {
		return nil, err
	}

	reference := fctl.GetString(flags, referenceFlag)

	tx, err := internal.CreateTransaction(ledgerClient, ctx, operations.CreateTransactionRequest{
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
		Ledger: fctl.GetString(flags, internal.LedgerFlag),
	})
	if err != nil {
		return nil, err
	}
	c.store.Transaction = internal.NewExportTransaction(tx)
	return c, nil
}

func (c *SendController) Render() error {
	return internal.PrintTransaction(c.config.GetOut(), c.store.Transaction)
}

func NewSendCommand() *cobra.Command {
	c := NewSendConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.RangeArgs(3, 4)),
		fctl.WithController[*SendStore](NewSendController(c)),
	)
}
