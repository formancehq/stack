package transactions

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	const (
		amountVarFlag  = "amount-var"
		portionVarFlag = "portion-var"
		accountVarFlag = "account-var"
		metadataFlag   = "metadata"
		referenceFlag  = "reference"
		timestampFlag  = "timestamp"
	)
	return fctl.NewCommand("num -|<filename>",
		fctl.WithShortDescription("Execute a numscript script on a ledger"),
		fctl.WithDescription(`More help on variables can be found here: https://docs.formance.com/oss/ledger/reference/numscript/variables`),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithConfirmFlag(),
		fctl.WithStringSliceFlag(amountVarFlag, []string{""}, "Pass a variable of type 'amount'"),
		fctl.WithStringSliceFlag(portionVarFlag, []string{""}, "Pass a variable of type 'portion'"),
		fctl.WithStringSliceFlag(accountVarFlag, []string{""}, "Pass a variable of type 'account'"),
		fctl.WithStringSliceFlag(metadataFlag, []string{""}, "Metadata to use"),
		fctl.WithStringFlag(timestampFlag, "", "Timestamp to use (format RFC3339)"),
		fctl.WithStringFlag(referenceFlag, "", "Reference to add to the generated transaction"),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}

			organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
			if err != nil {
				return err
			}

			stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
			if err != nil {
				return err
			}

			script, err := fctl.ReadFile(cmd, stack, args[0])
			if err != nil {
				return err
			}

			if !fctl.CheckStackApprobation(cmd, stack, "You are about to apply a numscript") {
				return fctl.ErrMissingApproval
			}

			client, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return err
			}

			vars := map[string]interface{}{}
			for _, v := range fctl.GetStringSlice(cmd, accountVarFlag) {
				parts := strings.SplitN(v, "=", 2)
				if len(parts) == 1 {
					return fmt.Errorf("malformed var: %s", v)
				}
				vars[parts[0]] = parts[1]
			}
			for _, v := range fctl.GetStringSlice(cmd, portionVarFlag) {
				parts := strings.SplitN(v, "=", 2)
				if len(parts) == 1 {
					return fmt.Errorf("malformed var: %s", v)
				}
				vars[parts[0]] = parts[1]
			}
			for _, v := range fctl.GetStringSlice(cmd, amountVarFlag) {
				parts := strings.SplitN(v, "=", 2)
				if len(parts) == 1 {
					return fmt.Errorf("malformed var: %s", v)
				}

				amountParts := strings.SplitN(parts[1], "/", 2)
				if len(amountParts) != 2 {
					return fmt.Errorf("malformed var: %s", v)
				}

				amount, err := strconv.ParseInt(amountParts[0], 10, 64)
				if err != nil {
					return fmt.Errorf("malformed var: %s", v)
				}

				vars[parts[0]] = map[string]any{
					"amount": amount,
					"asset":  amountParts[1],
				}
			}

			timestampStr := fctl.GetString(cmd, timestampFlag)
			var (
				timestamp time.Time
			)
			if timestampStr != "" {
				timestamp, err = time.Parse(time.RFC3339Nano, timestampStr)
				if err != nil {
					return err
				}
			}

			reference := fctl.GetString(cmd, referenceFlag)

			metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, metadataFlag))
			if err != nil {
				return err
			}

			ledger := fctl.GetString(cmd, internal.LedgerFlag)

			tx, err := internal.CreateTransaction(client, cmd.Context(), ledger, operations.CreateTransactionRequest{
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
				return err
			}

			return internal.PrintTransaction(cmd.OutOrStdout(), *tx)
		}),
	)
}
