package regions

import (
	"time"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ShowStore struct {
	Region membershipclient.AnyRegion `json:"region"`
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
	return fctl.NewCommand("show <region-id>",
		fctl.WithAliases("sh", "s"),
		fctl.WithShortDescription("Show region details"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*ShowStore](NewShowController()),
	)
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	response, _, err := apiClient.DefaultApi.GetRegion(cmd.Context(), organizationID, args[0]).Execute()
	if err != nil {
		return nil, err
	}

	c.store.Region = response.Data

	return c, nil
}

func (c *ShowController) Render(cmd *cobra.Command, args []string) error {
	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
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
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()

}
