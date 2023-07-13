package fctl

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/TylerBrock/colorjson"
	"github.com/formancehq/fctl/membershipclient"
	"github.com/pkg/errors"
	"github.com/segmentio/analytics-go/v3"
	"github.com/segmentio/ksuid"
	"github.com/spf13/cobra"
)

const (
	stackFlag              = "stack"
	organizationFlag       = "organization"
	DefaultSegmentWriteKey = ""
	outputFlag             = "output"
)

var (
	ErrOrganizationNotSpecified   = errors.New("organization not specified")
	ErrMultipleOrganizationsFound = errors.New("found more than one organization and no organization specified")
)

type StackOrganizationConfig struct {
	OrganizationID string
	Stack          *membershipclient.Stack
	Config         *Config
}

func GetStackOrganizationConfig(flags *flag.FlagSet, ctx context.Context) (*StackOrganizationConfig, error) {
	cfg, err := GetConfig(flags)
	if err != nil {
		return nil, err
	}

	organizationID, err := ResolveOrganizationID(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := ResolveStack(flags, ctx, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	return &StackOrganizationConfig{
		OrganizationID: organizationID,
		Stack:          stack,
		Config:         cfg,
	}, nil
}

func GetStackOrganizationConfigApprobation(flags *flag.FlagSet, ctx context.Context, disclaimer string, args ...any) (*StackOrganizationConfig, error) {
	soc, err := GetStackOrganizationConfig(flags, ctx)
	if err != nil {
		return nil, err
	}

	if !CheckStackApprobation(flags, soc.Stack, disclaimer, args...) {
		return nil, ErrMissingApproval
	}

	return soc, nil
}

func GetSelectedOrganization(flags *flag.FlagSet) string {
	return GetString(flags, organizationFlag)
}

func RetrieveOrganizationIDFromFlagOrProfile(flags *flag.FlagSet, cfg *Config) (string, error) {
	if id := GetSelectedOrganization(flags); id != "" {
		return id, nil
	}

	if defaultOrganization := GetCurrentProfile(flags, cfg).GetDefaultOrganization(); defaultOrganization != "" {
		return defaultOrganization, nil
	}
	return "", ErrOrganizationNotSpecified
}

func ResolveOrganizationID(flags *flag.FlagSet, ctx context.Context, cfg *Config) (string, error) {
	if id, err := RetrieveOrganizationIDFromFlagOrProfile(flags, cfg); err == nil {
		return id, nil
	}

	client, err := NewMembershipClient(flags, ctx, cfg)
	if err != nil {
		return "", err
	}

	organizations, _, err := client.DefaultApi.ListOrganizations(ctx).Execute()
	if err != nil {
		return "", errors.Wrap(err, "listing organizations")
	}

	if len(organizations.Data) == 0 {
		return "", errors.New("no organizations found")
	}

	if len(organizations.Data) > 1 {
		return "", ErrMultipleOrganizationsFound
	}

	return organizations.Data[0].Id, nil
}

func GetSelectedStackID(flags *flag.FlagSet) string {
	return GetString(flags, stackFlag)
}

func ResolveStack(flags *flag.FlagSet, ctx context.Context, cfg *Config, organizationID string) (*membershipclient.Stack, error) {
	client, err := NewMembershipClient(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}
	if id := GetSelectedStackID(flags); id != "" {
		response, _, err := client.DefaultApi.ReadStack(ctx, organizationID, id).Execute()
		if err != nil {
			return nil, err
		}

		return response.Data, nil
	}

	stacks, _, err := client.DefaultApi.ListStacks(ctx, organizationID).Execute()
	if err != nil {
		return nil, errors.Wrap(err, "listing stacks")
	}
	if len(stacks.Data) == 0 {
		return nil, errors.New("no stacks found")
	}
	if len(stacks.Data) > 1 {
		return nil, errors.New("found more than one stack and no stack specified")
	}
	return &(stacks.Data[0]), nil
}

type CommandOption interface {
	apply(cmd *cobra.Command)
}
type CommandOptionFn func(cmd *cobra.Command)

func (fn CommandOptionFn) apply(cmd *cobra.Command) {
	fn(cmd)
}

func WithGoFlagSet(flags *flag.FlagSet) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.Flags().AddGoFlagSet(flags)
	}
}

func WithPersistentStringFlag(name, defaultValue, help string) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.PersistentFlags().String(name, defaultValue, help)
	}
}

func WithStringFlag(name, defaultValue, help string) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.Flags().String(name, defaultValue, help)
	}
}

func WithPersistentStringPFlag(name, short, defaultValue, help string) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.PersistentFlags().StringP(name, short, defaultValue, help)
	}
}

func WithBoolFlag(name string, defaultValue bool, help string) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.Flags().Bool(name, defaultValue, help)
	}
}

func WithAliases(aliases ...string) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.Aliases = aliases
	}
}

func WithPersistentBoolPFlag(name, short string, defaultValue bool, help string) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolP(name, short, defaultValue, help)
	}
}

func WithPersistentBoolFlag(name string, defaultValue bool, help string) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.PersistentFlags().Bool(name, defaultValue, help)
	}
}

func WithIntFlag(name string, defaultValue int, help string) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.Flags().Int(name, defaultValue, help)
	}
}

func WithStringSliceFlag(name string, defaultValue []string, help string) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.Flags().StringSlice(name, defaultValue, help)
	}
}

func WithHiddenFlag(name string) CommandOptionFn {
	return func(cmd *cobra.Command) {
		_ = cmd.Flags().MarkHidden(name)
	}
}

func WithHidden() CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.Hidden = true
	}
}

func WithRunE(fn func(cmd *cobra.Command, args []string) error) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.RunE = fn
	}
}

func WithPreRunE(fn func(cmd *cobra.Command, args []string) error) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.PreRunE = fn
	}
}

func WithGlobalFlags(flags *flag.FlagSet) *flag.FlagSet {

	if flags == nil {
		flags = flag.NewFlagSet("global", flag.ContinueOnError)
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	flags.Bool(InsecureTlsFlag, false, "Insecure TLS")
	flags.Bool(TelemetryFlag, false, "Telemetry enabled")
	flags.String(ProfileFlag, "", "config profile to use")
	flags.String(FileFlag, fmt.Sprintf("%s/.formance/fctl.config", homedir), "Debug mode")
	flags.Bool(DebugFlag, false, "Debug mode")
	flags.String(outputFlag, "plain", "Output format (plain, json)")

	return flags
}

func WithController[T any](c Controller[T]) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
			if c.GetContext() == nil {
				c.SetContext(cmd.Context())
			}

			if len(args) > 0 {
				c.SetArgs(args)
			}

			return nil
		}
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			renderable, err := c.Run()

			// If the controller return an argument error, we want to print the usage
			// of the command instead of the error message.
			// if errors.Is(err, ErrArgument) {
			// 	_ = cmd.help()
			// 	return nil
			// }

			if err != nil {
				return err
			}

			err = WithRender(c.GetFlags(), args, c, renderable)

			if err != nil {
				return err
			}

			return nil
		}
	}
}
func WithRender[T any](flags *flag.FlagSet, args []string, c Controller[T], r Renderable) error {
	flag := GetString(flags, OutputFlag)

	switch flag {
	case "json":
		// Inject into export struct
		export := ExportedData{
			Data: c.GetStore(),
		}

		// Marshal to JSON then print to stdout
		out, err := json.Marshal(export)
		if err != nil {
			return err
		}

		raw := make(map[string]any)
		if err := json.Unmarshal(out, &raw); err == nil {
			f := colorjson.NewFormatter()
			f.Indent = 2
			colorized, err := f.Marshal(raw)
			if err != nil {
				panic(err)
			}
			fmt.Print(string(colorized))
			return nil
		} else {
			fmt.Print(out)
			return nil
		}
	default:
		return r.Render()
	}
}

func WithChildCommands(cmds ...*cobra.Command) CommandOptionFn {
	return func(cmd *cobra.Command) {
		for _, child := range cmds {
			cmd.AddCommand(child)
		}
	}
}

func WithShortDescription(v string) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.Short = v
	}
}

func WithArgs(p cobra.PositionalArgs) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.Args = p
	}
}

func WithValidArgs(validArgs ...string) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.ValidArgs = validArgs
	}
}

func WithValidArgsFunction(fn func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.ValidArgsFunction = fn
	}
}

func WithDescription(v string) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.Long = v
	}
}

func WithSilenceUsage() CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.SilenceUsage = true
	}
}

func WithSilenceError() CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.SilenceErrors = true
	}
}

func NewStackCommand(use string, opts ...CommandOption) *cobra.Command {
	cmd := NewMembershipCommand(use,
		append(opts,
			WithPersistentStringFlag(stackFlag, "", "Specific stack (not required if only one stack is present)"),
		)...,
	)
	cmd.RegisterFlagCompletionFunc("stack", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		flags := ConvertPFlagSetToFlagSet(cmd.Flags())

		cfg, err := GetConfig(flags)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		profile := GetCurrentProfile(flags, cfg)

		claims, err := profile.GetUserInfo()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		selectedOrganization := GetSelectedOrganization(flags)
		if selectedOrganization == "" {
			selectedOrganization = profile.defaultOrganization
		}

		ret := make([]string, 0)
		for _, org := range claims.Org {
			if selectedOrganization != "" && selectedOrganization != org.ID {
				continue
			}
			for _, stack := range org.Stacks {
				ret = append(ret, fmt.Sprintf("%s\t%s", stack.ID, stack.DisplayName))
			}
		}

		return ret, cobra.ShellCompDirectiveDefault
	})
	return cmd
}

func NewMembershipCommand(use string, opts ...CommandOption) *cobra.Command {
	cmd := NewCommand(use,
		append(opts,
			WithPersistentStringFlag(organizationFlag, "", "Selected organization (not required if only one organization is present)"),
		)...,
	)
	cmd.RegisterFlagCompletionFunc("organization", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		flags := ConvertPFlagSetToFlagSet(cmd.Flags())
		cfg, err := GetConfig(flags)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		profile := GetCurrentProfile(flags, cfg)

		claims, err := profile.GetUserInfo()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		ret := make([]string, 0)
		for _, org := range claims.Org {
			ret = append(ret, fmt.Sprintf("%s\t%s", org.ID, org.DisplayName))
		}

		return ret, cobra.ShellCompDirectiveDefault
	})
	return cmd
}

func NewCommand(use string, opts ...CommandOption) *cobra.Command {
	cmd := &cobra.Command{
		Use: use,
		PersistentPostRun: func(cmd *cobra.Command, args []string) {

			flags := ConvertPFlagSetToFlagSet(cmd.Flags())

			if GetBool(flags, TelemetryFlag) {
				cfg, err := GetConfig(flags)
				if err != nil {
					return
				}

				if cfg.GetUniqueID() == "" {
					uniqueID := ksuid.New().String()
					cfg.SetUniqueID(uniqueID)
					err = cfg.Persist()
					if err != nil {
						return
					}
				}

				client := NewAnalyticsClient()
				defer client.Close()
				err = client.Enqueue(analytics.Track{
					UserId: cfg.GetUniqueID(),
					Event:  "fctl usages",
					Properties: analytics.NewProperties().
						Set("command", cmd.Name()).
						Set("args", args),
				})
				if err != nil {
					return
				}
			}
		},
	}
	for _, opt := range opts {
		opt.apply(cmd)
	}
	return cmd
}

func NewAnalyticsClient() analytics.Client {
	client := analytics.New(DefaultSegmentWriteKey)
	return client
}
