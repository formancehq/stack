package regions

import (
	"flag"
	"time"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useList   = "list"
	shortList = "List all private regions"
)

type ListStore struct {
	Regions []membershipclient.AnyRegion `json:"regions"`
}

func NewListStore() *ListStore {
	return &ListStore{}
}

func NewListConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	return fctl.NewControllerConfig(
		useList,
		shortList,
		shortList,
		[]string{
			"ls", "l",
		},
		flags,
		fctl.Organization,
	)
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

type ListController struct {
	store  *ListStore
	config *fctl.ControllerConfig
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

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	regionsResponse, _, err := apiClient.DefaultApi.ListRegions(ctx, organizationID).Execute()
	if err != nil {
		return nil, err
	}

	c.store.Regions = regionsResponse.Data

	return c, nil
}

func (c *ListController) Render() error {
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
		}
	})
	tableData = fctl.Prepend(tableData, []string{"ID", "Name", "Base url", "Public", "Active", "Last ping", "Owner"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()
}

func NewListCommand() *cobra.Command {

	config := NewListConfig()
	return fctl.NewCommand("list",
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController(config)),
	)
}
