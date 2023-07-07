package wallets

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStore struct {
	Wallets []shared.Wallet `json:"wallets"`
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
		fctl.WithShortDescription("List all wallets"),
		fctl.WithAliases("ls", "l"),
		fctl.WithStringSliceFlag(c.metadataFlag, []string{""}, "Metadata to use"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](c),
	)
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

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

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(cmd, c.metadataFlag))
	if err != nil {
		return nil, err
	}

	response, err := client.Wallets.ListWallets(cmd.Context(), operations.ListWalletsRequest{
		Metadata: metadata,
	})
	if err != nil {
		return nil, errors.Wrap(err, "listing wallet")
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Wallets = response.ListWalletsResponse.Cursor.Data

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	if len(c.store.Wallets) == 0 {
		fctl.Println("No wallets found.")
		return nil
	}

	if err := pterm.DefaultTable.
		WithHasHeader(true).
		WithWriter(cmd.OutOrStdout()).
		WithData(
			fctl.Prepend(
				fctl.Map(c.store.Wallets,
					func(src shared.Wallet) []string {
						return []string{
							src.ID,
							src.Name,
						}
					}),
				[]string{"ID", "Name"},
			),
		).Render(); err != nil {
		return errors.Wrap(err, "rendering table")
	}
	return nil
}
