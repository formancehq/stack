package webhooks

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func RunActivateWebhookfunc(cmd *cobra.Command, args []string) error {
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

	if !fctl.CheckStackApprobation(cmd, stack, "You are bout to activate a webhook") {
		return fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return errors.Wrap(err, "creating stack client")
	}

	request := operations.ActivateConfigRequest{
		ID: args[0],
	}
	response, err := client.Webhooks.ActivateConfig(cmd.Context(), request)
	if err != nil {
		return errors.Wrap(err, "activating config")
	}

	if response.ErrorResponse != nil {
		return fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	output := &OutputActivateWebhook{
		Success: response.ConfigResponse.Data.Active,
	}

	fctl.SetSharedData(output, nil, nil, nil)

	return nil
}

type OutputActivateWebhook struct {
	Success bool `json:"success"`
}

func DisplayOutputActivateWebhook(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Config activated successfully")

	return nil
}

func NewActivateCommand() *cobra.Command {
	return fctl.NewCommand("activate <config-id>",
		fctl.WithShortDescription("Activate one config"),
		fctl.WithAliases("ac", "a"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithRunE(RunActivateWebhookfunc),
		fctl.WrapOutputPostRunE(DisplayOutputActivateWebhook),
	)
}
