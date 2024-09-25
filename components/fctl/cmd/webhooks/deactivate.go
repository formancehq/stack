package webhooks

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DesactivateWebhookStore struct {
	Success bool `json:"success"`
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

	request := operations.DeactivateConfigRequest{
		ID: args[0],
	}
	response, err := store.Client().Webhooks.V1.DeactivateConfig(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "deactivating config")
	}

	c.store.Success = !response.ConfigResponse.Data.Active

	return c, nil
}

func (c *DesactivateWebhookController) Render(cmd *cobra.Command, args []string) error {

	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Config deactivated successfully")

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
