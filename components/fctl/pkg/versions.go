package fctl

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

func CheckMembershipVersion(version string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		store := GetMembershipStore(cmd.Context())
		serverVersion := MembershipServerInfo(cmd.Context(), store.Client())
		if !semver.IsValid(serverVersion) {
			return nil
		}

		if semver.Compare(serverVersion, version) >= 0 {
			return nil
		}

		return fmt.Errorf("unsupported membership server version: %s", version)
	}

}
