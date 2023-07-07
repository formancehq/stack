package profiles

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type Profile struct {
	Name   string `json:"name"`
	Active string `json:"active"`
}
type ProfilesListStore struct {
	Profiles []*Profile `json:"profiles"`
}
type ProfilesListController struct {
	store *ProfilesListStore
}

var _ fctl.Controller[*ProfilesListStore] = (*ProfilesListController)(nil)

func NewDefaultProfilesListStore() *ProfilesListStore {
	return &ProfilesListStore{
		Profiles: []*Profile{},
	}
}

func NewProfilesListController() *ProfilesListController {
	return &ProfilesListController{
		store: NewDefaultProfilesListStore(),
	}
}

func (c *ProfilesListController) GetStore() *ProfilesListStore {
	return c.store
}

func (c *ProfilesListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	p := fctl.MapKeys(cfg.GetProfiles())
	currentProfileName := fctl.GetCurrentProfileName(cmd, cfg)

	for _, k := range p {
		c.store.Profiles = append(c.store.Profiles, &Profile{
			Name: k,
			Active: func(k string) string {
				if currentProfileName == k {
					return "Yes"
				}
				return "No"
			}(k),
		})
	}

	return c, nil

}

func (c *ProfilesListController) Render(cmd *cobra.Command, args []string) error {
	tableData := fctl.Map(c.store.Profiles, func(p *Profile) []string {
		return []string{
			p.Name,
			p.Active,
		}
	})
	tableData = fctl.Prepend(tableData, []string{"Name", "Active"})

	pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
	return nil
}

func NewListCommand() *cobra.Command {
	return fctl.NewCommand("list",
		fctl.WithAliases("l"),
		fctl.WithShortDescription("List profiles"),
		fctl.WithController[*ProfilesListStore](NewProfilesListController()),
	)
}
