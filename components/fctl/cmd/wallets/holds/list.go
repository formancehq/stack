package holds

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
	useList   = "list"
	shortList = "List holds of a wallets"
)

type ListStore struct {
	Holds []shared.Hold `json:"holds"`
}
type ListController struct {
	store  *ListStore
	config *fctl.ControllerConfig
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewListStore() *ListStore {
	return &ListStore{}
}
func NewListConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	internal.WithTargetingWalletByName(flags)
	internal.WithTargetingWalletByID(flags)
	fctl.WithMetadataFlag(flags)
	return fctl.NewControllerConfig(
		useList,
		shortList,
		shortList,
		[]string{
			"ls", "l",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)

}
func NewListController(config *fctl.ControllerConfig) *ListController {
	return &ListController{
		store:  NewListStore(),
		config: config,
	}
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *ListController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	out := c.config.GetOut()
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

	walletID, err := internal.RetrieveWalletID(flags, ctx, client)
	if err != nil {
		return nil, err
	}

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(flags, fctl.MetadataFlag))
	if err != nil {
		return nil, err
	}

	request := operations.GetHoldsRequest{
		WalletID: &walletID,
		Metadata: metadata,
	}
	response, err := client.Wallets.GetHolds(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "getting holds")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Holds = response.GetHoldsResponse.Cursor.Data

	return c, nil
}

func (c *ListController) Render() error {
	if len(c.store.Holds) == 0 {
		fmt.Fprintln(c.config.GetOut(), "No holds found.")
		return nil
	}

	if err := pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(c.config.GetOut()).
		WithData(
			fctl.Prepend(
				fctl.Map(c.store.Holds,
					func(src shared.Hold) []string {
						return []string{
							src.ID,
							src.WalletID,
							src.Description,
							fctl.MetadataAsShortString(src.Metadata),
						}
					}),
				[]string{"ID", "Wallet ID", "Description", "Metadata"},
			),
		).Render(); err != nil {
		return errors.Wrap(err, "rendering table")
	}

	return nil

}

func NewListCommand() *cobra.Command {
	c := NewListConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.RangeArgs(0, 1)),
		fctl.WithController[*ListStore](NewListController(c)),
	)
}
