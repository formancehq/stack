package payments

import (
	"fmt"
	"time"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewListPaymentsCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithAliases("ls"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithShortDescription("List payments"),
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

			client, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return err
			}

			paymentsCursor, _, err := client.PaymentsApi.ListPayments(cmd.Context()).Execute()
			if err != nil {
				return fctl.WrapError(err, "listing payments")
			}

			tableData := fctl.Map(paymentsCursor.Cursor.Data, func(payment formance.Payment) []string {
				return []string{
					payment.Id,
					payment.Type,
					fmt.Sprint(payment.InitialAmount),
					payment.Asset,
					string(payment.Status),
					payment.Scheme,
					payment.Reference,
					payment.AccountID,
					string(payment.Provider),
					payment.CreatedAt.Format(time.RFC3339),
				}
			})
			tableData = fctl.Prepend(tableData, []string{"ID", "Type", "Amount", "Asset", "Status",
				"Scheme", "Reference", "Account ID", "Provider", "Created at"})
			return pterm.DefaultTable.
				WithHasHeader().
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render()
		}),
	)
}
