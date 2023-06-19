package webhooks

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func RunChangeSecret(cmd *cobra.Command, args []string) error {
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

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to change a webhook secret") {
		return fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return errors.Wrap(err, "creating stack client")
	}

	secret := ""
	if len(args) > 1 {
		secret = args[1]
	}

	response, err := client.Webhooks.
		ChangeConfigSecret(cmd.Context(), operations.ChangeConfigSecretRequest{
			ConfigChangeSecret: &shared.ConfigChangeSecret{
				Secret: secret,
			},
			ID: args[0],
		})
	if err != nil {
		return errors.Wrap(err, "changing secret")
	}

	if response.ErrorResponse != nil {
		return fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	output := &ChangeSecretOutput{
		Secret: response.ConfigResponse.Data.Secret,
		ID:     response.ConfigResponse.Data.ID,
	}

	fctl.SetSharedData(output, nil, nil, nil)

	return nil
}

type ChangeSecretOutput struct {
	Secret string `json:"secret"`
	ID     string `json:"id"`
}

func DisplayWebhooks(cmd *cobra.Command, args []string) error {
	Data, ok := fctl.GetSharedData().(*ChangeSecretOutput)
	if !ok {
		return errors.New("unable to get shared data")
	}

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln(
		"Config '%s' updated successfully with new secret", Data.ID)
	return nil
}

func NewChangeSecretCommand() *cobra.Command {
	return fctl.NewCommand("change-secret <config-id> <secret>",
		fctl.WithShortDescription("Change the signing secret of a config. You can bring your own secret. If not passed or empty, a secret is automatically generated. The format is a string of bytes of size 24, base64 encoded. (larger size after encoding)"),
		fctl.WithConfirmFlag(),
		fctl.WithAliases("cs"),
		fctl.WithArgs(cobra.RangeArgs(1, 2)),
		fctl.WithRunE(RunChangeSecret),
		fctl.WrapOutputPostRunE(DisplayWebhooks),
	)
}
