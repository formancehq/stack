package profiles

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ProfilesUseStore struct {
	Success bool `json:"success"`
}
type ProfilesUseController struct {
	store *ProfilesUseStore
}

var _ fctl.Controller[*ProfilesUseStore] = (*ProfilesUseController)(nil)

func NewDefaultProfilesUseStore() *ProfilesUseStore {
	return &ProfilesUseStore{
		Success: false,
	}
}

func NewProfilesUseController() *ProfilesUseController {
	return &ProfilesUseController{
		store: NewDefaultProfilesUseStore(),
	}
}

func ProfileNamesAutoCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ret, err := fctl.ListProfiles(cmd, toComplete)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveError
	}

	return ret, cobra.ShellCompDirectiveDefault
}

func (c *ProfilesUseController) GetStore() *ProfilesUseStore {
	return c.store
}

func (c *ProfilesUseController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	config, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	config.SetCurrentProfileName(args[0])
	if err := config.Persist(); err != nil {
		return nil, errors.Wrap(err, "Updating config")
	}

	c.store.Success = true
	return c, nil
}

func (c *ProfilesUseController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Selected profile updated!")
	return nil
}

func NewUseCommand() *cobra.Command {
	return fctl.NewCommand("use <name>",
		fctl.WithAliases("u"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithShortDescription("Use profile"),
		fctl.WithValidArgsFunction(ProfileNamesAutoCompletion),
		fctl.WithController[*ProfilesUseStore](NewProfilesUseController()),
	)
}
