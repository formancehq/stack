package balances

import (
	"flag"
	"fmt"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useCreate   = "create <balance-name>"
	shortCreate = "Create a balance"
)

type CreateStore struct {
	BalanceName string `json:"balanceName"`
}
type CreateController struct {
	store  *CreateStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*CreateStore] = (*CreateController)(nil)

func NewCreateStore() *CreateStore {
	return &CreateStore{}
}
func NewCreateConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useCreate, flag.ExitOnError)

	fctl.WithConfirmFlag(flags)

	internal.WithTargetingWalletByName(flags)
	internal.WithTargetingWalletByID(flags)

	c := fctl.NewControllerConfig(
		useCreate,
		shortCreate,
		shortCreate,
		[]string{
			"c", "cr",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)

	return c
}
func NewCreateController(config *fctl.ControllerConfig) *CreateController {
	return &CreateController{
		store:  NewCreateStore(),
		config: config,
	}
}

func (c *CreateController) GetStore() *CreateStore {
	return c.store
}

func (c *CreateController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *CreateController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()
	out := c.config.GetOut()

	if len(args) != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", len(args))
	}

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config")
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	walletID, err := internal.RequireWalletID(flags, ctx, client)
	if err != nil {
		return nil, err
	}

	request := operations.CreateBalanceRequest{
		ID: walletID,
		CreateBalanceRequest: &shared.CreateBalanceRequest{
			Name: args[0],
		},
	}
	response, err := client.Wallets.CreateBalance(ctx, request)
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

func (c *CreateController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln(
		"Balance created successfully with name: %s", c.store.BalanceName)
	return nil

}
func NewCreateCommand() *cobra.Command {
	c := NewCreateConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*CreateStore](NewCreateController(c)),
	)
}
