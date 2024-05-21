package stack

import (
	"fmt"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/pointer"

	"github.com/formancehq/fctl/cmd/stack/internal"
	"github.com/formancehq/fctl/cmd/stack/store"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	unprotectFlag = "unprotect"
	regionFlag    = "region"
	nowaitFlag    = "no-wait"
	versionFlag   = "version"
)

type StackCreateStore struct {
	Stack    *membershipclient.Stack
	Versions *shared.GetVersionsResponse
}

type StackCreateController struct {
	store   *StackCreateStore
	profile *fctl.Profile
}

var _ fctl.Controller[*StackCreateStore] = (*StackCreateController)(nil)

func NewDefaultStackCreateStore() *StackCreateStore {
	return &StackCreateStore{
		Stack:    &membershipclient.Stack{},
		Versions: &shared.GetVersionsResponse{},
	}
}
func NewStackCreateController() *StackCreateController {
	return &StackCreateController{
		store: NewDefaultStackCreateStore(),
	}
}

func NewCreateCommand() *cobra.Command {
	return fctl.NewMembershipCommand("create [name]",
		fctl.WithShortDescription("Create a new stack"),
		fctl.WithAliases("c", "cr"),
		fctl.WithArgs(cobra.RangeArgs(0, 1)),
		fctl.WithBoolFlag(unprotectFlag, false, "Unprotect stacks (no confirmation on write commands)"),
		fctl.WithStringFlag(regionFlag, "", "Region on which deploy the stack"),
		fctl.WithStringFlag(versionFlag, "", "Version of the stack"),
		fctl.WithBoolFlag(nowaitFlag, false, "Not wait stack availability"),
		fctl.WithController[*StackCreateStore](NewStackCreateController()),
	)
}
func (c *StackCreateController) GetStore() *StackCreateStore {
	return c.store
}

func (c *StackCreateController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	var err error
	store := store.GetStore(cmd.Context())

	protected := !fctl.GetBool(cmd, unprotectFlag)
	metadata := map[string]string{
		fctl.ProtectedStackMetadata: fctl.BoolPointerToString(&protected),
	}

	name := ""
	if len(args) > 0 {
		name = args[0]
	} else {
		name, err = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter a name")
		if err != nil {
			return nil, err
		}
	}

	region := fctl.GetString(cmd, regionFlag)
	if region == "" {
		regions, _, err := store.Client().ListRegions(cmd.Context(), store.OrganizationId()).Execute()
		if err != nil {
			return nil, errors.Wrap(err, "listing regions")
		}

		var options []string
		for _, region := range regions.Data {
			privacy := "Private"
			if region.Public {
				privacy = "Public "
			}
			name := "<noname>"
			if region.Name != "" {
				name = region.Name
			}
			options = append(options, fmt.Sprintf("%s | %s | %s", region.Id, privacy, name))
		}

		printer := pterm.DefaultInteractiveSelect.WithOptions(options)
		selectedOption, err := printer.Show("Please select a region")
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(options); i++ {
			if selectedOption == options[i] {
				region = regions.Data[i].Id
				break
			}
		}
	}

	req := membershipclient.CreateStackRequest{
		Name:     name,
		Metadata: pointer.For(metadata),
		RegionID: region,
	}

	availableVersions, httpResponse, err := store.Client().GetRegionVersions(cmd.Context(), store.OrganizationId(), region).Execute()
	if err != nil {
		return nil, errors.Wrap(err, "retrieving available versions")
	}

	if httpResponse.StatusCode > 300 {
		return nil, err
	}

	specifiedVersion := fctl.GetString(cmd, versionFlag)
	if specifiedVersion == "" {
		var options []string
		for _, version := range availableVersions.Data {
			options = append(options, version.Name)
		}

		selectedOption := ""
		if len(options) > 0 {
			printer := pterm.DefaultInteractiveSelect.WithOptions(options)
			selectedOption, err = printer.Show("Please select a version")
			if err != nil {
				return nil, err
			}
		}

		specifiedVersion = selectedOption
	}
	req.Version = pointer.For(specifiedVersion)

	stackResponse, _, err := store.Client().
		CreateStack(cmd.Context(), store.OrganizationId()).
		CreateStackRequest(req).
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, "creating stack")
	}

	if !fctl.GetBool(cmd, nowaitFlag) {
		spinner, err := pterm.DefaultSpinner.Start("Waiting services availability")
		if err != nil {
			return nil, err
		}

		stack, err := waitStackReady(cmd, store.MembershipClient, stackResponse.Data.OrganizationId, stackResponse.Data.Id)
		if err != nil {
			return nil, err
		}
		c.store.Stack = stack

		if err := spinner.Stop(); err != nil {
			return nil, err
		}
	} else {
		c.store.Stack = stackResponse.Data
	}

	dashboard := "https://console.formance.cloud"
	serverInfo, err := fctl.MembershipServerInfo(cmd.Context(), store.Client())
	if err != nil {
		return nil, err
	}
	if v := serverInfo.ConsoleURL; v != nil {
		dashboard = *v
	}

	fctl.BasicTextCyan.WithWriter(cmd.OutOrStdout()).Println("Your dashboard will be reachable on: " + dashboard)

	stackClient, err := fctl.NewStackClient(cmd, store.Config, stackResponse.Data)
	if err != nil {
		return nil, err
	}

	versions, err := stackClient.GetVersions(cmd.Context())
	if err != nil {
		return nil, err
	}
	if versions.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when reading versions", versions.StatusCode)
	}

	c.store.Versions = versions.GetVersionsResponse
	c.profile = store.Config.GetProfile(fctl.GetCurrentProfileName(cmd, store.Config))

	return c, nil
}

func (c *StackCreateController) Render(cmd *cobra.Command, args []string) error {
	return internal.PrintStackInformation(cmd.OutOrStdout(), c.profile, c.store.Stack, c.store.Versions)
}
