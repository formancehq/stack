package stack

import (
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/fctl/cmd/stack/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	unprotectFlag = "unprotect"
	regionFlag    = "region"
	nowaitFlag    = "no-wait"
)

type StackCreate struct {
	Stack    *membershipclient.Stack
	Versions *shared.GetVersionsResponse
}

func NewCreateCommand() *cobra.Command {
	return fctl.NewMembershipCommand("create [name]",
		fctl.WithShortDescription("Create a new stack"),
		fctl.WithAliases("c", "cr"),
		fctl.WithArgs(cobra.RangeArgs(0, 1)),
		fctl.WithBoolFlag(unprotectFlag, false, "Unprotect stacks (no confirmation on write commands)"),
		fctl.WithStringFlag(regionFlag, "", "Region on which deploy the stack"),
		fctl.WithBoolFlag(nowaitFlag, false, "Not wait stack availability"),
		fctl.WithRunE(createStackCommand),
		fctl.WrapOutputPostRunE(viewStackCreate),
	)
}

func waitStackReady(cmd *cobra.Command, profile *fctl.Profile, stack *membershipclient.Stack) error {
	baseUrlStr := profile.ServicesBaseUrl(stack).String()
	authServerUrl := fmt.Sprintf("%s/api/auth", baseUrlStr)
	for {
		req, err := http.NewRequestWithContext(cmd.Context(), http.MethodGet,
			fmt.Sprintf(authServerUrl+"/.well-known/openid-configuration"), nil)
		if err != nil {
			return err
		}
		rsp, err := fctl.GetHttpClient(cmd, map[string][]string{}).Do(req)
		if err == nil && rsp.StatusCode == http.StatusOK {
			break
		}
		select {
		case <-cmd.Context().Done():
			return cmd.Context().Err()
		case <-time.After(time.Second):
		}
	}
	return nil
}

func createStackCommand(cmd *cobra.Command, args []string) error {

	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return err
	}

	organization, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return err
	}

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return err
	}

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
			return err
		}
	}

	region := fctl.GetString(cmd, regionFlag)
	if region == "" {
		regions, _, err := apiClient.DefaultApi.ListRegions(cmd.Context(), organization).Execute()
		if err != nil {
			return errors.Wrap(err, "listing regions")
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
			return err
		}
		for i := 0; i < len(options); i++ {
			if selectedOption == options[i] {
				region = regions.Data[i].Id
				break
			}
		}
	}

	stackResponse, _, err := apiClient.DefaultApi.CreateStack(cmd.Context(), organization).CreateStackRequest(membershipclient.CreateStackRequest{
		Name:     name,
		Metadata: metadata,
		RegionID: region,
	}).Execute()
	if err != nil {
		return errors.Wrap(err, "creating stack")
	}

	profile := fctl.GetCurrentProfile(cmd, cfg)

	if !fctl.GetBool(cmd, nowaitFlag) {
		spinner, err := pterm.DefaultSpinner.Start("Waiting services availability")
		if err != nil {
			return err
		}

		if err := waitStackReady(cmd, profile, stackResponse.Data); err != nil {
			return err
		}

		if err := spinner.Stop(); err != nil {
			return err
		}
	}

	fctl.BasicTextCyan.WithWriter(cmd.OutOrStdout()).Printfln("Your dashboard will be reachable on: %s",
		profile.ServicesBaseUrl(stackResponse.Data).String())

	stackClient, err := fctl.NewStackClient(cmd, cfg, stackResponse.Data)
	if err != nil {
		return err
	}

	versions, err := stackClient.GetVersions(cmd.Context())
	if err != nil {
		return err
	}
	if versions.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d when reading versions", versions.StatusCode)
	}

	fctl.SetSharedData(&StackCreate{
		Stack:    stackResponse.Data,
		Versions: versions.GetVersionsResponse,
	}, profile, nil)

	return nil
}

func viewStackCreate(cmd *cobra.Command, args []string) error {

	data := fctl.GetSharedData().(*StackCreate)

	return internal.PrintStackInformation(cmd.OutOrStdout(), fctl.GetSharedProfile(), data.Stack, data.Versions)
}
