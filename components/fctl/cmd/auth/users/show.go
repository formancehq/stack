package users

import (
	"flag"
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useShow   = "show <user-id>"
	shortShow = "Show user details"
)

type ShowStore struct {
	User *User `json:"user,omitempty"`
}

func NewShowStore() *ShowStore {
	return &ShowStore{}
}

func NewShowConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useShow, flag.ExitOnError)
	return fctl.NewControllerConfig(
		useShow,
		shortShow,
		shortShow,
		[]string{
			"sh",
		},
		flags,
		fctl.Organization, fctl.Stack,
	)
}

var _ fctl.Controller[*ShowStore] = (*ShowController)(nil)

type ShowController struct {
	store  *ShowStore
	config *fctl.ControllerConfig
}

func NewShowController(config *fctl.ControllerConfig) *ShowController {
	return &ShowController{
		store:  NewShowStore(),
		config: config,
	}
}

func (c *ShowController) GetStore() *ShowStore {
	return c.store
}

func (c *ShowController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *ShowController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()
	out := c.config.GetOut()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, err
	}

	request := operations.ReadUserRequest{
		UserID: args[0],
	}
	readUserResponse, err := client.Auth.ReadUser(ctx, request)
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

func (c *ShowController) Render() error {
	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("ID"), c.store.User.ID})
	tableData = append(tableData, []string{pterm.LightCyan("Membership ID"), c.store.User.Subject})
	tableData = append(tableData, []string{pterm.LightCyan("Email"), c.store.User.Email})

	return pterm.DefaultTable.
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()

}

func NewShowCommand() *cobra.Command {

	config := NewShowConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*ShowStore](NewShowController(config)),
	)
}
