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

type DeleteWebhookStore struct {
	ErrorResponse *shared.ErrorResponse `json:"err"`
	Success       bool                  `json:"success"`
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
	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "fctl.GetConfig")
	}
	c.config = cfg

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to delete a webhook") {
		return nil, fctl.ErrMissingApproval
	}

	webhookClient, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	request := operations.DeleteConfigRequest{
		ID: args[0],
	}
	response, err := webhookClient.Webhooks.DeleteConfig(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "deleting config")
	}

	if response.ErrorResponse != nil {
		if response.ErrorResponse.ErrorCode == "NOT_FOUND" {
			c.store.ErrorResponse = response.ErrorResponse
			c.store.Success = false
			return c, nil
		}
		return nil, fmt.Errorf("%s: %s", response.ErrorResponse.ErrorCode, response.ErrorResponse.ErrorMessage)
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
