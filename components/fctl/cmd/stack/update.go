package stack

import (
	"fmt"

	"github.com/formancehq/fctl/cmd/stack/internal"
	"github.com/formancehq/fctl/cmd/stack/store"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

const (
	nameFlag = "name"
)

type StackUpdateStore struct {
	Stack *membershipclient.Stack
}

type StackUpdateController struct {
	store   *StackUpdateStore
	profile *fctl.Profile
}

var _ fctl.Controller[*StackUpdateStore] = (*StackUpdateController)(nil)

func NewDefaultStackUpdateStore() *StackUpdateStore {
	return &StackUpdateStore{
		Stack: &membershipclient.Stack{},
	}
}
func NewStackUpdateController() *StackUpdateController {
	return &StackUpdateController{
		store: NewDefaultStackUpdateStore(),
	}
}

func NewUpdateCommand() *cobra.Command {
	return fctl.NewMembershipCommand("update <stack-id>",
		fctl.WithShortDescription("Update a created stack, name, or metadata"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithPreRunE(func(cmd *cobra.Command, args []string) error {
			store := store.GetStore(cmd.Context())

			version := fctl.MembershipServerInfo(cmd.Context(), store.Client())
			if !semver.IsValid(version) {
				return nil
			}

			if semver.Compare(version, "v0.27.1") >= 0 {
				return nil
			}

			return fmt.Errorf("unsupported membership server version: %s", version)
		}),
		fctl.WithBoolFlag(unprotectFlag, false, "Unprotect stacks (no confirmation on write commands)"),
		fctl.WithStringFlag(nameFlag, "", "Name of the stack"),
		fctl.WithController[*StackUpdateStore](NewStackUpdateController()),
	)
}
func (c *StackUpdateController) GetStore() *StackUpdateStore {
	return c.store
}

func (c *StackUpdateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := store.GetStore(cmd.Context())
	c.profile = store.Config.GetProfile(fctl.GetCurrentProfileName(cmd, store.Config))

	protected := !fctl.GetBool(cmd, unprotectFlag)
	metadata := map[string]string{
		fctl.ProtectedStackMetadata: fctl.BoolPointerToString(&protected),
	}

	stack, res, err := store.Client().GetStack(cmd.Context(), store.OrganizationId(), args[0]).Execute()
	if err != nil {
		return nil, errors.Wrap(err, "retrieving stack")
	}
	if res.StatusCode > 300 {
		return nil, errors.New("stack not found")
	}

	name := fctl.GetString(cmd, nameFlag)
	if name == "" {
		name = stack.Data.Name
	}

	req := membershipclient.UpdateStackRequest{
		Name:     name,
		Metadata: metadata,
	}

	stackResponse, _, err := store.Client().
		UpdateStack(cmd.Context(), store.OrganizationId(), args[0]).
		UpdateStackRequest(req).
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, "updating stack")
	}

	c.store.Stack = stackResponse.Data

	return c, nil
}

func (c *StackUpdateController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintStackInformation(cmd.OutOrStdout(), c.profile, c.store.Stack, nil)
}
