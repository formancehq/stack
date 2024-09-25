package pools

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type RemoveAccountStore struct {
	PoolID    string `json:"poolID"`
	AccountID string `json:"accountID"`
	Success   bool   `json:"success"`
}
type RemoveAccountController struct {
	PaymentsVersion versions.Version

	store *RemoveAccountStore
}

func (c *RemoveAccountController) SetVersion(version versions.Version) {
	c.PaymentsVersion = version
}

var _ fctl.Controller[*RemoveAccountStore] = (*RemoveAccountController)(nil)

func NewRemoveAccountStore() *RemoveAccountStore {
	return &RemoveAccountStore{}
}

func NewRemoveAccountController() *RemoveAccountController {
	return &RemoveAccountController{
		store: NewRemoveAccountStore(),
	}
}

func NewRemoveAccountCommand() *cobra.Command {
	c := NewRemoveAccountController()
	return fctl.NewCommand("remove-account <poolID> <accountID>",
		fctl.WithShortDescription("Remove account from pool"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithAliases("add", "a"),
		fctl.WithController[*RemoveAccountStore](c),
	)
}

func (c *RemoveAccountController) GetStore() *RemoveAccountStore {
	return c.store
}

func (c *RemoveAccountController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	if err := versions.GetPaymentsVersion(cmd, args, c); err != nil {
		return nil, err
	}

	if c.PaymentsVersion < versions.V1 {
		return nil, fmt.Errorf("pools are only supported in >= v1.0.0")
	}

	response, err := store.Client().Payments.V1.RemoveAccountFromPool(cmd.Context(), operations.RemoveAccountFromPoolRequest{
		PoolID:    args[0],
		AccountID: args[1],
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.PoolID = args[0]
	c.store.AccountID = args[1]

	return c, nil
}

func (c *RemoveAccountController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Successfully removed '%s' to '%s'", c.store.AccountID, c.store.PoolID)

	return nil
}
