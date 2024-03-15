package webhooks

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/webhooks/store"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DeleteWebhookStore struct {
	ErrorResponse *sdkerrors.WebhooksErrorResponse `json:"error"`
	Success       bool                             `json:"success"`
}

type DeleteWebhookController struct {
	store  *DeleteWebhookStore
	config *fctl.Config
}

var _ fctl.Controller[*DeleteWebhookStore] = (*DeleteWebhookController)(nil)

func NewDefaultDeleteWebhookStore() *DeleteWebhookStore {
	return &DeleteWebhookStore{
		Success: true,
	}
}

func NewDeleteWebhookController() *DeleteWebhookController {
	return &DeleteWebhookController{
		store:  NewDefaultDeleteWebhookStore(),
		config: nil,
	}
}

func (c *DeleteWebhookController) GetStore() *DeleteWebhookStore {
	return c.store
}

func (c *DeleteWebhookController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := store.GetStore(cmd.Context())
	c.config = store.Config

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to delete a webhook") {
		return nil, fctl.ErrMissingApproval
	}

	request := operations.DeleteConfigRequest{
		ID: args[0],
	}
	response, err := store.Client().Webhooks.DeleteConfig(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "deleting config")
	}

	if response.WebhooksErrorResponse != nil {
		if response.WebhooksErrorResponse.ErrorCode == "NOT_FOUND" {
			c.store.ErrorResponse = response.WebhooksErrorResponse
			c.store.Success = false
			return c, nil
		}
		return nil, fmt.Errorf("%s: %s", response.WebhooksErrorResponse.ErrorCode, response.WebhooksErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = response.StatusCode == 200

	return c, nil
}

func (c *DeleteWebhookController) Render(cmd *cobra.Command, args []string) error {
	if !c.store.Success {
		pterm.Warning.WithShowLineNumber(false).Printfln("Config %s not found", args[0])
		return nil
	}

	if c.store.ErrorResponse != nil {
		pterm.Warning.WithShowLineNumber(false).Printfln(c.store.ErrorResponse.ErrorMessage)
		return nil
	}

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Config deleted successfully")

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
