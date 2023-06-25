package wallets

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/wallets/internal"
	"github.com/formancehq/fctl/cmd/wallets/internal/views"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type ShowStore struct {
	Wallet shared.WalletWithBalances `json:"wallet"`
}
type ShowController struct {
	store *ShowStore
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

func NewDefaultShowStore() *ShowStore {
	return &ShowStore{}
}

func NewShowController() *ShowController {
	return &ShowController{
		store: NewDefaultShowStore(),
	}
}

func NewShowCommand() *cobra.Command {
	c := NewShowController()
	return fctl.NewCommand("show",
		fctl.WithShortDescription("Show a wallets"),
		fctl.WithAliases("sh"),
		fctl.WithConfirmFlag(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		internal.WithTargetingWalletByID(),
		internal.WithTargetingWalletByName(),
		fctl.WithController[*ShowStore](c),
	)
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	walletID, err := internal.RetrieveWalletID(cmd, client)
	if err != nil {
		return nil, err
	}
	if walletID == "" {
		return nil, errors.New("You need to specify wallet id using --id or --name flags")
	}

	response, err := client.Wallets.GetWallet(cmd.Context(), operations.GetWalletRequest{
		ID: walletID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "getting wallet")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Wallet = response.GetWalletResponse.Data

	return c, nil
}

func (c *ShowController) Render(cmd *cobra.Command, args []string) error {
	return views.PrintWalletWithMetadata(cmd.OutOrStdout(), c.store.Wallet)
}
