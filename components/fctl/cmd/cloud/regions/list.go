package regions

import (
	"time"

	"github.com/formancehq/fctl/cmd/cloud/store"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListStore struct {
	Regions []membershipclient.AnyRegion `json:"regions"`
}
type ListController struct {
	store *ListStore
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

func NewDefaultListStore() *ListStore {
	return &ListStore{}
}

func NewListController() *ListController {
	return &ListController{
		store: NewDefaultListStore(),
	}
}

func NewListCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List regions"),
		fctl.WithController[*ListStore](NewListController()),
	)
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := store.GetStore(cmd.Context())

	organizationID, err := fctl.ResolveOrganizationID(cmd, store.Config, store.Client())
	if err != nil {
		return nil, err
	}

	regionsResponse, _, err := store.Client().ListRegions(cmd.Context(), organizationID).Execute()
	if err != nil {
		return nil, err
	}

	c.store.Regions = regionsResponse.Data

	return c, nil
}

func (c *ListController) Render(cmd *cobra.Command, args []string) error {
	tableData := fctl.Map(c.store.Regions, func(i membershipclient.AnyRegion) []string {
		return []string{
			i.Id,
			i.Name,
			i.BaseUrl,
			fctl.BoolToString(i.Public),
			fctl.BoolToString(i.Active),
			func() string {
				if i.LastPing != nil {
					return i.LastPing.Format(time.RFC3339)
				}
				return ""
			}(),
			func() string {
				if i.Creator != nil {
					return i.Creator.Email
				}
				return "Formance Cloud"
			}(),
			func() string {
				if i.Version == nil {
					return ""
				}
				return *i.Version
			}(),
		}
	})
	tableData = fctl.Prepend(tableData, []string{"ID", "Name", "Base url", "Public", "Active", "Last ping", "Owner", "Version"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}
