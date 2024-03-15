package stack

import (
	"context"
	"fmt"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type UpgradeStore struct {
	Stack *membershipclient.Stack
}

type UpgradeController struct {
	store   *UpgradeStore
	profile *fctl.Profile
}

var _ fctl.Controller[*UpgradeStore] = (*UpgradeController)(nil)

func NewDefaultUpgradeStore() *UpgradeStore {
	return &UpgradeStore{
		Stack: &membershipclient.Stack{},
	}
}
func NewUpgradeController() *UpgradeController {
	return &UpgradeController{
		store: NewDefaultUpgradeStore(),
	}
}

func NewUpgradeCommand() *cobra.Command {
	return fctl.NewMembershipCommand("upgrade <stack-id> <version>",
		fctl.WithShortDescription("Upgrade a stack to specified version"),
		fctl.WithBoolFlag(nowaitFlag, false, "Wait stack availability"),
		fctl.WithArgs(cobra.RangeArgs(1, 2)),
		fctl.WithPreRunE(func(cmd *cobra.Command, args []string) error {
			cfg, err := fctl.GetConfig(cmd)
			if err != nil {
				return err
			}

			client, err := fctl.NewMembershipClient(cmd, cfg)
			if err != nil {
				return err
			}

			version := fctl.MembershipServerInfo(cmd.Context(), client.APIClient)
			if !semver.IsValid(version) {
				return nil
			}

			if semver.Compare(version, "v0.27.1") >= 0 {
				return nil
			}

			return fmt.Errorf("unsupported membership server version: %s", version)
		}),
		fctl.WithController[*UpgradeStore](NewUpgradeController()),
	)
}
func (c *UpgradeController) GetStore() *UpgradeStore {
	return c.store
}

func (c *UpgradeController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(cmd, cfg)
	c.profile = profile

	organization, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}

	stack, res, err := client.DefaultApi.GetStack(cmd.Context(), organization, args[0]).Execute()
	if err != nil {
		return nil, errors.Wrap(err, "retrieving stack")
	}

	if res.StatusCode > 300 {
		return nil, err
	}

	req := membershipclient.StackVersion{
		Version: nil,
	}
	specifiedVersion := fctl.GetString(cmd, versionFlag)
	if specifiedVersion == "" {
		upgradeOpts, err := retrieveUpgradableVersion(cmd.Context(), organization, *stack.Data, client.APIClient)
		if err != nil {
			return nil, err
		}
		printer := pterm.DefaultInteractiveSelect.WithOptions(upgradeOpts)
		selectedOption, err := printer.Show("Please select a version")
		if err != nil {
			return nil, err
		}

		specifiedVersion = selectedOption
	}

	if specifiedVersion != *stack.Data.Version {
		if !fctl.CheckStackApprobation(cmd, stack.Data, "Disclaimer: You are about to migrate the stack '%s' from '%s' to '%s'. It might take some time to fully migrate", stack.Data.Name, *stack.Data.Version, specifiedVersion) {
			return nil, fctl.ErrMissingApproval
		}
	} else {
		pterm.Warning.WithWriter(cmd.OutOrStdout()).Printfln("Stack is already at version %s", specifiedVersion)
		return nil, nil
	}
	req.Version = pointer.For(specifiedVersion)

	res, err = client.DefaultApi.
		UpgradeStack(cmd.Context(), organization, args[0]).StackVersion(req).
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, "upgrading stack")
	}

	if res.StatusCode > 300 {
		return nil, err
	}

	if !fctl.GetBool(cmd, nowaitFlag) {
		spinner, err := pterm.DefaultSpinner.Start("Waiting services availability")
		if err != nil {
			return nil, err
		}

		stack, err := waitStackReady(cmd, client, stack.Data.OrganizationId, stack.Data.Id)
		if err != nil {
			return nil, err
		}
		c.store.Stack = stack

		if err := spinner.Stop(); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *UpgradeController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Stack upgrade progressing.")
	return nil
}

func retrieveUpgradableVersion(ctx context.Context, organization string, stack membershipclient.Stack, apiClient *membershipclient.APIClient) ([]string, error) {
	availableVersions, httpResponse, err := apiClient.DefaultApi.GetRegionVersions(ctx, organization, stack.RegionID).Execute()
	if httpResponse == nil {
		return nil, err
	}

	var upgradeOptions []string
	for _, version := range availableVersions.Data {
		if version.Name == *stack.Version {
			continue
		}
		if !semver.IsValid(version.Name) || semver.Compare(version.Name, *stack.Version) >= 1 {
			upgradeOptions = append(upgradeOptions, version.Name)
		}
	}
	return upgradeOptions, nil
}
