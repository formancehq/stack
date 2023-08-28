package internal

import (
	"flag"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

func ProfileNamesAutoCompletion(flags *flag.FlagSet, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ret, err := fctl.ListProfiles(flags, toComplete)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveError
	}

	return ret, cobra.ShellCompDirectiveDefault
}

func ProfileCobraAutoCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	flags := fctl.ConvertPFlagSetToFlagSet(cmd.Flags())
	return ProfileNamesAutoCompletion(flags, args, toComplete)
}
