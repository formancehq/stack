package transactions

import (
	"fmt"
	"strings"
	"time"

	internal "github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	const (
		listTransactionsPageSizeFlag   = "page-size"
		listTransactionsMetadataFlag   = "metadata"
		listTransactionsReferenceFlag  = "reference"
		listTransactionAccountFlag     = "account"
		listTransactionDestinationFlag = "dst"
		listTransactionSourceFlag      = "src"
		listTransactionsAfterFlag      = "after"
		listTransactionsEndTimeFlag    = "end"
		listTransactionsStartTimeFlag  = "start"
	)

	return fctl.NewCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List transactions"),
		fctl.WithStringFlag(listTransactionAccountFlag, "", "Filter on account"),
		fctl.WithStringFlag(listTransactionDestinationFlag, "", "Filter on destination account"),
		fctl.WithStringFlag(listTransactionsAfterFlag, "", "Filter results after given tx id"),
		fctl.WithStringFlag(listTransactionsEndTimeFlag, "", "Consider transactions before date"),
		fctl.WithStringFlag(listTransactionsStartTimeFlag, "", "Consider transactions after date"),
		fctl.WithStringFlag(listTransactionSourceFlag, "", "Filter on source account"),
		fctl.WithStringFlag(listTransactionsReferenceFlag, "", "Filter on reference"),
		fctl.WithStringSliceFlag(listTransactionsMetadataFlag, []string{}, "Filter transactions with metadata"),
		fctl.WithIntFlag(listTransactionsPageSizeFlag, 5, "Page size"),
		// SDK not generating correct requests
		fctl.WithHiddenFlag(listTransactionsMetadataFlag),
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

			ledgerClient, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return err
			}

			metadata := map[string]interface{}{}
			for _, v := range fctl.GetStringSlice(cmd, listTransactionsMetadataFlag) {
				parts := strings.SplitN(v, "=", 2)
				if len(parts) == 1 {
					return fmt.Errorf("malformed metadata: %s", v)
				}
				metadata[parts[0]] = parts[1]
			}

			ledger := fctl.GetString(cmd, internal.LedgerFlag)
			rsp, _, err := ledgerClient.TransactionsApi.
				ListTransactions(cmd.Context(), ledger).
				PageSize(int32(fctl.GetInt(cmd, listTransactionsPageSizeFlag))).
				Reference(fctl.GetString(cmd, listTransactionsReferenceFlag)).
				Account(fctl.GetString(cmd, listTransactionAccountFlag)).
				Destination(fctl.GetString(cmd, listTransactionDestinationFlag)).
				Source(fctl.GetString(cmd, listTransactionSourceFlag)).
				After(fctl.GetString(cmd, listTransactionsAfterFlag)).
				EndTime(fctl.GetString(cmd, listTransactionsEndTimeFlag)).
				StartTime(fctl.GetString(cmd, listTransactionsStartTimeFlag)).
				Metadata(metadata).
				Execute()
			if err != nil {
				return err
			}

			tableData := fctl.Map(rsp.Cursor.Data, func(tx formance.Transaction) []string {
				return []string{
					fmt.Sprintf("%d", tx.Txid),
					func() string {
						if tx.Reference == nil {
							return ""
						}
						return *tx.Reference
					}(),
					tx.Timestamp.Format(time.RFC3339),
				}
			})
			tableData = fctl.Prepend(tableData, []string{"ID", "Reference", "Date"})
			return pterm.DefaultTable.
				WithHasHeader().
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render()
		}),
	)
}
