package profiles

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ProfilesShowStore struct {
	MembershipURI       string `json:"membership_uri"`
	DefaultOrganization string `json:"default_organization"`
}
type ProfilesShowController struct {
	store *ProfilesShowStore
}

var _ fctl.Controller[*ProfilesShowStore] = (*ProfilesShowController)(nil)

func NewDefaultProfilesShowStore() *ProfilesShowStore {
	return &ProfilesShowStore{
		MembershipURI:       "",
		DefaultOrganization: "",
	}
}

func NewProfilesShowController() *ProfilesShowController {
	return &ProfilesShowController{
		store: NewDefaultProfilesShowStore(),
	}
}

func (c *ProfilesShowController) GetStore() *ProfilesShowStore {
	return c.store
}

func (c *ProfilesShowController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	config, err := fctl.GetConfig(cmd)
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

func (c *ProfilesShowController) Render(cmd *cobra.Command, args []string) error {

	tableData := pterm.TableData{}
	tableData = append(tableData, []string{pterm.LightCyan("Membership URI"), c.store.MembershipURI})
	tableData = append(tableData, []string{pterm.LightCyan("Default organization"), c.store.DefaultOrganization})
	return pterm.DefaultTable.
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}

func NewShowCommand() *cobra.Command {
	return fctl.NewCommand("show <name>",
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithAliases("s"),
		fctl.WithShortDescription("Show profile"),
		fctl.WithValidArgsFunction(ProfileNamesAutoCompletion),
		fctl.WithController[*ProfilesShowStore](NewProfilesShowController()),
	)
}
