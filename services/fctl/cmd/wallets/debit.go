package wallets

import (
	"strconv"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewDebitWalletCommand() *cobra.Command {
	const (
		pendingFlag     = "pending"
		metadataFlag    = "metadata"
		descriptionFlag = "description"
	)
	return fctl.NewCommand("debit <amount> <asset>",
		fctl.WithShortDescription("Debit a wallet"),
		fctl.WithAliases("deb"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.RangeArgs(2, 3)),
		fctl.WithStringFlag(descriptionFlag, "", "Debit description"),
		fctl.WithBoolFlag(pendingFlag, false, "Create a pending debit"),
		fctl.WithStringSliceFlag(metadataFlag, []string{""}, "Metadata to use"),
		internal.WithTargetingWalletByName(),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return errors.Wrap(err, "fctl.GetConfig")
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
			walletID, err := internal.RetrieveWalletIDFromName(cmd, client)
			if err != nil {
				return err
			}

			if walletID == "" {
				return errors.New("You need to specify wallet id using --id or --name flags")
			}

			description := fctl.GetString(cmd, descriptionFlag)

			amount, err := strconv.ParseInt(amountStr, 10, 32)
			if err != nil {
				return errors.Wrap(err, "parsing amount")
			}

			hold, _, err := client.WalletsApi.DebitWallet(cmd.Context(), walletID).DebitWalletRequest(formance.DebitWalletRequest{
				Amount: formance.Monetary{
					Asset:  asset,
					Amount: amount,
				},
				Pending:     &pending,
				Metadata:    metadata,
				Description: &description,
			}).Execute()
			if err != nil {
				return errors.Wrap(err, "Debiting wallets")
			}

			if hold != nil && hold.Data.Id != "" {
				fctl.Success(cmd.OutOrStdout(), "Wallet debited successfully with hold id '%s'!", hold.Data.Id)
			} else {
				fctl.Success(cmd.OutOrStdout(), "Wallet debited successfully!")
			}

			return nil
		}),
	)
}
