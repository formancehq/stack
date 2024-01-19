package fctl

import (
	"context"
	"fmt"

	"github.com/formancehq/fctl/membershipclient"
	formance "github.com/formancehq/formance-sdk-go/v2"
	"github.com/spf13/cobra"
)

func getVersion(cmd *cobra.Command) string {
	for cmd != nil {
		if cmd.Version != "" {
			return cmd.Version
		}
		cmd = cmd.Parent()
	}
	return "cmd.Version"
}

func NewMembershipClient(cmd *cobra.Command, cfg *Config) (*membershipclient.APIClient, error) {
	profile := GetCurrentProfile(cmd, cfg)
	httpClient := GetHttpClient(cmd, map[string][]string{})
	configuration := membershipclient.NewConfiguration()
	token, err := profile.GetToken(cmd.Context(), httpClient)
	if err != nil {
		return nil, err
	}
	configuration.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	configuration.HTTPClient = httpClient
	configuration.UserAgent = "fctl/" + getVersion(cmd)
	configuration.Servers[0].URL = profile.GetMembershipURI()
	return membershipclient.NewAPIClient(configuration), nil
}

func MembershipServerInfo(ctx context.Context, client *membershipclient.APIClient) string {
	serverInfo, response, err := client.DefaultApi.GetServerInfo(ctx).Execute()
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	if response.StatusCode != 200 {
		return fmt.Sprintf("Error: %s", response.Status)
	}
	return serverInfo.Version
}

func NewStackClient(cmd *cobra.Command, cfg *Config, stack *membershipclient.Stack) (*formance.Formance, error) {
	profile := GetCurrentProfile(cmd, cfg)
	httpClient := GetHttpClient(cmd, map[string][]string{})

	token, err := profile.GetStackToken(cmd.Context(), httpClient, stack)
	if err != nil {
		return nil, err
	}

	return formance.New(
		formance.WithServerURL(stack.Uri),
		formance.WithClient(
			GetHttpClient(cmd, map[string][]string{
				"Authorization": {fmt.Sprintf("Bearer %s", token)},
				"User-Agent":    {"fctl/" + getVersion(cmd)},
			}),
		),
	), nil
}
