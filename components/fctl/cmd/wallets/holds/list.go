package holds

import (
	"github.com/formancehq/fctl/cmd/wallets/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStore struct {
	Holds []shared.Hold `json:"holds"`
}
type ListController struct {
	store        *ListStore
	metadataFlag string
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewDefaultListStore() *ListStore {
	return &ListStore{}
}

func NewListController() *ListController {
	return &ListController{
		store:        NewDefaultListStore(),
		metadataFlag: "metadata",
	}
}

func NewListCommand() *cobra.Command {
	c := NewListController()
	return fctl.NewCommand("list",
		fctl.WithShortDescription("List holds of a wallets"),
		fctl.WithAliases("ls", "l"),
		fctl.WithArgs(cobra.RangeArgs(0, 1)),
		internal.WithTargetingWalletByName(),
		internal.WithTargetingWalletByID(),
		fctl.WithStringSliceFlag(c.metadataFlag, []string{""}, "Metadata to use"),
		fctl.WithController[*ListStore](c),
	)
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := fctl.GetStackStore(cmd.Context())

	walletID, err := internal.RetrieveWalletID(cmd, store.Client())
	if err != nil {
		return nil, err
	}

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, c.metadataFlag))
	if err != nil {
		return nil, err
	}

	request := operations.GetHoldsRequest{
		WalletID: &walletID,
		Metadata: metadata,
	}
	response, err := store.Client().Wallets.V1.GetHolds(cmd.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "getting holds")
	}

	c.store.Holds = response.GetHoldsResponse.Cursor.Data

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	if len(c.store.Holds) == 0 {
		fctl.Println("No holds found.")
		return nil
	}

	if err := pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(cmd.OutOrStdout()).
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
