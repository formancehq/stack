package hooks

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DeleteWebhookStore struct {
	ErrorResponse  error `json:"error"`
	Success   bool                             `json:"success"`
}

type DeleteWebhookController struct {
	store  *DeleteWebhookStore
}

var _ fctl.Controller[*DeleteWebhookStore] = (*DeleteWebhookController)(nil)

func NewDefaultDeleteWebhookStore() *DeleteWebhookStore {
	return &DeleteWebhookStore{
		Success: false,
	}
}

func NewDeleteWebhookController() *DeleteWebhookController {
	return &DeleteWebhookController{
		store:  NewDefaultDeleteWebhookStore(),
	}
}

func (c *DeleteWebhookController) GetStore() *DeleteWebhookStore {
	return c.store
}

func (c *DeleteWebhookController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to delete a webhook") {
		return nil, fctl.ErrMissingApproval
	}

	request := operations.DeleteHookRequest{
		HookID: args[0],
	}
	_, err := store.Client().Webhooks.DeleteHook(cmd.Context(), request)
	
	if err!= nil {
		
		c.store.ErrorResponse = err
		
	} else {
		c.store.Success = true
	}

	return c, nil
}

func (c *DeleteWebhookController) Render(cmd *cobra.Command, args []string) error {

	if c.store.ErrorResponse != nil {
		pterm.Warning.WithShowLineNumber(false).Printfln(c.store.ErrorResponse.Error())
		return nil
	}


	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Hook deleted successfully")

	return nil
}

func NewDeleteCommand() *cobra.Command {
	return fctl.NewCommand("delete <config-id>",
		fctl.WithShortDescription("Delete a config"),
		fctl.WithConfirmFlag(),
		fctl.WithAliases("del"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*DeleteWebhookStore](NewDeleteWebhookController()),
	)
}
