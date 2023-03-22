package holds

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewVoidCommand() *cobra.Command {
	return fctl.NewCommand("void <hold-id>",
		fctl.WithShortDescription("Void a hold"),
		fctl.WithAliases("v"),
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

			stackClient, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return errors.Wrap(err, "creating stack client")
			}

			request := operations.VoidHoldRequest{
				HoldID: args[0],
			}
			response, err := stackClient.Wallets.VoidHold(cmd.Context(), request)
			if err != nil {
				return errors.Wrap(err, "voiding hold")
			}

			if response.WalletsErrorResponse != nil {
				return fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
			}

			if response.StatusCode >= 300 {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}

			pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Hold '%s' voided!", args[0])

			return nil
		}),
	)
}
