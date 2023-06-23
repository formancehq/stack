package users

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ShowStore struct {
	User *User `json:"user,omitempty"`
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
	return fctl.NewCommand("show <user-id>",
		fctl.WithAliases("s"),
		fctl.WithShortDescription("Show user"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*ShowStore](NewShowController()),
	)
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(cmd, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(cmd, cfg, stack)
	if err != nil {
		return nil, err
	}

	request := operations.ReadUserRequest{
		UserID: args[0],
	}
	readUserResponse, err := client.Auth.ReadUser(cmd.Context(), request)
	if err != nil {
		return nil, err
	}

	if readUserResponse.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", readUserResponse.StatusCode)
	}

	c.store.User = &User{
		ID:      *readUserResponse.ReadUserResponse.Data.ID,
		Subject: *readUserResponse.ReadUserResponse.Data.Subject,
		Email:   *readUserResponse.ReadUserResponse.Data.Email,
	}

	return c, nil
}

func (c *ShowController) Render(cmd *cobra.Command, args []string) error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.User.ID})
	tableData = append(tableData, []string{pterm.LightCyan("Membership ID"), c.store.User.Subject})
	tableData = append(tableData, []string{pterm.LightCyan("Email"), c.store.User.Email})

	return pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()

}
