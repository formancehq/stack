package pools

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/payments/versions"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type AddAccountStore struct {
	PoolID    string `json:"poolID"`
	AccountID string `json:"accountID"`
	Success   bool   `json:"success"`
}
type AddAccountController struct {
	PaymentsVersion versions.Version

	store *AddAccountStore
}

func (c *AddAccountController) SetVersion(version versions.Version) {
	c.PaymentsVersion = version
}

var _ fctl.Controller[*AddAccountStore] = (*AddAccountController)(nil)

func NewAddAccountStore() *AddAccountStore {
	return &AddAccountStore{}
}

func NewAddAccountController() *AddAccountController {
	return &AddAccountController{
		store: NewAddAccountStore(),
	}
}

func NewAddAccountCommand() *cobra.Command {
	c := NewAddAccountController()
	return fctl.NewCommand("add-account <poolID> <accountID>",
		fctl.WithShortDescription("Add account to pool"),
		fctl.WithArgs(cobra.ExactArgs(2)),
		fctl.WithAliases("add", "a"),
		fctl.WithController[*AddAccountStore](c),
	)
}

func (c *AddAccountController) GetStore() *AddAccountStore {
	return c.store
}

func (c *AddAccountController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := fctl.GetStackStore(cmd.Context())

	if err := versions.GetPaymentsVersion(cmd, args, c); err != nil {
		return nil, err
	}

	if c.PaymentsVersion < versions.V1 {
		return nil, fmt.Errorf("pools are only supported in >= v1.0.0")
	}

	response, err := store.Client().Payments.AddAccountToPool(cmd.Context(), operations.AddAccountToPoolRequest{
		PoolID: args[0],
		AddAccountToPoolRequest: shared.AddAccountToPoolRequest{
			AccountID: args[1],
		},
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

func (c *AddAccountController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Successfully added '%s' to '%s'", c.store.AccountID, c.store.PoolID)

	return nil
}
