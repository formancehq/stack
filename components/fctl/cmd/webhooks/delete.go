package webhooks

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func RunDeleteCommand(cmd *cobra.Command, args []string) error {
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

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to delete a webhook") {
		return fctl.ErrMissingApproval
	}

	webhookClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return errors.Wrap(err, "creating stack client")
	}

	request := operations.DeleteConfigRequest{
		ID: args[0],
	}
	response, err := webhookClient.Webhooks.DeleteConfig(cmd.Context(), request)
	if err != nil {
		return errors.Wrap(err, "deleting config")
	}

	if response.ErrorResponse != nil {
		return fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	output := &OutputDeleteWebhook{
		Success: response.StatusCode == 200,
	}

	fctl.SetSharedData(output, nil, cfg, nil)

	return nil
}

type OutputDeleteWebhook struct {
	Success bool `json:"success"`
}

func DisplayDeleteCommand(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Config deleted successfully")

	return nil
}

func NewDeleteCommand() *cobra.Command {
	return fctl.NewCommand("delete <config-id>",
		fctl.WithShortDescription("Delete a config"),
		fctl.WithConfirmFlag(),
		fctl.WithAliases("del"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithRunE(RunDeleteCommand),
		fctl.WrapOutputPostRunE(DisplayDeleteCommand),
	)
}
