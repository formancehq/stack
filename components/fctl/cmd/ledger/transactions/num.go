package transactions

import (
	"flag"
	"fmt"
	"math/big"
	"strings"
	"time"

	internal "github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/spf13/cobra"
)

const (
	amountVarFlag  = "amount-var"
	portionVarFlag = "portion-var"
	accountVarFlag = "account-var"
	timestampFlag  = "timestamp"
)

const (
	useNum         = "num -|<filename>"
	shortNum       = "Execute a numscript script on a ledger"
	descriptionNum = `More help on variables can be found here: https://docs.formance.com/oss/ledger/reference/numscript/variables`
)

type NumStore struct {
	Transaction *internal.ExportTransaction `json:"transaction"`
}

func NewNumStore() *NumStore {
	return &NumStore{}
}
func NewNumConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)
	flags.String(amountVarFlag, "", "Pass a variable of type 'amount'")
	flags.String(portionVarFlag, "", "Pass a variable of type 'portion'")
	flags.String(accountVarFlag, "", "Pass a variable of type 'account'")
	flags.String(timestampFlag, "", "Timestamp to use (format RFC3339)")
	flags.String(internal.ReferenceFlag, "", "Reference to add to the generated transaction")
	flags.String(internal.MetadataFlag, "", "Filter accounts with metadata") //  experimental feature: Should be hidden
	return fctl.NewControllerConfig(
		useNum,
		descriptionNum,
		shortNum,
		[]string{},
		flags,
		fctl.Organization, fctl.Stack, fctl.Ledger,
	)
}

var _ fctl.Controller[*NumStore] = (*NumController)(nil)

type NumController struct {
	store  *NumStore
	config *fctl.ControllerConfig
}

func NewNumController(config *fctl.ControllerConfig) *NumController {
	return &NumController{
		store:  NewNumStore(),
		config: config,
	}
}

func (c *NumController) GetStore() *NumStore {
	return c.store
}

func (c *NumController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *NumController) Run() (fctl.Renderable, error) {

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

	script, err := fctl.ReadFile(flags, stack, args[0])
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(flags, stack, "You are about to apply a numscript") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	vars := map[string]interface{}{}
	for _, v := range fctl.GetStringSlice(flags, accountVarFlag) {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 1 {
			return nil, fmt.Errorf("malformed var: %s", v)
		}
		vars[parts[0]] = parts[1]
	}
	for _, v := range fctl.GetStringSlice(flags, portionVarFlag) {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 1 {
			return nil, fmt.Errorf("malformed var: %s", v)
		}
		vars[parts[0]] = parts[1]
	}
	for _, v := range fctl.GetStringSlice(flags, amountVarFlag) {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 1 {
			return nil, fmt.Errorf("malformed var: %s", v)
		}

		amountParts := strings.SplitN(parts[1], "/", 2)
		if len(amountParts) != 2 {
			return nil, fmt.Errorf("malformed var: %s", v)
		}

		amount, ok := big.NewInt(0).SetString(amountParts[0], 10)
		if !ok {
			return nil, fmt.Errorf("unable to parse '%s' as big int", amountParts[0])
		}

		vars[parts[0]] = map[string]any{
			"amount": amount,
			"asset":  amountParts[1],
		}
	}

	timestampStr := fctl.GetString(flags, timestampFlag)
	var (
		timestamp time.Time
	)
	if timestampStr != "" {
		timestamp, err = time.Parse(time.RFC3339Nano, timestampStr)
		if err != nil {
			return nil, err
		}
	}

	reference := fctl.GetString(flags, internal.ReferenceFlag)

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(flags, fctl.MetadataFlag))
	if err != nil {
		return nil, err
	}

	ledger := fctl.GetString(flags, internal.LedgerFlag)

	tx, err := internal.CreateTransaction(client, ctx, operations.CreateTransactionRequest{
		PostTransaction: shared.PostTransaction{
			Metadata:  metadata,
			Reference: &reference,
			Script: &shared.PostTransactionScript{
				Plain: script,
				Vars:  vars,
			},
			Timestamp: func() *time.Time {
				if timestamp.IsZero() {
					return nil
				}
				return &timestamp
			}(),
		},
		Ledger: ledger,
	})
	if err != nil {
		return nil, err
	}

	c.store.Transaction = internal.NewExportTransaction(tx)

	return c, nil
}

func (c *NumController) Render() error {
	return internal.PrintTransaction(c.config.GetOut(), c.store.Transaction)
}

func NewNumCommand() *cobra.Command {
	c := NewNumConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*NumStore](NewNumController(c)),
	)
}
