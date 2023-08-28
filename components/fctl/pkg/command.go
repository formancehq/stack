package fctl

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strconv"

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

func GetStackOrganizationConfig(flags *flag.FlagSet, ctx context.Context, out io.Writer) (*StackOrganizationConfig, error) {
	cfg, err := GetConfig(flags)
	if err != nil {
		return nil, err
	}

	organizationID, err := ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	return &StackOrganizationConfig{
		OrganizationID: organizationID,
		Stack:          stack,
		Config:         cfg,
	}, nil
}

func GetStackOrganizationConfigApprobation(flags *flag.FlagSet, ctx context.Context, disclaimer string, out io.Writer, args ...any) (*StackOrganizationConfig, error) {
	soc, err := GetStackOrganizationConfig(flags, ctx, out)
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

func ResolveOrganizationID(flags *flag.FlagSet, ctx context.Context, cfg *Config, out io.Writer) (string, error) {
	if id, err := RetrieveOrganizationIDFromFlagOrProfile(flags, cfg); err == nil {
		return id, nil
	}

	client, err := NewMembershipClient(flags, ctx, cfg, out)
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

func ResolveStack(flags *flag.FlagSet, ctx context.Context, cfg *Config, organizationID string, out io.Writer) (*membershipclient.Stack, error) {
	client, err := NewMembershipClient(flags, ctx, cfg, out)
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

// WithCommandScopesFlags adds flags to the command that will be used to
// set the scopes only cobra side for display purpose.
// The function is used to display the scopes in the help
func WithCommandScopesFlags(flags ...*flag.Flag) CommandOptionFn {
	return func(cmd *cobra.Command) {
		for _, f := range flags {
			cmd.PersistentFlags().AddGoFlag(f)
		}
	}
}

// WithGoPersistentFlagSet is intended to be used only with fValue[T] flags who are defined in pkg/scopes.go
// AND WithGoPersistentFlagSet is used to define short flag for scpecific flag only
func WithGoPersistentFlagSet(flags *flag.FlagSet) CommandOptionFn {
	return func(cmd *cobra.Command) {
		flags.VisitAll(func(f *flag.Flag) {

			switch f.Name {
			case "config":
				cmd.PersistentFlags().StringVarP(&configFlagV.value, f.Name, "c", f.DefValue, f.Usage)
			case "profile":
				cmd.PersistentFlags().StringVarP(&profileFlagV.value, f.Name, "p", f.DefValue, f.Usage)
			case "output":
				cmd.PersistentFlags().StringVarP(&outputFlagV.value, f.Name, "o", f.DefValue, f.Usage)
			case "debug":
				defaultV, err := strconv.ParseBool(f.DefValue)
				if err != nil {
					panic(err)
				}
				cmd.PersistentFlags().BoolVarP(&debugFlagV.value, f.Name, "d", defaultV, f.Usage)
			}
			cmd.PersistentFlags().AddGoFlag(f)

		})

	}
}

func WithAliases(aliases ...string) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.Aliases = aliases
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

func WithDeprecated(message string) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.Deprecated = message
	}
}

func withOrganizationCompletion() CommandOptionFn {
	return func(cmd *cobra.Command) {
		err := cmd.RegisterFlagCompletionFunc(organizationFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			flags := ConvertPFlagSetToFlagSet(cmd.PersistentFlags())
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
		if err != nil {
			panic(errors.Wrap(err, "failed to register flag completion function"))
		}

	}
}

func withStackCompletion() CommandOptionFn {
	return func(cmd *cobra.Command) {
		err := cmd.RegisterFlagCompletionFunc("stack", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

		if err != nil {
			panic(errors.Wrap(err, "could not register stack flag completion"))
		}

	}
}

// We dont want to add again the pFlagSet, because it is already added by the root command
func configureCobraWithControllerConfig(cmd *cobra.Command, config *ControllerConfig) *cobra.Command {
	config.SetOut(cmd.OutOrStdout())
	cmd.Aliases = append(cmd.Aliases, config.GetAliases()...)
	cmd.Use = config.GetUse()
	cmd.Short = config.GetShortDescription()
	cmd.Long = config.GetDescription()

	//Adding flags
	scopes := config.GetScopes()
	cmd.Flags().AddGoFlagSet(scopes)
	cmd.Flags().AddGoFlagSet(config.GetFlags())

	//Adding completion flags
	if cmd.Flags().Lookup(stackFlag) != nil {
		withStackCompletion()(cmd)
	}

	if cmd.Flags().Lookup(organizationFlag) != nil {
		withOrganizationCompletion()(cmd)
	}

	//Hide Specific flags
	if cmd.Flags().Lookup("metadata") != nil {
		err := cmd.Flags().MarkHidden("metadata")
		if err != nil {
			panic(errors.Wrap(err, "failed to hide metadata flag"))
		}
	}

	return cmd
}
func WithController[T any](c Controller[T]) CommandOptionFn {
	return func(cmd *cobra.Command) {
		config := c.GetConfig()
		cmd = configureCobraWithControllerConfig(cmd, config)
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			if config.GetContext() == nil {
				config.SetContext(cmd.Context())
			}

			config.SetArgs(args)

			renderer, err := c.Run()

			// If the controller return an argument error, we want to print the usage
			// of the command instead of the error message.
			// if errors.Is(err, ErrArgument) {
			// 	_ = cmd.help()
			// 	return nil
			// }

			if err != nil {
				return err
			}

			err = render(config.GetPFlags(), c, renderer)

			if err != nil {
				return err
			}

			return nil
		}
	}
}
func render[T any](flags *flag.FlagSet, c Controller[T], r Renderable) error {
	f := GetString(flags, OutputFlag)
	getConfig := c.GetConfig()
	outWriter := getConfig.GetOut()
	switch f {
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
			_, err = fmt.Fprintln(outWriter, string(colorized))
			if err != nil {
				return err
			}
			return nil
		} else {
			_, err := fmt.Fprintln(outWriter, out)
			if err != nil {
				return err
			}
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
