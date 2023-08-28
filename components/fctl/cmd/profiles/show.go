package profiles

import (
	"flag"

	"github.com/formancehq/fctl/cmd/profiles/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useShow         = "show <name>"
	shortShow       = "Show a profile"
	descriptionShow = "Show a profile"
)

type ShowStore struct {
	MembershipURI       string `json:"membershipUri"`
	DefaultOrganization string `json:"defaultOrganization"`
}

func NewShowStore() *ShowStore {
	return &ShowStore{
		MembershipURI:       "",
		DefaultOrganization: "",
	}
}

func NewShowConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useShow, flag.ExitOnError)

	return fctl.NewControllerConfig(
		useShow,
		descriptionShow,
		shortShow,
		[]string{
			"s",
		},
		flags,
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

	args := c.config.GetArgs()
	if len(args) < 1 {
		return nil, errors.New("Profile: invalid number of arguments")
	}

	config, err := fctl.GetConfig(c.config.GetAllFLags())
	if err != nil {
		return nil, err
	}

	p := config.GetProfile(args[0])
	if p == nil {
		return nil, errors.New("not found")
	}

	c.store.DefaultOrganization = p.GetDefaultOrganization()
	c.store.MembershipURI = p.GetMembershipURI()

	return c, nil
}

func (c *ShowController) Render() error {

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Membership URI"), c.store.MembershipURI})
	tableData = append(tableData, []string{pterm.LightCyan("Default organization"), c.store.DefaultOrganization})
	return pterm.DefaultTable.
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()
}

func NewShowCommand() *cobra.Command {
	config := NewShowConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithValidArgsFunction(internal.ProfileCobraAutoCompletion),
		fctl.WithController[*ShowStore](NewShowController(config)),
	)
}
