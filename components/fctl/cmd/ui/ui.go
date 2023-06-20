package ui

import (
	"fmt"
	"os/exec"
	"runtime"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

type UiStruct struct {
	UIUrl        string `json:"stack_url"`
	FoundBrowser bool   `json:"browser_found"`
}

type UiController struct {
	store *UiStruct
}

var _ fctl.Controller[*UiStruct] = (*UiController)(nil)

func NewDefaultUiStore() *UiStruct {
	return &UiStruct{
		UIUrl:        "",
		FoundBrowser: false,
	}
}

func NewUiController() *UiController {
	return &UiController{
		store: NewDefaultUiStore(),
	}
}

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
func (c *UiController) GetStore() *UiStruct {
	return c.store
}

func (c *UiController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	organization, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organization)
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(cmd, cfg)
	stackUrl := profile.ServicesBaseUrl(stack)

	c.store.UIUrl = stackUrl.String()

	if err := openUrl(c.store.UIUrl); err != nil {
		c.store.FoundBrowser = true
	}

	return c, nil
}

func (c *UiController) Render(cmd *cobra.Command, args []string) error {

	fmt.Println("Opening url: ", c.store.UIUrl)

	return nil
}

func NewCommand() *cobra.Command {
	return fctl.NewStackCommand("ui",
		fctl.WithShortDescription("Open UI"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*UiStruct](NewUiController()),
	)
}
