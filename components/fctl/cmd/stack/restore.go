package stack

import (
	"fmt"
	"net/http"

	"github.com/formancehq/fctl/cmd/stack/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
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

	var stack *membershipclient.Stack
	if len(args) == 1 {
		rsp, _, err := apiClient.DefaultApi.GetStack(cmd.Context(), organization, args[0]).Execute()
		if err != nil {
			return nil, err
		}
		stack = rsp.Data
	}

	if stack == nil {
		return nil, errors.New("Stack not found")
	}

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to restore stack '%s'", stack.Name) {
		return nil, fctl.ErrMissingApproval
	}

	response, _, err := apiClient.DefaultApi.
		RestoreStack(cmd.Context(), organization, args[0]).
		Execute()
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(cmd, cfg)

	if !fctl.GetBool(cmd, nowaitFlag) {
		stack, err = waitStackReady(cmd, apiClient, profile, response.Data.OrganizationId, response.Data.Id)
		if err != nil {
			return nil, err
		}

		c.store.Stack = stack
	} else {
		c.store.Stack = response.Data
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

	c.store.Versions = versions.GetVersionsResponse
	c.config = cfg

	return c, nil
}

func (c *StackRestoreController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintStackInformation(cmd.OutOrStdout(), fctl.GetCurrentProfile(cmd, c.config), c.store.Stack, c.store.Versions)
}
