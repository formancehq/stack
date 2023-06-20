package version

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	Version = "develop"
)

type VersionStore struct {
	Version   string `json:"version" yaml:"version"`
	BuildDate string `json:"buildDate" yaml:"buildDate"`
	Commit    string `json:"commit" yaml:"commit"`
}
type VersionController struct {
	store *VersionStore
}

var _ fctl.Controller[*VersionStore] = (*VersionController)(nil)

func NewDefaultVersionStore() *VersionStore {
	return &VersionStore{
		Version:   "develop",
		BuildDate: "-",
		Commit:    "-",
	}
}

func NewVersionController() *VersionController {
	return &VersionController{
		store: NewDefaultVersionStore(),
	}
}

func NewCommand() *cobra.Command {
	return fctl.NewCommand("version",
		fctl.WithShortDescription("Get version"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*VersionStore](NewVersionController()),
	)
}

func (c *VersionController) GetStore() *VersionStore {
	return c.store
}

func (c *VersionController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	return c, nil
}

// TODO: This need to use the ui.NewListModel
func (c *VersionController) Render(cmd *cobra.Command, args []string) error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Version"), c.store.Version})
	tableData = append(tableData, []string{pterm.LightCyan("Date"), c.store.BuildDate})
	tableData = append(tableData, []string{pterm.LightCyan("Commit"), c.store.Commit})
	return pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}
