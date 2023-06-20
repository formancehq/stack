package wallets

import (
	"fmt"
	"strconv"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewDebitWalletCommand() *cobra.Command {
	const (
		pendingFlag     = "pending"
		metadataFlag    = "metadata"
		descriptionFlag = "description"
		balanceFlag     = "balance"
		destinationFlag = "destination"
	)
	return fctl.NewCommand("debit <amount> <asset>",
		fctl.WithShortDescription("Debit a wallet"),
		fctl.WithAliases("deb"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.RangeArgs(2, 3)),
		fctl.WithStringFlag(descriptionFlag, "", "Debit description"),
		fctl.WithBoolFlag(pendingFlag, false, "Create a pending debit"),
		fctl.WithStringSliceFlag(metadataFlag, []string{""}, "Metadata to use"),
		fctl.WithStringSliceFlag(balanceFlag, []string{""}, "Balance to debit"),
		fctl.WithStringFlag(destinationFlag, "",
			`Use --destination account=<account> | --destination wallet=id:<wallet-id>[/<balance>] | --destination wallet=name:<wallet-name>[/<balance>]`),
		internal.WithTargetingWalletByName(),
		internal.WithTargetingWalletByID(),
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

			if !fctl.CheckStackApprobation(cmd, stack, "You are about to debit a wallets") {
				return fctl.ErrMissingApproval
			}

			client, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return errors.Wrap(err, "creating stack client")
			}

			pending := fctl.GetBool(cmd, pendingFlag)

			metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, metadataFlag))
			if err != nil {
				return err
			}

			amountStr := args[0]
			asset := args[1]
			walletID, err := internal.RequireWalletID(cmd, client)
			if err != nil {
				return err
			}

			description := fctl.GetString(cmd, descriptionFlag)

			amount, err := strconv.ParseInt(amountStr, 10, 32)
			if err != nil {
				return errors.Wrap(err, "parsing amount")
			}

			var destination *shared.Subject
			if destinationStr := fctl.GetString(cmd, destinationFlag); destinationStr != "" {
				destination, err = internal.ParseSubject(destinationStr, cmd, client)
				if err != nil {
					return err
				}
			}

			response, err := client.Wallets.DebitWallet(cmd.Context(), operations.DebitWalletRequest{
				DebitWalletRequest: &shared.DebitWalletRequest{
					Amount: shared.Monetary{
						Asset:  asset,
						Amount: amount,
					},
					Pending:     &pending,
					Metadata:    metadata,
					Description: &description,
					Destination: destination,
					Balances:    fctl.GetStringSlice(cmd, balanceFlag),
				},
				ID: walletID,
			})
			if err != nil {
				return errors.Wrap(err, "debiting wallet")
			}

			if response.WalletsErrorResponse != nil {
				return fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
			}

			if response.StatusCode >= 300 {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}

			if response != nil && response.DebitWalletResponse.Data.ID != "" {
				pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Wallet debited successfully with hold id '%s'!", response.DebitWalletResponse.Data.ID)
			} else {
				pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Wallet debited successfully!")
			}

			return nil
		}),
	)
}
