package stack

import (
	"net/http"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
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
		fctl.WithArgs(cobra.RangeArgs(1, 2)),
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

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return nil, err
	}

	if len(args) < 1 {
		return nil, errors.New("stack-id is required")
	}

	stack, res, err := apiClient.DefaultApi.GetStack(cmd.Context(), organization, args[0]).Execute()
	if err != nil {
		return nil, errors.Wrap(err, "retrieving stack")
	}
	if res.StatusCode > 300 {
		return nil, errors.New("stack not found")
	}

	availableVersions, httpResponse, err := apiClient.DefaultApi.GetRegionVersions(cmd.Context(), organization, stack.Data.RegionID).Execute()
	if httpResponse == nil {
		return nil, err
	}

	req := membershipclient.StackVersion{
		Version: nil,
	}
	specifiedVersion := fctl.GetString(cmd, versionFlag)
	switch {
	case httpResponse.StatusCode == http.StatusOK && specifiedVersion == "":
		var options []string
		for _, version := range availableVersions.Data {
			options = append(options, version.Name)
		}

		printer := pterm.DefaultInteractiveSelect.WithOptions(options)
		selectedOption, err := printer.Show("Please select a version")
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(options); i++ {
			if selectedOption == options[i] {
				specifiedVersion = availableVersions.Data[i].Name
				break
			}
		}
	case httpResponse.StatusCode != http.StatusOK && specifiedVersion == "":
		// nothing to do, we cannot set a specific version for this membership version
	case httpResponse.StatusCode == http.StatusOK && specifiedVersion != "":
		// nothing to do, let membership handle the case
	case httpResponse.StatusCode != http.StatusOK && specifiedVersion != "":
		return nil, errors.New("--version flag can not be used with the actual membership api")
	}

	if specifiedVersion != "" {
		if specifiedVersion != *stack.Data.Version {
			if !fctl.CheckStackApprobation(cmd, stack.Data, "Disclaimer: You are about to migrate the stack '%s' from '%s' to '%s'. It might take some time to fully migrate", stack.Data.Name, *stack.Data.Version, specifiedVersion) {
				return nil, fctl.ErrMissingApproval
			}
		}
		req.Version = pointer.For(specifiedVersion)
	}

	res, err = apiClient.DefaultApi.
		UpgradeStack(cmd.Context(), organization, args[0]).StackVersion(req).
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, "upgrading stack")
	}

	if res.StatusCode > 300 {
		return nil, err
	}

	return c, nil
}

func (c *UpgradeController) Render(cmd *cobra.Command, args []string) error {
	pterm.Success.WithWriter(cmd.OutOrStdout()).Printfln("Stack upgrade progressing.")
	return nil
}
