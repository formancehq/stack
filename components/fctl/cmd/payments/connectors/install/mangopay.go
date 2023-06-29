package install

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/connectors/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewMangoPayCommand() *cobra.Command {
	const (
		endpointFlag = "endpoint"

		defaultEndpoint = "https://api.sandbox.mangopay.com"
	)
	return fctl.NewCommand(internal.MangoPayConnector+" <clientID> <apiKey>",
		fctl.WithShortDescription("Install a MangoPay connector"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithStringFlag(endpointFlag, defaultEndpoint, "API endpoint"),
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

			if !fctl.CheckStackApprobation(cmd, stack, "You are about to install connector '%s'", internal.MangoPayConnector) {
				return fctl.ErrMissingApproval
			}

			paymentsClient, err := fctl.NewStackClient(cmd, cfg, stack)
			if err != nil {
				return err
			}

			request := operations.InstallConnectorRequest{
				Connector: shared.ConnectorMangopay,
				RequestBody: shared.MangoPayConfig{
					ClientID: args[0],
					APIKey:   args[1],
					Endpoint: fctl.GetString(cmd, endpointFlag),
				},
			}
			response, err := paymentsClient.Payments.InstallConnector(cmd.Context(), request)
			if err != nil {
				return errors.Wrap(err, "installing connector")
			}

			if response.StatusCode >= 300 {
				return fmt.Errorf("unexpected status code: %d", response.StatusCode)
			}

			pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Connector installed!")

			return nil
		}),
	)
}
