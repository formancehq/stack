package cmd

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	Version = "develop"
)

type VersionStruct struct {
	Version   string `json:"version" yaml:"version"`
	BuildDate string `json:"buildDate" yaml:"buildDate"`
	Commit    string `json:"commit" yaml:"commit"`
}
type VersionController struct {
	store *fctl.SharedStore
}

func NewVersion() *VersionController {
	return &VersionController{
		store: fctl.NewSharedStore(),
	}
}

func NewVersionCommand() *cobra.Command {
	return fctl.NewCommand("version",
		fctl.WithShortDescription("Get version"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController(NewVersion()),
	)
}

func (c *VersionController) GetStore() *fctl.SharedStore {
	return c.store
}

func (c *VersionController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	version := &VersionStruct{
		Version:   "develop",
		BuildDate: "-",
		Commit:    "-",
	}

	c.GetStore().SetData(version)

	return c, nil
}

// TODO: This need to use the ui.NewListModel
func (c *VersionController) Render(cmd *cobra.Command, args []string) error {
	data := c.GetStore().GetData().(*VersionStruct)

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Version"), data.Version})
	tableData = append(tableData, []string{pterm.LightCyan("Date"), data.BuildDate})
	tableData = append(tableData, []string{pterm.LightCyan("Commit"), data.Commit})
	return pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}
