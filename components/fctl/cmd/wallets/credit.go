package wallets

import (
	"fmt"
	"strconv"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	cobra "github.com/spf13/cobra"
)

func NewCreditWalletCommand() *cobra.Command {
	const (
		metadataFlag = "metadata"
		balanceFlag  = "balance"
		sourceFlag   = "source"
	)
	return fctl.NewCommand("credit <amount> <asset>",
		fctl.WithShortDescription("Credit a wallets"),
		fctl.WithAliases("cr"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithStringSliceFlag(metadataFlag, []string{""}, "Metadata to use"),
		fctl.WithStringFlag(balanceFlag, "", "Balance to credit"),
		fctl.WithStringSliceFlag(sourceFlag, []string{}, `Use --source account=<account> | --source wallet=id:<wallet-id>[/<balance>] | --source wallet=name:<wallet-name>[/<balance>]`),
		internal.WithTargetingWalletByName(),
		internal.WithTargetingWalletByID(),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return errors.Wrap(err, "reading config")
			}

			organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
			if err != nil {
				return err
			}

			stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
			if err != nil {
				return err
			}

			if !fctl.CheckStackApprobation(cmd, stack, "You are about to credit a wallets") {
				return fctl.ErrMissingApproval
			}

			client, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return errors.Wrap(err, "creating stack client")
			}

			amountStr := args[0]
			asset := args[1]
			walletID, err := internal.RetrieveWalletID(cmd, client)
			if err != nil {
				return err
			}

			if walletID == "" {
				return errors.New("You need to specify wallet id using --id or --name flags")
			}

			amount, err := strconv.ParseUint(amountStr, 10, 64)
			if err != nil {
				return errors.Wrap(err, "parsing amount")
			}

			metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, metadataFlag))
			if err != nil {
				return err
			}

			sources := make([]shared.Subject, 0)
			for _, sourceStr := range fctl.GetStringSlice(cmd, sourceFlag) {
				source, err := internal.ParseSubject(sourceStr, cmd, client)
				if err != nil {
					return err
				}
				sources = append(sources, *source)
			}

			request := operations.CreditWalletRequest{
				ID: walletID,
				CreditWalletRequest: &shared.CreditWalletRequest{
					Amount: shared.Monetary{
						Asset:  asset,
						Amount: int64(amount),
					},
					Metadata: metadata,
					Sources:  sources,
					Balance:  formance.String(fctl.GetString(cmd, balanceFlag)),
				},
			}
			response, err := client.Wallets.CreditWallet(cmd.Context(), request)
			if err != nil {
				return errors.Wrap(err, "crediting wallet")
			}

			if response.WalletsErrorResponse != nil {
				return fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
			}

			if response.StatusCode >= 300 {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}

			pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Wallet credited successfully!")

			return nil
		}),
	)
}
