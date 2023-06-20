package accounts

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("show <address>",
		fctl.WithShortDescription("Show account"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithAliases("sh", "s"),
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

			ledger := fctl.GetString(cmd, internal.LedgerFlag)
			response, err := ledgerClient.Ledger.GetAccount(cmd.Context(), operations.GetAccountRequest{
				Address: args[0],
				Ledger:  ledger,
			})
			if err != nil {
				return err
			}

			if response.ErrorResponse != nil {
				return fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
			}

			if response.StatusCode >= 300 {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}

			fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
			if response.AccountResponse.Data.Volumes != nil && len(response.AccountResponse.Data.Volumes) > 0 {
				tableData := pterm.TableData{}
				tableData = append(tableData, []string{"Asset", "Input", "Output"})
				for asset, volumes := range response.AccountResponse.Data.Volumes {
					input := volumes["input"]
					output := volumes["output"]
					tableData = append(tableData, []string{pterm.LightCyan(asset), fmt.Sprint(input), fmt.Sprint(output)})
				}
				if err := pterm.DefaultTable.
					WithHasHeader(true).
					WithWriter(cmd.OutOrStdout()).
					WithData(tableData).
					Render(); err != nil {
					return err
				}
			} else {
				fctl.Println("No balances.")
			}

			fmt.Fprintln(cmd.OutOrStdout())

			if err := fctl.PrintMetadata(cmd.OutOrStdout(), response.AccountResponse.Data.Metadata); err != nil {
				return err
			}

			return nil
		}),
	)
}
