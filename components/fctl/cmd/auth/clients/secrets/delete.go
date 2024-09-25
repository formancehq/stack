package secrets

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DeleteStore struct {
	Success  bool   `json:"success"`
	SecretId string `json:"secretId"`
}
type DeleteController struct {
	store *DeleteStore
}

var _ fctl.Controller[*DeleteStore] = (*DeleteController)(nil)

func NewDefaultDeleteStore() *DeleteStore {
	return &DeleteStore{}
}

func NewDeleteController() *DeleteController {
	return &DeleteController{
		store: NewDefaultDeleteStore(),
	}
}

func NewDeleteCommand() *cobra.Command {
	return fctl.NewCommand("delete <client-id> <secret-id>",
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithAliases("d"),
		fctl.WithShortDescription("Delete secret"),
		fctl.WithConfirmFlag(),
		fctl.WithController[*DeleteStore](NewDeleteController()),
	)
}

func (c *DeleteController) GetStore() *DeleteStore {
	return c.store
}

func (c *DeleteController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := fctl.GetStackStore(cmd.Context())

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to delete a client secret") {
		return nil, fctl.ErrMissingApproval
	}

	request := operations.DeleteSecretRequest{
		ClientID: args[0],
		SecretID: args[1],
	}
	response, err := store.Client().Auth.V1.DeleteSecret(cmd.Context(), request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.SecretId = args[1]
	c.store.Success = true

	return c, nil
}

func (c *DeleteController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Secret %s successfully deleted!", c.store.SecretId)

	return nil

}
