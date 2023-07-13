package stack

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	deletedFlag = "deleted"
	useList     = "list"
	description = "List stacks"
)

type Stack struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Dashboard string  `json:"dashboard"`
	RegionID  string  `json:"region"`
	DeletedAt *string `json:"deletedAt"`
}

type StackListStore struct {
	Stacks []Stack `json:"stacks"`
}

func NewDefaultStackListStore() *StackListStore {
	return &StackListStore{
		Stacks: []Stack{},
	}
}

type StackListControllerConfig struct {
	context     context.Context
	use         string
	description string
	aliases     []string
	out         io.Writer
	flags       *flag.FlagSet
	args        []string
}

func NewStackListControllerConfig() *StackListControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	flags.Bool(deletedFlag, false, "Show deleted stacks")

	fctl.WithGlobalFlags(flags)

	return &StackListControllerConfig{
		context:     nil,
		use:         useList,
		description: description,
		aliases: []string{
			"ls",
			"l",
		},
		out:   os.Stdout,
		flags: flags,
		args:  []string{},
	}
}

var _ fctl.Controller[*StackListStore] = (*StackListController)(nil)

type StackListController struct {
	store   *StackListStore
	profile *fctl.Profile
	config  StackListControllerConfig
}

func NewStackListController(config StackListControllerConfig) *StackListController {
	return &StackListController{
		store:  NewDefaultStackListStore(),
		config: config,
	}
}

func (c *StackListController) GetFlags() *flag.FlagSet {
	return c.config.flags
}

func (c *StackListController) GetContext() context.Context {
	return c.config.context
}

func (c *StackListController) SetContext(ctx context.Context) {
	c.config.context = ctx
}

func (c *StackListController) GetStore() *StackListStore {
	return c.store
}

func (c *StackListController) SetArgs(args []string) {
	c.config.args = append([]string{}, args...)
}

func (c *StackListController) Run() (fctl.Renderable, error) {
	flags := c.config.flags
	ctx := c.config.context

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(flags, cfg)

	organization, err := fctl.ResolveOrganizationID(flags, ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "searching default organization")
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}

	rsp, _, err := apiClient.DefaultApi.ListStacks(ctx, organization).
		Deleted(fctl.GetBool(flags, deletedFlag)).
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

func (c *StackListController) Render() error {
	if len(c.store.Stacks) == 0 {
		fmt.Fprintln(os.Stdout, "No stacks found.")
		return nil
	}

	tableData := fctl.Map(c.store.Stacks, func(stack Stack) []string {
		data := []string{
			stack.Id,
			stack.Name,
			stack.Dashboard,
			stack.RegionID,
		}
		if fctl.GetBool(c.config.flags, deletedFlag) {
			if stack.DeletedAt != nil {
				data = append(data, *stack.DeletedAt)
			} else {
				data = append(data, "")
			}
		}
		return data
	})

	headers := []string{"ID", "Name", "Dashboard", "Region"}
	if fctl.GetBool(c.config.flags, deletedFlag) {
		headers = append(headers, "Deleted at")
	}

	tableData = fctl.Prepend(tableData, headers)

	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(os.Stdout).
		WithData(tableData).
		Render()
}

func NewListCommand() *cobra.Command {
	config := NewStackListControllerConfig()

	options := []fctl.CommandOption{
		fctl.WithAliases(config.aliases...),
		fctl.WithShortDescription(config.description),
		fctl.WithArgs(cobra.ExactArgs(0)), //////////////// <--- This is used by cobra to validate the number of arguments passed to the command
		fctl.WithGoFlagSet(config.flags),
		fctl.WithController[*StackListStore](NewStackListController(*config)),
	}

	return fctl.NewMembershipCommand(config.use, options...)
}
