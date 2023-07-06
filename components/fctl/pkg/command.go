package fctl

import (
	"github.com/formancehq/fctl/membershipclient"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
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

func GetSelectedOrganization(cmd *cobra.Command) string {
	return GetString(cmd, organizationFlag)
}

func RetrieveOrganizationIDFromFlagOrProfile(cmd *cobra.Command, cfg *Config) (string, error) {
	if id := GetSelectedOrganization(cmd); id != "" {
		return id, nil
	}

	if defaultOrganization := GetCurrentProfile(cmd, cfg).GetDefaultOrganization(); defaultOrganization != "" {
		return defaultOrganization, nil
	}
	return "", ErrOrganizationNotSpecified
}

func ResolveOrganizationID(cmd *cobra.Command, cfg *Config) (string, error) {
	if id, err := RetrieveOrganizationIDFromFlagOrProfile(cmd, cfg); err == nil {
		return id, nil
	}

	client, err := NewMembershipClient(cmd, cfg)
	if err != nil {
		return "", err
	}

	organizations, _, err := client.DefaultApi.ListOrganizations(cmd.Context()).Execute()
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

func GetSelectedStackID(cmd *cobra.Command) string {
	return GetString(cmd, stackFlag)
}

func ResolveStack(cmd *cobra.Command, cfg *Config, organizationID string) (*membershipclient.Stack, error) {
	client, err := NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}
	if id := GetSelectedStackID(cmd); id != "" {
		response, _, err := client.DefaultApi.ReadStack(cmd.Context(), organizationID, id).Execute()
		if err != nil {
			return nil, err
		}

		return response.Data, nil
	}

	stacks, _, err := client.DefaultApi.ListStacks(cmd.Context(), organizationID).Execute()
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

func WrapOutputPostRunE(fn func(cmd *cobra.Command, args []string) error) CommandOptionFn {
	return func(cmd *cobra.Command) {
		cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
			flags := GetString(cmd, OutputFlag)

			if flags == "table" {
				return fn(cmd, args)
			}

			switch flags {
			case "json":
				// Marshal to JSON then print to stdout
				data, err := ShareStoreToJson()
				if (err) != nil {
					return err
				}
				pterm.Println(string(data))
				return nil
			}

			return cmd.Help()
		}
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

func WithConfirmFlag() CommandOptionFn {
	return WithBoolFlag(confirmFlag, false, "Confirm action")
}

func NewStackCommand(use string, opts ...CommandOption) *cobra.Command {
	cmd := NewMembershipCommand(use,
		append(opts,
			WithPersistentStringFlag(stackFlag, "", "Specific stack (not required if only one stack is present)"),
		)...,
	)
	cmd.RegisterFlagCompletionFunc("stack", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		cfg, err := GetConfig(cmd)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		profile := GetCurrentProfile(cmd, cfg)

		claims, err := profile.GetUserInfo(cmd)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		selectedOrganization := GetSelectedOrganization(cmd)
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

		cfg, err := GetConfig(cmd)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		profile := GetCurrentProfile(cmd, cfg)

		claims, err := profile.GetUserInfo(cmd)
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
			if GetBool(cmd, TelemetryFlag) {
				cfg, err := GetConfig(cmd)
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
