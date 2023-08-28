package stack

import (
	"flag"

	"github.com/formancehq/fctl/cmd/stack/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useDelete   = "delete (<stack-id> | --name=<stack-name>)"
	shortDelete = "Delete a stack"
)

type DeletedStore struct {
	Stack  *membershipclient.Stack `json:"stack"`
	Status string                  `json:"status"`
}

func NewDeletedStore() *DeletedStore {
	return &DeletedStore{
		Stack:  &membershipclient.Stack{},
		Status: "",
	}
}

func NewDeleteConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useDelete, flag.ExitOnError)
	flags.String(internal.StackNameFlag, "", "Stack to remove")
	fctl.WithConfirmFlag(flags)

	return fctl.NewControllerConfig(
		useDelete,
		shortDelete,
		shortDelete,
		[]string{
			"delete",
			"del",
			"rm",
		},
		flags,
		fctl.Organization,
	)
}

var _ fctl.Controller[*DeletedStore] = (*StackDeleteController)(nil)

type StackDeleteController struct {
	store  *DeletedStore
	config *fctl.ControllerConfig
}

func NewDeleteController(config *fctl.ControllerConfig) *StackDeleteController {
	return &StackDeleteController{
		store:  NewDeletedStore(),
		config: config,
	}
}

func (c *StackDeleteController) GetStore() *DeletedStore {
	return c.store
}

func (c *StackDeleteController) GetConfig() *fctl.ControllerConfig {
	return c.config
}

func (c *StackDeleteController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}
	organization, err := fctl.ResolveOrganizationID(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, errors.Wrap(err, "searching default organization")
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	var stack *membershipclient.Stack
	if len(c.config.GetArgs()) == 1 {
		if fctl.GetString(flags, internal.StackNameFlag) != "" {
			return nil, errors.New("need either an id of a name specified using --name flag")
		}

		rsp, _, err := apiClient.DefaultApi.ReadStack(ctx, organization, c.config.GetArgs()[0]).Execute()
		if err != nil {
			return nil, err
		}
		stack = rsp.Data
	} else {
		if fctl.GetString(flags, internal.StackNameFlag) == "" {
			return nil, errors.New("need either an id of a name specified using --name flag")
		}
		stacks, _, err := apiClient.DefaultApi.ListStacks(ctx, organization).Execute()
		if err != nil {
			return nil, errors.Wrap(err, "listing stacks")
		}
		for _, s := range stacks.Data {
			if s.Name == fctl.GetString(flags, internal.StackNameFlag) {
				stack = &s
				break
			}
		}
	}
	if stack == nil {
		return nil, errors.New("Stack not found")
	}

	if !fctl.CheckStackApprobation(flags, stack, "You are about to delete stack '%s'", stack.Name) {
		return nil, fctl.ErrMissingApproval
	}

	if _, err := apiClient.DefaultApi.DeleteStack(ctx, organization, stack.Id).Execute(); err != nil {
		return nil, errors.Wrap(err, "deleting stack")
	}

	c.store.Stack = stack
	c.store.Status = "OK"

	return c, nil
}

func (c *StackDeleteController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Stack deleted.")
	return nil
}

func NewDeleteCommand() *cobra.Command {
	config := NewDeleteConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.MaximumNArgs(1)),
		fctl.WithController[*DeletedStore](NewDeleteController(config)),
	)
}
