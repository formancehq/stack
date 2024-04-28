package fctl

import (
	"context"
	"fmt"

	"github.com/formancehq/fctl/membershipclient"
	formance "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/stack/libs/go-libs/logging"
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

type MembershipClient struct {
	profile *Profile
	*membershipclient.APIClient
}

func (c *MembershipClient) GetProfile() *Profile {
	return c.profile
}

func (c *MembershipClient) RefreshIfNeeded(cmd *cobra.Command) error {
	logging.Debug("Refreshing membership client")
	token, err := c.profile.GetToken(cmd.Context(), c.GetConfig().HTTPClient)
	if err != nil {
		return err
	}
	config := c.GetConfig()
	config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	c.APIClient = membershipclient.NewAPIClient(config)
	return err
}

func NewMembershipClient(cmd *cobra.Command, cfg *Config) (*MembershipClient, error) {
	profile := GetCurrentProfile(cmd, cfg)
	httpClient := GetHttpClient(cmd, map[string][]string{})
	configuration := membershipclient.NewConfiguration()
	configuration.HTTPClient = httpClient
	configuration.UserAgent = "fctl/" + getVersion(cmd)
	configuration.Servers[0].URL = profile.GetMembershipURI()
	client := &MembershipClient{
		APIClient: membershipclient.NewAPIClient(configuration),
		profile:   profile,
	}
	err := client.RefreshIfNeeded(cmd)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func MembershipServerInfo(ctx context.Context, client *membershipclient.DefaultApiService) string {
	serverInfo, response, err := client.GetServerInfo(ctx).Execute()
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
