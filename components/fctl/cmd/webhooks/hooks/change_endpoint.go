package hooks

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ChangeEndpointWebhookStore struct {
	Success bool `json:"success"`
	ErrorResponse error `json:"error"`
	
}
type ChangeEndpointWebhookController struct {
	store *ChangeEndpointWebhookStore
}

func NewDefaultChangeEndpointWebhookVersionStore() *ChangeEndpointWebhookStore {
	return &ChangeEndpointWebhookStore{
		Success: false,
	}
}
func NewChangeEndpointWebhookController() *ChangeEndpointWebhookController {
	return &ChangeEndpointWebhookController{
		store: NewDefaultChangeEndpointWebhookVersionStore(),
	}
}
func (c *ChangeEndpointWebhookController) GetStore() *ChangeEndpointWebhookStore {
	return c.store
}

func (c *ChangeEndpointWebhookController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())
	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are bout to change endpoint for a webhook") {
		return nil, fctl.ErrMissingApproval
	}

	request := operations.UpdateEndpointHookRequest{
		RequestBody: operations.UpdateEndpointHookRequestBody{
			Endpoint: &args[1],
		},
		HookID: args[0],
	}

	_, err := store.Client().Webhooks.UpdateEndpointHook(cmd.Context(), request)

	if err != nil {
		c.store.ErrorResponse = err
	} else {
		c.store.Success = true
	}

	return c, nil
}

func (c *ChangeEndpointWebhookController) Render(cmd *cobra.Command, args []string) error {
	if c.store.ErrorResponse != nil {
		pterm.Warning.WithShowLineNumber(false).Printfln(c.store.ErrorResponse.Error())
		return nil
	}
	
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Endpoint changed with success")

	return nil
}

func NewChangeEndpointCommand() *cobra.Command {
	return fctl.NewCommand("set-endpoint <hook-id> <new-endpoint>",
		fctl.WithShortDescription("Change endpoint for one hook"),
		fctl.WithAliases("se", "sendpoint"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithController[*ChangeEndpointWebhookStore](NewChangeEndpointWebhookController()),
	)
}
