package webhooks

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ActivateWebhookStore struct {
	Success bool `json:"success"`
}
type ActivateWebhookController struct {
	store *ActivateWebhookStore
}

func NewDefaultVersionStore() *ActivateWebhookStore {
	return &ActivateWebhookStore{
		Success: true,
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

	request := operations.ActivateConfigRequest{
		ID: args[0],
	}

	_, err := store.Client().Webhooks.V1.ActivateConfig(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "activating config")
	}

	return c, nil
}

func (*ActivateWebhookController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Config activated successfully")

	return nil
}

func NewActivateCommand() *cobra.Command {
	return fctl.NewCommand("activate <config-id>",
		fctl.WithShortDescription("Activate one config"),
		fctl.WithAliases("ac", "a"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*ActivateWebhookStore](NewActivateWebhookController()),
	)
}
