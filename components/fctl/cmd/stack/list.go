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

type Stack struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Dashboard string  `json:"dashboard"`
	RegionID  string  `json:"region"`
	DeletedAt *string `json:"deleted_at"`
}
type StackListStore struct {
	Stacks []Stack `json:"stacks"`
}

type StackListController struct {
	store   *StackListStore
	profile *fctl.Profile
}

var _ fctl.Controller[*StackListStore] = (*StackListController)(nil)

func NewDefaultStackListStore() *StackListStore {
	return &StackListStore{
		Stacks: []Stack{},
	}
}

func NewStackListController() *StackListController {
	return &StackListController{
		store: NewDefaultStackListStore(),
	}
}

func NewListCommand() *cobra.Command {
	return fctl.NewMembershipCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List stacks"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithBoolFlag(deletedFlag, false, "Display deleted stacks"),
		fctl.WithController[*StackListStore](NewStackListController()),
	)
}
func (c *StackListController) GetStore() *StackListStore {
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

	c.profile = profile
	if len(rsp.Data) == 0 {
		return c, nil
	}

	c.store.Stacks = fctl.Map(rsp.Data, func(stack membershipclient.Stack) Stack {
		return Stack{
			Id:        stack.Id,
			Name:      stack.Name,
			Dashboard: c.profile.ServicesBaseUrl(&stack).String(),
			RegionID:  stack.RegionID,
			DeletedAt: func() *string {
				if stack.DeletedAt != nil {
					t := stack.DeletedAt.Format(time.RFC3339)
					return &t
				}
				return nil
			}(),
		}
	})

	return c, nil
}

func (c *StackListController) Render(cmd *cobra.Command, args []string) error {
	if len(c.store.Stacks) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No stacks found.")
		return nil
	}

	tableData := fctl.Map(c.store.Stacks, func(stack Stack) []string {
		data := []string{
			stack.Id,
			stack.Name,
			stack.Dashboard,
			stack.RegionID,
		}
		if fctl.GetBool(cmd, deletedFlag) {
			if stack.DeletedAt != nil {
				data = append(data, *stack.DeletedAt)
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
