package webhooks

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func DesactivateWebHookCommand(cmd *cobra.Command, args []string) error {
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

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to deactivate a webhook") {
		return fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return errors.Wrap(err, "creating stack client")
	}

	request := operations.DeactivateConfigRequest{
		ID: args[0],
	}
	response, err := client.Webhooks.DeactivateConfig(cmd.Context(), request)
	if err != nil {
		return errors.Wrap(err, "deactivating config")
	}

	output := &OutputDesactivateWebhook{
		Success: !response.ConfigResponse.Data.Active,
	}

	// Check if there is an error
	if response.ErrorResponse != nil {
		return fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	// Check if the status code is >= 300
	if response.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	fctl.SetSharedData(output, nil, cfg, nil)

	return nil
}

type OutputDesactivateWebhook struct {
	Success bool `json:"success"`
}

func OutputDesactivateCommand(cmd *cobra.Command, args []string) error {

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Config deactivated successfully")

	return nil
}

func NewDeactivateCommand() *cobra.Command {
	return fctl.NewCommand("deactivate <config-id>",
		fctl.WithShortDescription("Deactivate one config"),
		fctl.WithConfirmFlag(),
		fctl.WithAliases("deac"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithRunE(DesactivateWebHookCommand),
		fctl.WrapOutputPostRunE(OutputDesactivateCommand),
	)
}
