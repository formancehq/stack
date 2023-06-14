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

func NewVersionCommand() *cobra.Command {
	return fctl.NewCommand("version",
		fctl.WithShortDescription("Get version"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithRunE(versionCommand),
		fctl.WrapOutputPostRunE(view),
	)
}

func versionCommand(cmd *cobra.Command, args []string) error {

	version := &VersionStruct{
		Version:   "develop",
		BuildDate: "-",
		Commit:    "-",
	}

	fctl.SetSharedData(version, nil, nil)

	return nil
}

func view(cmd *cobra.Command, args []string) error {
	data := fctl.GetSharedData().(*VersionStruct)

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Version"), data.Version})
	tableData = append(tableData, []string{pterm.LightCyan("Date"), data.BuildDate})
	tableData = append(tableData, []string{pterm.LightCyan("Commit"), data.Commit})
	return pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}
