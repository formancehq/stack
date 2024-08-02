package hooks

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DesactivateWebhookStore struct {
	Success bool `json:"success"`
	ErrorResponse error `json:"error"`
	
}

type DesactivateWebhookController struct {
	store *DesactivateWebhookStore
}

var _ fctl.Controller[*DesactivateWebhookStore] = (*DesactivateWebhookController)(nil)

func NewDefaultDesactivateWebhookStore() *DesactivateWebhookStore {
	return &DesactivateWebhookStore{
		Success: true,
	}
}

func NewDesactivateWebhookController() *DesactivateWebhookController {
	return &DesactivateWebhookController{
		store: NewDefaultDesactivateWebhookStore(),
	}
}
func (c *DesactivateWebhookController) GetStore() *DesactivateWebhookStore {
	return c.store
}

func (c *DesactivateWebhookController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to deactivate a webhook") {
		return nil, fctl.ErrMissingApproval
	}

	request := operations.DeactivateHookRequest{
		HookID: args[0],
	}
	_, err := store.Client().Webhooks.DeactivateHook(cmd.Context(), request)
	
	if err!= nil {
		c.store.ErrorResponse = err
	} else {
		c.store.Success = true
	}

	return c, nil
}

func (c *DesactivateWebhookController) Render(cmd *cobra.Command, args []string) error {
	if c.store.ErrorResponse != nil {
		pterm.Warning.WithShowLineNumber(false).Printfln(c.store.ErrorResponse.Error())
		return nil
	}


	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Hook deactivated successfully")

	return nil
}

func NewDeactivateCommand() *cobra.Command {
	return fctl.NewCommand("deactivate <config-id>",
		fctl.WithShortDescription("Deactivate one config"),
		fctl.WithConfirmFlag(),
		fctl.WithAliases("deac"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*DesactivateWebhookStore](NewDesactivateWebhookController()),
	)
}
