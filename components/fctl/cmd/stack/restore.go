package stack

import (
	"fmt"
	"net/http"

	"github.com/formancehq/fctl/cmd/stack/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type StackRestoreStore struct {
	Stack    *membershipclient.Stack     `json:"stack"`
	Versions *shared.GetVersionsResponse `json:"versions"`
}
type StackRestoreController struct {
	store  *StackRestoreStore
	config *fctl.Config
}

var _ fctl.Controller[*StackRestoreStore] = (*StackRestoreController)(nil)

func NewDefaultVersionStore() *StackRestoreStore {
	return &StackRestoreStore{
		Stack:    &membershipclient.Stack{},
		Versions: &shared.GetVersionsResponse{},
	}
}

func NewStackRestoreController() *StackRestoreController {
	return &StackRestoreController{
		store: NewDefaultVersionStore(),
	}
}

func NewRestoreStackCommand() *cobra.Command {
	const stackNameFlag = "name"

	return fctl.NewMembershipCommand("restore <stack-id>",
		fctl.WithShortDescription("Restore a deleted stack"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithStringFlag(stackNameFlag, "", ""),
		fctl.WithController[*StackRestoreStore](NewStackRestoreController()),
	)
}
func (c *StackRestoreController) GetStore() *StackRestoreStore {
	return c.store
}

func (c *StackRestoreController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	organization, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "searching default organization")
	}

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}

	response, _, err := apiClient.DefaultApi.
		RestoreStack(cmd.Context(), organization, args[0]).
		Execute()
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(cmd, cfg)

	if err := waitStackReady(cmd, profile, response.Data); err != nil {
		return nil, err
	}

	stackClient, err := fctl.NewStackClient(cmd, cfg, response.Data)
	if err != nil {
		return nil, err
	}

	versions, err := stackClient.GetVersions(cmd.Context())
	if err != nil {
		return nil, err
	}

	if versions.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when reading versions", versions.StatusCode)
	}

	c.store.Stack = response.Data
	c.store.Versions = versions.GetVersionsResponse
	c.config = cfg

	return c, nil
}

func (c *StackRestoreController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintStackInformation(cmd.OutOrStdout(), fctl.GetCurrentProfile(cmd, c.config), c.store.Stack, c.store.Versions)
}
