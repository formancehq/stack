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

type RestoreStackInformation struct {
	Stack    *membershipclient.Stack     `json:"stack"`
	Versions *shared.GetVersionsResponse `json:"versions"`
}
type StackRestore struct {
	store *fctl.SharedStore
}

func NewStackRestoreController() *StackRestore {
	return &StackRestore{
		store: fctl.NewSharedStore(),
	}
}

func NewRestoreStackCommand() *cobra.Command {
	const stackNameFlag = "name"

	return fctl.NewMembershipCommand("restore <stack-id>",
		fctl.WithShortDescription("Restore a deleted stack"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithStringFlag(stackNameFlag, "", ""),
		fctl.WithController(NewStackRestoreController()),
	)
}
func (c *StackRestore) GetStore() *fctl.SharedStore {
	return c.store
}

func (c *StackRestore) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
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

	c.store.SetData(&RestoreStackInformation{
		Stack:    response.Data,
		Versions: versions.GetVersionsResponse,
	})

	c.store.SetConfig(cfg)

	return c, nil
}

func (c *StackRestore) Render(cmd *cobra.Command, args []string) error {
	data := c.store.GetData().(*RestoreStackInformation)

	return internal.PrintStackInformation(cmd.OutOrStdout(), fctl.GetCurrentProfile(cmd, c.store.GetConfig()), data.Stack, data.Versions)
}
