package holds

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewConfirmCommand() *cobra.Command {
	return fctl.NewCommand("confirm WALLET_ID HOLD_ID",
		fctl.WithShortDescription("Confirm a hold"),
		fctl.WithAliases("c", "conf"),
		fctl.WithArgs(cobra.ExactArgs(2)),
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

			stackClient, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return errors.Wrap(err, "creating stack client")
			}

			_, err = stackClient.WalletsApi.ConfirmHold(cmd.Context(), args[0], args[1]).Execute()
			if err != nil {
				return errors.Wrap(err, "listing wallets")
			}

			fctl.Success(cmd.OutOrStdout(), "Hold '%s' confirmed!", args[1])

			return nil
		}),
	)
}
