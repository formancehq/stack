package profiles

import (
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ProfilesDeleteStore struct {
	Success bool `json:"success"`
}
type ProfileDeleteController struct {
	store *ProfilesDeleteStore
}

var _ fctl.Controller[*ProfilesDeleteStore] = (*ProfileDeleteController)(nil)

func NewDefaultDeleteProfileStore() *ProfilesDeleteStore {
	return &ProfilesDeleteStore{
		Success: false,
	}
}

func NewDeleteProfileController() *ProfileDeleteController {
	return &ProfileDeleteController{
		store: NewDefaultDeleteProfileStore(),
	}
}

func (c *ProfileDeleteController) GetStore() *ProfilesDeleteStore {
	return c.store
}

func (c *ProfileDeleteController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	config, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}
	if err := config.DeleteProfile(args[0]); err != nil {
		return nil, err
	}

	if err := config.Persist(); err != nil {
		return nil, errors.Wrap(err, "updating config")
	}

	c.store.Success = true

	return c, nil
}

func (c *ProfileDeleteController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Profile deleted!")
	return nil
}

func NewDeleteCommand() *cobra.Command {
	return fctl.NewCommand("delete <name>",
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithShortDescription("Delete a profile"),
		fctl.WithValidArgsFunction(ProfileNamesAutoCompletion),
		fctl.WithController[*ProfilesDeleteStore](NewDeleteProfileController()),
	)
}
