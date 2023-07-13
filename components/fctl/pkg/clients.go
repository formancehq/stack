package fctl

import (
	"context"
	"flag"
	"fmt"

	"github.com/formancehq/fctl/membershipclient"
	"github.com/formancehq/formance-sdk-go"
)

func getVersion() string {
	return Version
}

func NewMembershipClient(flags *flag.FlagSet, ctx context.Context, cfg *Config) (*membershipclient.APIClient, error) {
	profile := GetCurrentProfile(flags, cfg)

	httpClient := GetHttpClient(flags, map[string][]string{})

	configuration := membershipclient.NewConfiguration()

	token, err := profile.GetToken(ctx, httpClient)
	if err != nil {
		return nil, err
	}

	configuration.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	configuration.HTTPClient = httpClient
	configuration.UserAgent = "fctl/" + getVersion()
	configuration.Servers[0].URL = profile.GetMembershipURI()

	return membershipclient.NewAPIClient(configuration), nil
}

func NewStackClient(flags *flag.FlagSet, ctx context.Context, cfg *Config, stack *membershipclient.Stack) (*formance.Formance, error) {
	profile := GetCurrentProfile(flags, cfg)
	httpClient := GetHttpClient(flags, map[string][]string{})

	token, err := profile.GetStackToken(ctx, httpClient, stack)
	if err != nil {
		return nil, err
	}

	return formance.New(
		formance.WithServerURL(stack.Uri),
		formance.WithClient(
			GetHttpClient(flags, map[string][]string{
				"Authorization": {fmt.Sprintf("Bearer %s", token)},
				"User-Agent":    {"fctl/" + getVersion()},
			}),
		),
	), nil
}
