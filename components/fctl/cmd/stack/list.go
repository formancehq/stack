package stack

import (
	"fmt"
	"time"

	"github.com/formancehq/fctl/cmd/stack/store"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	allFlag     = "all"
	deletedFlag = "deleted"
)

type Stack struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Dashboard    string  `json:"dashboard"`
	RegionID     string  `json:"region"`
	DisabledAt   *string `json:"disabledAt"`
	DeletedAt    *string `json:"deletedAt"`
	AuditEnabled string  `json:"auditEnabled"`
	Status       string  `json:"status"`
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
		fctl.WithBoolFlag(allFlag, false, "Display deleted stacks"),
		fctl.WithDeprecatedFlag(deletedFlag, "Use --all instead"),
		fctl.WithController[*StackListStore](NewStackListController()),
	)
}
func (c *StackListController) GetStore() *StackListStore {
	return c.store
}

func (c *StackListController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	store := store.GetStore(cmd.Context())
	c.profile = store.Config.GetProfile(fctl.GetCurrentProfileName(cmd, store.Config))

	rsp, _, err := store.Client().ListStacks(cmd.Context(), store.OrganizationId()).
		All(fctl.GetBool(cmd, allFlag)).
		Deleted(fctl.GetBool(cmd, deletedFlag)).
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, "listing stacks")
	}

	if len(rsp.Data) == 0 {
		return c, nil
	}

	c.store.Stacks = fctl.Map(rsp.Data, func(stack membershipclient.Stack) Stack {
		return Stack{
			Id:           stack.Id,
			Name:         stack.Name,
			Dashboard:    "https://console.formance.cloud",
			RegionID:     stack.RegionID,
			Status:       stack.State,
			AuditEnabled: fctl.BoolPointerToString(stack.AuditEnabled),
			DisabledAt: func() *string {
				if stack.DisabledAt != nil {
					t := stack.DisabledAt.Format(time.RFC3339)
					return &t
				}
				return nil
			}(),
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
			stack.Status,
			stack.AuditEnabled,
		}
		if fctl.GetBool(cmd, allFlag) {
			if stack.DisabledAt != nil {
				data = append(data, *stack.DisabledAt)
			} else {
				data = append(data, "")
			}

			if stack.DeletedAt != nil {
				data = append(data, *stack.DeletedAt)
			} else {
				if stack.Status != "DELETED" {
					data = append(data, "")
				} else {
					data = append(data, "<retention period>")
				}
			}
		}

		return data
	})

	headers := []string{"ID", "Name", "Dashboard", "Region", "Status", "Audit Enabled"}
	if fctl.GetBool(cmd, allFlag) {
		headers = append(headers, "Disabled At", "Deleted At")
	}
	tableData = fctl.Prepend(tableData, headers)

	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(cmd.OutOrStdout()).
		WithData(tableData).
		Render()
}
