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

type ChangeSecretOutput struct {
	Secret string `json:"secret"`
	ID     string `json:"id"`
}

type ChangeSecretWebhook struct {
	store *fctl.SharedStore
}

func NewChangeSecretWebhook() *ChangeSecretWebhook {
	return &ChangeSecretWebhook{
		store: fctl.NewSharedStore(),
	}
}
func (c *ChangeSecretWebhook) GetStore() *fctl.SharedStore {
	return c.store
}
func (c *ChangeSecretWebhook) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "fctl.GetConfig")
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to change a webhook secret") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
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
		return nil, errors.Wrap(err, "changing secret")
	}

	if response.ErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	output := &ChangeSecretOutput{
		Secret: response.ConfigResponse.Data.Secret,
		ID:     response.ConfigResponse.Data.ID,
	}

	c.GetStore().SetData(output)

	return c, nil
}

func (c *ChangeSecretWebhook) Render(cmd *cobra.Command, args []string) error {
	Data, ok := c.GetStore().GetData().(*ChangeSecretOutput)
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
		fctl.WithController(NewChangeSecretWebhook()),
		// fctl.WrapOutputPostRunE(DisplayWebhooks),
	)
}
