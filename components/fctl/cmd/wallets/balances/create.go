package balances

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type CreateStore struct {
	BalanceName string `json:"balance_name"`
}
type CreateController struct {
	store *CreateStore
}

var _ fctl.Controller[*CreateStore] = (*CreateController)(nil)

func NewDefaultCreateStore() *CreateStore {
	return &CreateStore{}
}

func NewCreateController() *CreateController {
	return &CreateController{
		store: NewDefaultCreateStore(),
	}
}

func NewCreateCommand() *cobra.Command {
	return fctl.NewCommand("create <balance-name>",
		fctl.WithShortDescription("Create a new balance"),
		fctl.WithAliases("c", "cr"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		internal.WithTargetingWalletByID(),
		internal.WithTargetingWalletByName(),
		fctl.WithController[*CreateStore](NewCreateController()),
	)
}

func (c *CreateController) GetStore() *CreateStore {
	return c.store
}

func (c *CreateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config")
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	walletID, err := internal.RequireWalletID(cmd, client)
	if err != nil {
		return nil, err
	}

	request := operations.CreateBalanceRequest{
		ID: walletID,
		CreateBalanceRequest: &shared.CreateBalanceRequest{
			Name: args[0],
		},
	}
	response, err := client.Wallets.CreateBalance(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "creating balance")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.BalanceName = response.CreateBalanceResponse.Data.Name
	return c, nil
}

func (c *CreateController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln(
		"Balance created successfully with name: %s", c.store.BalanceName)
	return nil

}
