package stack

import (
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DeletedStack struct {
	Stack  *membershipclient.Stack `json:"stack"`
	Status string                  `json:"status"`
}

func NewDeleteCommand() *cobra.Command {
	const (
		stackNameFlag = "name"
	)
	return fctl.NewMembershipCommand("delete (<stack-id> | --name=<stack-name>)",
		fctl.WithConfirmFlag(),
		fctl.WithShortDescription("Delete a stack"),
		fctl.WithAliases("del", "d"),
		fctl.WithArgs(cobra.MaximumNArgs(1)),
		fctl.WithOutputFlag(),
		fctl.WithStringFlag(stackNameFlag, "", "Stack to remove"),
		fctl.WithRunE(deleteCommand),
		fctl.WrapOutputPostRunE(viewDelete),
	)
}

func deleteCommand(cmd *cobra.Command, args []string) error {
	const (
		stackNameFlag = "name"
	)

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return err
	}
	organization, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return errors.Wrap(err, "searching default organization")
	}

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return err
	}

	var stack *membershipclient.Stack
	if len(args) == 1 {
		if fctl.GetString(cmd, stackNameFlag) != "" {
			return errors.New("need either an id of a name specified using --name flag")
		}

		rsp, _, err := apiClient.DefaultApi.ReadStack(cmd.Context(), organization, args[0]).Execute()
		if err != nil {
			return err
		}
		stack = rsp.Data
	} else {
		if fctl.GetString(cmd, stackNameFlag) == "" {
			return errors.New("need either an id of a name specified using --name flag")
		}
		stacks, _, err := apiClient.DefaultApi.ListStacks(cmd.Context(), organization).Execute()
		if err != nil {
			return errors.Wrap(err, "listing stacks")
		}
		for _, s := range stacks.Data {
			if s.Name == fctl.GetString(cmd, stackNameFlag) {
				stack = &s
				break
			}
		}
	}
	if stack == nil {
		return errors.New("Stack not found")
	}

	if !fctl.CheckStackApprobation(cmd, stack, "You are about to delete stack '%s'", stack.Name) {
		return fctl.ErrMissingApproval
	}

	if _, err := apiClient.DefaultApi.DeleteStack(cmd.Context(), organization, stack.Id).Execute(); err != nil {
		return errors.Wrap(err, "deleting stack")
	}

	fctl.SetSharedData(&DeletedStack{
		Stack:  stack,
		Status: "OK",
	}, nil, nil)

	return nil
}

func viewDelete(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Stack deleted.")
	return nil
}
