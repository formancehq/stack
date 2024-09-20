package fctl

import (
	"fmt"

	"github.com/formancehq/fctl/membershipclient"
	"github.com/formancehq/go-libs/collectionutils"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

func CheckMembershipVersion(version string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		store := GetMembershipStore(cmd.Context())
		serverInfo, err := MembershipServerInfo(cmd.Context(), store.Client())
		if err != nil {
			return err
		}
		if !semver.IsValid(serverInfo.Version) {
			return nil
		}

		if semver.Compare(serverInfo.Version, version) >= 0 {
			return nil
		}

		return fmt.Errorf("unsupported membership server version: %s", version)
	}

}

func CheckMembershipCapabilities(capability membershipclient.Capability) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		store := GetMembershipStore(cmd.Context())
		serverInfo, err := MembershipServerInfo(cmd.Context(), store.Client())
		if err != nil {
			return err
		}

		if collectionutils.Contains(serverInfo.Capabilities, capability) {
			return nil
		}

		return fmt.Errorf("unsupported membership server capability: %s", capability)
	}
}
