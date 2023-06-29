package ledger

import (
	"fmt"

	internal "github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewBalancesCommand() *cobra.Command {
	const (
		addressFlag = "address"
	)
	return fctl.NewCommand("balances",
		fctl.WithAliases("balance", "bal", "b"),
		fctl.WithStringFlag(addressFlag, "", "Filter on specific address"),
		fctl.WithShortDescription("Read balances"),
		fctl.WithArgs(cobra.ExactArgs(0)),
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

			response, err := client.Ledger.GetBalances(
				cmd.Context(),
				operations.GetBalancesRequest{
					Address: fctl.Ptr(fctl.GetString(cmd, addressFlag)),
					Ledger:  fctl.GetString(cmd, internal.LedgerFlag),
				},
			)
			if err != nil {
				return err
			}

			if response.ErrorResponse != nil {
				return fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
			}

			if response.StatusCode >= 300 {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}

			balances := response.BalancesCursorResponse

			tableData := pterm.TableData{}
			tableData = append(tableData, []string{"Account", "Asset", "Balance"})
			for _, accountBalances := range balances.Cursor.Data {
				for account, volumes := range accountBalances {
					for asset, balance := range volumes {
						tableData = append(tableData, []string{
							account, asset, fmt.Sprint(balance),
						})
					}
				}
			}
			if err := pterm.DefaultTable.
				WithHasHeader(true).
				WithWriter(cmd.OutOrStdout()).
				WithData(tableData).
				Render(); err != nil {
				return err
			}

			return nil
		}),
	)
}
