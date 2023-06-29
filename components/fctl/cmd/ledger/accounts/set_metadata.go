package accounts

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/ledger/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewSetMetadataCommand() *cobra.Command {
	return fctl.NewCommand("set-metadata <address> [<key>=<value>...]",
		fctl.WithConfirmFlag(),
		fctl.WithShortDescription("Set metadata on address"),
		fctl.WithAliases("sm", "set-meta"),
		fctl.WithArgs(cobra.MinimumNArgs(2)),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {

			metadata, err := fctl.ParseMetadata(args[1:])
			if err != nil {
				return err
			}

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

			address := args[0]

			if !fctl.CheckStackApprobation(cmd, stack, "You are about to set a metadata on address '%s'", address) {
				return fctl.ErrMissingApproval
			}

			ledgerClient, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return err
			}

			request := operations.AddMetadataToAccountRequest{
				Ledger:      fctl.GetString(cmd, internal.LedgerFlag),
				Address:     address,
				RequestBody: metadata,
			}
			response, err := ledgerClient.Ledger.AddMetadataToAccount(cmd.Context(), request)
			if err != nil {
				return err
			}

			if response.ErrorResponse != nil {
				return fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
			}

			if response.StatusCode >= 300 {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}

			pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Metadata added!")
			return nil
		}),
	)
}
