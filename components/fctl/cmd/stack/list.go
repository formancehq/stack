package stack

import (
	"fmt"
	"time"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	deletedFlag = "deleted"
)

type StackListController struct {
	store *fctl.SharedStore
}

func NewStackListController() *StackListController {
	return &StackListController{
		store: fctl.NewSharedStore(),
	}
}

func NewListCommand() *cobra.Command {
	return fctl.NewMembershipCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List stacks"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithBoolFlag(deletedFlag, false, "Display deleted stacks"),
		fctl.WithController(NewStackListController()),
		// fctl.WrapOutputPostRunE(view),
	)
}
func (c *StackListController) GetStore() *fctl.SharedStore {
	return c.store
}

func (c *StackListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(cmd, cfg)

	organization, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "searching default organization")
	}

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}

	rsp, _, err := apiClient.DefaultApi.ListStacks(cmd.Context(), organization).
		Deleted(fctl.GetBool(cmd, deletedFlag)).
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, "listing stacks")
	}

	if len(rsp.Data) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No stacks found.")
		return nil, nil
	}

	c.store.SetData(rsp.Data)
	c.store.SetProfile(profile)

	return c, nil
}

func (c *StackListController) Render(cmd *cobra.Command, args []string) error {

	sharedStruct := c.store.GetData().([]membershipclient.Stack)
	tableData := fctl.Map(sharedStruct, func(stack membershipclient.Stack) []string {
		data := []string{
			stack.Id,
			stack.Name,
			c.store.GetProfile().ServicesBaseUrl(&stack).String(),
			stack.RegionID,
		}
		if fctl.GetBool(cmd, deletedFlag) {
			if stack.DeletedAt != nil {
				data = append(data, stack.DeletedAt.Format(time.RFC3339))
			} else {
				data = append(data, "")
			}
		}
		return data
	})

	headers := []string{"ID", "Name", "Dashboard", "Region"}
	if fctl.GetBool(cmd, deletedFlag) {
		headers = append(headers, "Deleted at")
	}

	tableData = fctl.Prepend(tableData, headers)

	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}
