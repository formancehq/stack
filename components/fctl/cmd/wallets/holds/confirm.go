package holds

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewConfirmCommand() *cobra.Command {
	const (
		finalFlag  = "final"
		amountFlag = "amount"
	)
	return fctl.NewCommand("confirm <hold-id>",
		fctl.WithShortDescription("Confirm a hold"),
		fctl.WithAliases("c", "conf"),
		fctl.WithArgs(cobra.RangeArgs(1, 2)),
		fctl.WithBoolFlag(finalFlag, false, "Is final debit (close hold)"),
		fctl.WithIntFlag(amountFlag, 0, "Amount to confirm"),
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

			stackClient, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return errors.Wrap(err, "creating stack client")
			}

			final := fctl.GetBool(cmd, finalFlag)
			amount := int64(fctl.GetInt(cmd, amountFlag))

			request := operations.ConfirmHoldRequest{
				HoldID: args[0],
				ConfirmHoldRequest: &shared.ConfirmHoldRequest{
					Amount: &amount,
					Final:  &final,
				},
			}
			response, err := stackClient.Wallets.ConfirmHold(cmd.Context(), request)
			if err != nil {
				return errors.Wrap(err, "confirming hold")
			}

			if response.WalletsErrorResponse != nil {
				return fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
			}

			if response.StatusCode >= 300 {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}

			pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Hold '%s' confirmed!", args[0])

			return nil
		}),
	)
}
