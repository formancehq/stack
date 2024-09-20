package regions

import (
	"fmt"
	"strings"
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
	store := fctl.GetMembershipStore(cmd.Context())

	organizationID, err := fctl.ResolveOrganizationID(cmd, store.Config, store.Client())
	if err != nil {
		return nil, err
	}

	response, _, err := store.Client().GetRegion(cmd.Context(), organizationID, args[0]).Execute()
	if err != nil {
		return nil, err
	}

	c.store.Region = response.Data

	return c, nil
}

func (c *ShowController) Render(cmd *cobra.Command, args []string) (err error) {
	fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Information")
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.Region.Id})
	tableData = append(tableData, []string{pterm.LightCyan("Name"), c.store.Region.Name})
	tableData = append(tableData, []string{pterm.LightCyan("Base URL"), c.store.Region.BaseUrl})
	tableData = append(tableData, []string{pterm.LightCyan("Active"), fctl.BoolToString(c.store.Region.Active)})
	tableData = append(tableData, []string{pterm.LightCyan("Public"), fctl.BoolToString(c.store.Region.Public)})

	if c.store.Region.Version != nil {
		tableData = append(tableData, []string{pterm.LightCyan("Version"), *c.store.Region.Version})
	}

	if c.store.Region.Creator != nil {
		tableData = append(tableData, []string{pterm.LightCyan("Creator"), c.store.Region.Creator.Email})
	}
	if c.store.Region.LastPing != nil {
		tableData = append(tableData, []string{pterm.LightCyan("Last ping"), c.store.Region.LastPing.Format(time.RFC3339)})
	}

	err = pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
	if err != nil {
		return
	}

	tableData = pterm.TableData{}
	capabilities, err := fctl.StructToMap(c.store.Region.Capabilities)
	if err != nil {
		return
	}
	if len(capabilities) > 0 {
		fctl.Section.WithWriter(cmd.OutOrStdout()).Println("Capabilities")
	}
	for key, value := range capabilities {
		data := []string{
			pterm.LightCyan(key),
		}

		var v []string
		if value != nil {
			c, ok := value.([]interface{})
			if ok {
				for _, val := range c {
					v = append(v, fmt.Sprintf("%v", val))
				}
			}
		}
		data = append(data, strings.Join(v, ", "))
		tableData = append(tableData, data)
	}

	return pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()

}
