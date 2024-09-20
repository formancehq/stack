package stack

import (
	"github.com/formancehq/fctl/cmd/stack/store"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	stackNameFlag = "name"
	forceFlag     = "force"
)

type DeletedStackStore struct {
	Stack  *membershipclient.Stack `json:"stack"`
	Status string                  `json:"status"`
}
type StackDeleteController struct {
	store *DeletedStackStore
}

var _ fctl.Controller[*DeletedStackStore] = (*StackDeleteController)(nil)

func NewDefaultDeletedStackStore() *DeletedStackStore {
	return &DeletedStackStore{
		Stack:  &membershipclient.Stack{},
		Status: "",
	}
}

func NewStackDeleteController() *StackDeleteController {
	return &StackDeleteController{
		store: NewDefaultDeletedStackStore(),
	}
}

func NewDeleteCommand() *cobra.Command {
	return fctl.NewMembershipCommand("delete (<stack-id> | --name=<stack-name>)",
		fctl.WithConfirmFlag(),
		fctl.WithShortDescription("Delete a stack"),
		fctl.WithAliases("del", "d"),
		fctl.WithArgs(cobra.MaximumNArgs(1)),
		fctl.WithStringFlag(stackNameFlag, "", "Stack to delete"),
		fctl.WithBoolFlag(forceFlag, false, "Force to delete a stack without retention policy"),
		fctl.WithController[*DeletedStackStore](NewStackDeleteController()),
	)
}
func (c *StackDeleteController) GetStore() *DeletedStackStore {
	return c.store
}

func (c *StackDeleteController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	store := store.GetStore(cmd.Context())

	var stack *membershipclient.Stack
	if len(args) == 1 {
		if fctl.GetString(cmd, stackNameFlag) != "" {
			return nil, errors.New("need either an id of a name specified using --name flag")
		}

		rsp, _, err := store.Client().GetStack(cmd.Context(), store.OrganizationId(), args[0]).Execute()
		if err != nil {
			return nil, err
		}
		stack = rsp.Data
	} else {
		if fctl.GetString(cmd, stackNameFlag) == "" {
			return nil, errors.New("need either an id of a name specified using --name flag")
		}
		stacks, _, err := store.Client().ListStacks(cmd.Context(), store.OrganizationId()).Execute()
		if err != nil {
			return nil, errors.Wrap(err, "listing stacks")
		}
		for _, s := range stacks.Data {
			if s.Name == fctl.GetString(cmd, stackNameFlag) {
				stack = &s
				break
			}
		}
	}
	if stack == nil {
		return nil, errors.New("Stack not found")
	}

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to delete stack '%s'", stack.Name) {
		return nil, fctl.ErrMissingApproval
	}

	query := store.Client().DeleteStack(cmd.Context(), store.OrganizationId(), stack.Id)
	if fctl.GetBool(cmd, forceFlag) {
		if isValid := fctl.CheckMembershipVersion("v0.27.1")(cmd, args); isValid != nil {
			return nil, isValid
		}
		query = query.Force(true)
	}

	_, err := query.Execute()
	if err != nil {
		return nil, errors.Wrap(err, "deleting stack")
	}

	c.store.Stack = stack
	c.store.Status = "OK"

	return c, nil
}

func (c *StackDeleteController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Stack deleted.")
	return nil
}
