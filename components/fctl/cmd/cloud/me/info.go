package me

import (
	"errors"
	"flag"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useInfo   = "info"
	shortInfo = "Display user information"
)

type InfoStore struct {
	Subject string `json:"subject"`
	Email   string `json:"email"`
}

func NewInfoStore() *InfoStore {
	return &InfoStore{}
}

func NewInfoConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useInfo, flag.ExitOnError)
	return fctl.NewControllerConfig(
		useInfo,
		shortInfo,
		shortInfo,
		[]string{
			"i", "in",
		},
		flags,
	)
}

var _ fctl.Controller[*InfoStore] = (*InfoController)(nil)

type InfoController struct {
	store  *InfoStore
	config *fctl.ControllerConfig
}

func NewInfoController(config *fctl.ControllerConfig) *InfoController {
	return &InfoController{
		store:  NewInfoStore(),
		config: config,
	}
}

func (c *InfoController) GetStore() *InfoStore {
	return c.store
}

func (c *InfoController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *InfoController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(flags, cfg)
	if !profile.IsConnected() {
		return nil, errors.New("Not logged. Use 'login' command before.")
	}

	userInfo, err := profile.GetUserInfo()
	if err != nil {
		return nil, err
	}

	c.store.Subject = userInfo.Subject
	c.store.Email = userInfo.Email

	return c, nil
}

func (c *InfoController) Render() error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Subject"), c.store.Subject})
	tableData = append(tableData, []string{pterm.LightCyan("Email"), c.store.Email})

	return pterm.DefaultTable.
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()

}

func NewInfoCommand() *cobra.Command {
	config := NewInfoConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*InfoStore](NewInfoController(config)),
	)
}
