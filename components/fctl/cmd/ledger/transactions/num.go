package transactions

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/formancehq/go-libs/collectionutils"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/spf13/cobra"
)

type NumStore struct {
	Transaction *shared.Transaction `json:"transaction"`
}
type NumController struct {
	store          *NumStore
	amountVarFlag  string
	portionVarFlag string
	accountVarFlag string
	metadataFlag   string
	referenceFlag  string
	timestampFlag  string
}

var _ fctl.Controller[*NumStore] = (*NumController)(nil)

func NewDefaultNumStore() *NumStore {
	return &NumStore{}
}

func NewNumController() *NumController {
	return &NumController{
		store:          NewDefaultNumStore(),
		amountVarFlag:  "amount-var",
		portionVarFlag: "portion-var",
		accountVarFlag: "account-var",
		metadataFlag:   "metadata",
		referenceFlag:  "reference",
		timestampFlag:  "timestamp",
	}
}

func NewNumCommand() *cobra.Command {
	c := NewNumController()

	return fctl.NewCommand("num -|<filename>",
		fctl.WithShortDescription("Execute a numscript script on a ledger"),
		fctl.WithDescription(`More help on variables can be found here: https://docs.formance.com/oss/ledger/reference/numscript/variables`),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithConfirmFlag(),
		fctl.WithStringSliceFlag(c.amountVarFlag, []string{""}, "Pass a variable of type 'amount'"),
		fctl.WithStringSliceFlag(c.portionVarFlag, []string{""}, "Pass a variable of type 'portion'"),
		fctl.WithStringSliceFlag(c.accountVarFlag, []string{""}, "Pass a variable of type 'account'"),
		fctl.WithStringSliceFlag(c.metadataFlag, []string{""}, "Metadata to use"),
		fctl.WithStringFlag(c.timestampFlag, "", "Timestamp to use (format RFC3339)"),
		fctl.WithStringFlag(c.referenceFlag, "", "Reference to add to the generated transaction"),
		fctl.WithController[*NumStore](c),
	)
}

func (c *NumController) GetStore() *NumStore {
	return c.store
}

func (c *NumController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := fctl.GetStackStore(cmd.Context())

	script, err := fctl.ReadFile(cmd, store.Stack(), args[0])
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to apply a numscript") {
		return nil, fctl.ErrMissingApproval
	}

	vars := map[string]interface{}{}
	for _, v := range fctl.GetStringSlice(cmd, c.accountVarFlag) {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 1 {
			return nil, fmt.Errorf("malformed var: %s", v)
		}
		vars[parts[0]] = parts[1]
	}
	for _, v := range fctl.GetStringSlice(cmd, c.portionVarFlag) {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 1 {
			return nil, fmt.Errorf("malformed var: %s", v)
		}
		vars[parts[0]] = parts[1]
	}
	for _, v := range fctl.GetStringSlice(cmd, c.amountVarFlag) {
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

	timestampStr := fctl.GetString(cmd, c.timestampFlag)
	var (
		timestamp time.Time
	)
	if timestampStr != "" {
		timestamp, err = time.Parse(time.RFC3339Nano, timestampStr)
		if err != nil {
			return nil, err
		}
	}

	reference := fctl.GetString(cmd, c.referenceFlag)

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, c.metadataFlag))
	if err != nil {
		return nil, err
	}

	ledger := fctl.GetString(cmd, internal.LedgerFlag)

	response, err := store.Client().Ledger.V1.CreateTransaction(cmd.Context(), operations.CreateTransactionRequest{
		PostTransaction: shared.PostTransaction{
			Metadata:  collectionutils.ConvertMap(metadata, collectionutils.ToAny[string]),
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

	c.store.Transaction = &response.TransactionsResponse.Data[0]

	return c, nil
}

func (c *NumController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintTransaction(cmd.OutOrStdout(), internal.WrapV1Transaction(*c.store.Transaction))
}
