package policies

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DeleteStore struct {
	PolicyID string `json:"policyID"`
	Success  bool   `json:"success"`
}

type DeleteController struct {
	store *DeleteStore
}

var _ fctl.Controller[*DeleteStore] = (*DeleteController)(nil)

func NewDeleteStore() *DeleteStore {
	return &DeleteStore{}
}

func NewDeleteController() *DeleteController {
	return &DeleteController{
		store: NewDeleteStore(),
	}
}
func NewDeleteCommand() *cobra.Command {
	c := NewDeleteController()
	return fctl.NewCommand("delete <policyID>",
		fctl.WithConfirmFlag(),
		fctl.WithAliases("d"),
		fctl.WithShortDescription("Delete a policy"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*DeleteStore](c),
	)
}

func (c *DeleteController) GetStore() *DeleteStore {
	return c.store
}

func (c *DeleteController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := fctl.GetStackStore(cmd.Context())

	if !fctl.CheckStackApprobation(cmd, store.Stack(), "You are about to delete '%s'", args[0]) {
		return nil, fctl.ErrMissingApproval
	}

	response, err := store.Client().Reconciliation.DeletePolicy(
		cmd.Context(),
		operations.DeletePolicyRequest{
			PolicyID: args[0],
		},
	)

	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.PolicyID = args[0]
	c.store.Success = true

	return c, nil
}

func (c *DeleteController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Policy %s Deleted!", c.store.PolicyID)
	return nil
}
