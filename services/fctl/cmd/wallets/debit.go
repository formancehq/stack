package wallets

import (
	"strconv"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewDebitWalletCommand() *cobra.Command {
	const (
		pendingFlag = "pending"
	)
	return fctl.NewCommand("debit ID AMOUNT ASSET",
		fctl.WithShortDescription("Debit a wallets"),
		fctl.WithAliases("deb"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(3)),
		fctl.WithBoolFlag(pendingFlag, false, "Create a pending debit"),
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

			amount, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return errors.Wrap(err, "parsing amount")
			}

			pending := fctl.GetBool(cmd, pendingFlag)
			hold, _, err := client.WalletsApi.DebitWallet(cmd.Context(), args[0]).DebitWalletRequest(formance.DebitWalletRequest{
				Amount: formance.Monetary{
					Asset:  args[2],
					Amount: float32(amount),
				},
				Pending: &pending,
			}).Execute()
			if err != nil {
				return errors.Wrap(err, "Debiting wallets")
			}

			if hold.Data.Id != "" {
				fctl.Success(cmd.OutOrStdout(), "Wallet debited successfully with hold id '%s'!", hold.Data.Id)
			} else {
				fctl.Success(cmd.OutOrStdout(), "Wallet debited successfully!")
			}

			return nil
		}),
	)
}
