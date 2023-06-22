package me

import (
	"errors"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type InfoStore struct {
	Subject string `json:"subject"`
	Email   string `json:"email"`
}
type InfoController struct {
	store *InfoStore
}

var _ fctl.Controller[*InfoStore] = (*InfoController)(nil)

func NewDefaultInfoStore() *InfoStore {
	return &InfoStore{}
}

func NewInfoController() *InfoController {
	return &InfoController{
		store: NewDefaultInfoStore(),
	}
}

func NewInfoCommand() *cobra.Command {
	return fctl.NewCommand("info",
		fctl.WithAliases("i", "in"),
		fctl.WithShortDescription("Display user information"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*InfoStore](NewInfoController()),
	)
}

func (c *InfoController) GetStore() *InfoStore {
	return c.store
}

func (c *InfoController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(cmd, cfg)
	if !profile.IsConnected() {
		return nil, errors.New("Not logged. Use 'login' command before.")
	}

	userInfo, err := profile.GetUserInfo(cmd)
	if err != nil {
		return nil, err
	}

	c.store.Subject = userInfo.Subject
	c.store.Email = userInfo.Email

	return c, nil
}

func (c *InfoController) Render(cmd *cobra.Command, args []string) error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Subject"), c.store.Subject})
	tableData = append(tableData, []string{pterm.LightCyan("Email"), c.store.Email})

	return pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()

}
