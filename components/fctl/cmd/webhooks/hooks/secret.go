package hooks

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ChangeSecretStore struct {
	Secret string `json:"secret"`
	ID     string `json:"id"`
	ErrorResponse error
	
}

type ChangeSecretWebhookController struct {
	store *ChangeSecretStore
}

var _ fctl.Controller[*ChangeSecretStore] = (*ChangeSecretWebhookController)(nil)

func NewDefaultChangeSecretStore() *ChangeSecretStore {
	return &ChangeSecretStore{
		Secret: "",
		ID:     "",
	}
}
func NewChangeSecretWebhookController() *ChangeSecretWebhookController {
	return &ChangeSecretWebhookController{
		store: NewDefaultChangeSecretStore(),
	}
}
func (c *ChangeSecretWebhookController) GetStore() *ChangeSecretStore {
	return c.store
}
func (c *ChangeSecretWebhookController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to change a webhook secret") {
		return nil, fctl.ErrMissingApproval
	}
	secret := ""
	if len(args) > 1 {
		secret = args[1]
	}

	request := operations.UpdateSecretHookRequest{
		RequestBody: operations.UpdateSecretHookRequestBody{
			Secret: &secret,
		},
		HookID: args[0],
	}

	response, err := store.Client().Webhooks.UpdateSecretHook(cmd.Context(), request)

		
	if err != nil {
		c.store.ErrorResponse = err
	} else {
		c.store.ID = response.V2HookResponse.Data.ID
		c.store.Secret = response.V2HookResponse.Data.Secret
	}

	

	return c, nil
}

func (c *ChangeSecretWebhookController) Render(cmd *cobra.Command, args []string) error {
	if c.store.ErrorResponse != nil {
		pterm.Warning.WithShowLineNumber(false).Printfln(c.store.ErrorResponse.Error())
		return nil
	}
	
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln(
		"Hook '%s' updated successfully with new secret : %s", c.store.ID, c.store.Secret)
	return nil
}

func NewChangeSecretCommand() *cobra.Command {
	return fctl.NewCommand("change-secret <hook-id> <secret>",
		fctl.WithShortDescription("Change the signing secret of a config. You can bring your own secret. If not passed or empty, a secret is automatically generated. The format is a string of bytes of size 24, base64 encoded. (larger size after encoding)"),
		fctl.WithConfirmFlag(),
		fctl.WithAliases("cs"),
		fctl.WithArgs(cobra.RangeArgs(1, 2)),
		fctl.WithController[*ChangeSecretStore](NewChangeSecretWebhookController()),
	)
}
