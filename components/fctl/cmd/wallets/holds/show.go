package holds

import (
	"fmt"
	"io"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("show <hold-id>",
		fctl.WithShortDescription("Show a hold"),
		fctl.WithAliases("sh"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return errors.Wrap(err, "retrieving config")
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
				return errors.Wrap(err, "creating stack client")
			}

			request := operations.GetHoldRequest{
				HoldID: args[0],
			}
			response, err := client.Wallets.GetHold(cmd.Context(), request)
			if err != nil {
				return errors.Wrap(err, "getting hold")
			}

			if response.WalletsErrorResponse != nil {
				return fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
			}

			if response.StatusCode >= 300 {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}

			return PrintHold(cmd.OutOrStdout(), response.GetHoldResponse.Data)
		}),
	)
}

func PrintHold(out io.Writer, hold shared.ExpandedDebitHold) error {
	fctl.Section.Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), fmt.Sprint(hold.ID)})
	tableData = append(tableData, []string{pterm.LightCyan("Wallet ID"), hold.WalletID})
	tableData = append(tableData, []string{pterm.LightCyan("Original amount"), fmt.Sprint(hold.OriginalAmount)})
	tableData = append(tableData, []string{pterm.LightCyan("Remaining"), fmt.Sprint(hold.Remaining)})

	if err := pterm.DefaultTable.
		WithWriter(out).
		WithData(tableData).
		Render(); err != nil {
		return err
	}

	return fctl.PrintMetadata(out, hold.Metadata)
}
