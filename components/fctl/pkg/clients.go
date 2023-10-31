package fctl

import (
	"fmt"

	"github.com/formancehq/fctl/membershipclient"
	formance "github.com/formancehq/formance-sdk-go"
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
