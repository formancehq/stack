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
	useShow   = "show <region-id>"
	shortShow = "Show region details with id"
)

type ShowStore struct {
	Region membershipclient.AnyRegion `json:"region"`
}

func NewShowStore() *ShowStore {
	return &ShowStore{}
}

func NewShowConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useShow, flag.ExitOnError)
	return fctl.NewControllerConfig(
		useShow,
		shortShow,
		shortShow,
		[]string{
			"sh", "s",
		},
		flags,
		fctl.Organization,
	)
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

type ShowController struct {
	store  *ShowStore
	config *fctl.ControllerConfig
}

func NewShowController(config *fctl.ControllerConfig) *ShowController {
	return &ShowController{
		store:  NewShowStore(),
		config: config,
	}
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *ShowController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()

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

	response, _, err := apiClient.DefaultApi.GetRegion(ctx, organizationID, args[0]).Execute()
	if err != nil {
		return nil, err
	}

	c.store.Region = response.Data

	return c, nil
}

func (c *ShowController) Render() error {
	fctl.Section.WithWriter(c.config.GetOut()).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.Region.Id})
	tableData = append(tableData, []string{pterm.LightCyan("Name"), c.store.Region.Name})
	tableData = append(tableData, []string{pterm.LightCyan("Base URL"), c.store.Region.BaseUrl})
	tableData = append(tableData, []string{pterm.LightCyan("Active: "), fctl.BoolToString(c.store.Region.Active)})
	tableData = append(tableData, []string{pterm.LightCyan("Public: "), fctl.BoolToString(c.store.Region.Public)})
	if c.store.Region.Creator != nil {
		tableData = append(tableData, []string{pterm.LightCyan("Creator"), c.store.Region.Creator.Email})
	}
	if c.store.Region.LastPing != nil {
		tableData = append(tableData, []string{pterm.LightCyan("Base URL"), c.store.Region.LastPing.Format(time.RFC3339)})
	}

	return pterm.DefaultTable.
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()

}
func NewShowCommand() *cobra.Command {

	config := NewShowConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*ShowStore](NewShowController(config)),
	)
}
