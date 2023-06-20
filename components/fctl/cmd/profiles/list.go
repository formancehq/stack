package profiles

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ProfilesListStore struct {
	Profiles []string `json:"profiles"`
}
type ProfilesListController struct {
	store              *ProfilesListStore
	currentProfileName string
}

var _ fctl.Controller[*ProfilesListStore] = (*ProfilesListController)(nil)

func NewDefaultProfilesListStore() *ProfilesListStore {
	return &ProfilesListStore{
		Profiles: []string{},
	}
}

func NewProfilesListController() *ProfilesListController {
	return &ProfilesListController{
		store:              NewDefaultProfilesListStore(),
		currentProfileName: "",
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

	c.store.Profiles = fctl.MapKeys(cfg.GetProfiles())

	c.currentProfileName = fctl.GetCurrentProfileName(cmd, cfg)

	return c, nil

}

func (c *ProfilesListController) Render(cmd *cobra.Command, args []string) error {
	tableData := fctl.Map(c.store.Profiles, func(p string) []string {
		isCurrent := "No"
		if p == c.currentProfileName {
			isCurrent = "Yes"
		}
		return []string{p, isCurrent}
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
