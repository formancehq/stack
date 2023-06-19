package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func openUrl(url string) error {
	var (
		cmd  string
		args []string
	)

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func RunUiCommand(cmd *cobra.Command, args []string) error {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return err
	}

	organization, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organization)
	if err != nil {
		return err
	}

	profile := fctl.GetCurrentProfile(cmd, cfg)
	stackUrl := profile.ServicesBaseUrl(stack)

	ui := UiOutput{
		UIUrl:        stackUrl.String(),
		FoundBrowser: false,
	}

	if err := openUrl(ui.UIUrl); err != nil {
		ui.FoundBrowser = true
	}

	fctl.SetSharedData(ui, profile, cfg, nil)
	return errors.Wrap(err, "opening ui")

}

type UiOutput struct {
	UIUrl        string `json:"stack_url"`
	FoundBrowser bool   `json:"browser_found"`
}

func DisplayOutput(cmd *cobra.Command, args []string) error {

	ui, ok := fctl.GetSharedData().(UiOutput)
	if !ok {
		return errors.New("invalid output")
	}

	fmt.Println("Opening url: ", ui.UIUrl)

	return nil
}

func NewUICommand() *cobra.Command {
	return fctl.NewStackCommand("ui",
		fctl.WithShortDescription("Open UI"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithRunE(RunUiCommand),
		fctl.WrapOutputPostRunE(DisplayOutput),
	)
}
