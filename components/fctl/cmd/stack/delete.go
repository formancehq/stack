package stack

import (
	"context"
	"flag"
	"io"
	"os"

	"github.com/formancehq/fctl/cmd/stack/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useDelete         = "delete (<stack-id> | --name=<stack-name>)"
	descriptionDelete = "Delete a stack"
)

type DeletedStackStore struct {
	Stack  *membershipclient.Stack `json:"stack"`
	Status string                  `json:"status"`
}

func NewDefaultDeletedStackStore() *DeletedStackStore {
	return &DeletedStackStore{
		Stack:  &membershipclient.Stack{},
		Status: "",
	}
}

type StackDeleteControllerConfig struct {
	context     context.Context
	use         string
	description string
	aliases     []string
	out         io.Writer
	flags       *flag.FlagSet
	args        []string
}

func NewStackDeleteControllerConfig() *StackDeleteControllerConfig {
	flags := flag.NewFlagSet(useDelete, flag.ExitOnError)
	flags.String(internal.StackNameFlag, "", "Stack to remove")
	fctl.WithConfirmFlag(flags)
	fctl.WithGlobalFlags(flags)

	return &StackDeleteControllerConfig{
		context:     nil,
		use:         useDelete,
		description: descriptionDelete,
		aliases: []string{
			"del",
			"d",
		},
		out:   os.Stdout,
		flags: flags,
		args:  []string{},
	}
}

var _ fctl.Controller[*DeletedStackStore] = (*StackDeleteController)(nil)

type StackDeleteController struct {
	store  *DeletedStackStore
	config StackDeleteControllerConfig
}

func NewStackDeleteController(config StackDeleteControllerConfig) *StackDeleteController {
	return &StackDeleteController{
		store:  NewDefaultDeletedStackStore(),
		config: config,
	}
}
func (c *StackDeleteController) GetFlags() *flag.FlagSet {
	return c.config.flags
}

func (c *StackDeleteController) GetContext() context.Context {
	return c.config.context
}

func (c *StackDeleteController) SetContext(ctx context.Context) {
	c.config.context = ctx
}

func (c *StackDeleteController) GetStore() *DeletedStackStore {
	return c.store
}

func (c *StackDeleteController) SetArgs(args []string) {
	c.config.args = append([]string{}, args...)

}

func (c *StackDeleteController) Run() (fctl.Renderable, error) {
	flags := c.config.flags
	ctx := c.config.context

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}
	organization, err := fctl.ResolveOrganizationID(flags, ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "searching default organization")
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}

	var stack *membershipclient.Stack
	if len(c.config.args) == 1 {
		if fctl.GetString(flags, internal.StackNameFlag) != "" {
			return nil, errors.New("need either an id of a name specified using --name flag")
		}

		rsp, _, err := apiClient.DefaultApi.ReadStack(ctx, organization, c.config.args[0]).Execute()
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
	pterm.Success.WithWriter(c.config.out).Printfln("Stack deleted.")
	return nil
}

func NewDeleteCommand() *cobra.Command {
	config := NewStackDeleteControllerConfig()
	return fctl.NewMembershipCommand(config.use,
		fctl.WithShortDescription(config.description),
		fctl.WithAliases(config.aliases...),
		fctl.WithArgs(cobra.MaximumNArgs(1)),
		fctl.WithGoFlagSet(config.flags),
		fctl.WithController[*DeletedStackStore](NewStackDeleteController(*config)),
	)
}
