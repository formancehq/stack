package instances

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewSendEventCommand() *cobra.Command {
	return fctl.NewCommand("send-event <instance-id> <event>",
		fctl.WithShortDescription("Send an event to an instance"),
		fctl.WithArgs(cobra.ExactArgs(2)),
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

			_, err = client.OrchestrationApi.SendEvent(cmd.Context(), args[0]).SendEventRequest(formance.SendEventRequest{
				Name: args[1],
			}).Execute()
			if err != nil {
				return err
			}

			pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Event '%s' sent", args[1])

			return nil
		}),
	)
}
