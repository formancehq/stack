package versions

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type Version int

const (
	V0 Version = iota
	V1
)

type VersionController interface {
	SetVersion(Version)
}

func GetPaymentsVersion(cmd *cobra.Command, args []string, controller VersionController) error {
	soc, err := fctl.GetStackOrganizationConfig(cmd)
	if err != nil {
		return err
	}

	client, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return errors.Wrap(err, "creating stack client")
	}

	response, err := client.Payments.PaymentsgetServerInfo(cmd.Context())
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	version := "v" + response.ServerInfo.Version

	if !semver.IsValid(version) {
		// last version for commits
		controller.SetVersion(V1)
		return nil
	}

	res := semver.Compare(version, "v1.0.0-rc.1")
	switch res {
	case 0, 1:
		controller.SetVersion(V1)
	default:
		controller.SetVersion(V0)
	}

	return nil
}
