package users

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

func NewCommand() *cobra.Command {
	return fctl.NewMembershipCommand("users",
		fctl.WithAliases("s"),
		fctl.WithShortDescription("Stack users management within an organization"),
		fctl.WithRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}

			apiClient, err := fctl.NewMembershipClient(cmd, cfg)
			if err != nil {
				return err
			}

			version := fctl.MembershipServerInfo(cmd.Context(), apiClient)
			if !semver.IsValid(version) {
				return nil
			}

			if semver.Compare(version, "v0.26.1") >= 0 {
				return nil
			}

			return fmt.Errorf("unsupported membership server version: %s", version)
		}),
		fctl.WithChildCommands(
			NewLinkCommand(),
			NewListCommand(),
			NewUnlinkCommand(),
		),
	)
}
