package hooks

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ActivateWebhookStore struct {
	Success bool `json:"success"`
	ErrorResponse error `json:"error"`
	
}
type ActivateWebhookController struct {
	store *ActivateWebhookStore
}

func NewDefaultVersionStore() *ActivateWebhookStore {
	return &ActivateWebhookStore{
		Success: false,
	}
}
func NewActivateWebhookController() *ActivateWebhookController {
	return &ActivateWebhookController{
		store: NewDefaultVersionStore(),
	}
}
func (c *ActivateWebhookController) GetStore() *ActivateWebhookStore {
	return c.store
}

func (c *ActivateWebhookController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())
	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are bout to activate a webhook") {
		return nil, fctl.ErrMissingApproval
	}

	request := operations.ActivateHookRequest{
		HookID: args[0],
	}
	_, err := store.Client().Webhooks.ActivateHook(cmd.Context(), request)

	if err!= nil {
		c.store.ErrorResponse = err
	} else {
		c.store.Success = true
	}

	return c, nil
}

func (c *ActivateWebhookController) Render(cmd *cobra.Command, args []string) error {
	
	if c.store.ErrorResponse != nil {
		pterm.Warning.WithShowLineNumber(false).Printfln(c.store.ErrorResponse.Error())
		return nil
	}
	
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Hook activated successfully")

	return nil
}

func NewActivateCommand() *cobra.Command {
	return fctl.NewCommand("activate <hook-id>",
		fctl.WithShortDescription("Activate one Hook"),
		fctl.WithAliases("ac", "a"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*ActivateWebhookStore](NewActivateWebhookController()),
	)
}
