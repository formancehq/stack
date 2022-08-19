package svix

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/numary/webhooks/constants"
	"github.com/spf13/viper"
	svix "github.com/svix/svix-webhooks/go"
)

func New(httpClient *http.Client) (svixClient *svix.Svix, svixAppId string, err error) {
	token := viper.GetString(constants.SvixTokenFlag)
	appName := viper.GetString(constants.SvixAppNameFlag)
	appId := viper.GetString(constants.SvixAppIdFlag)
	serverUrl := viper.GetString(constants.SvixServerUrlFlag)

	u, err := url.Parse(serverUrl)
	if err != nil {
		return nil, "", fmt.Errorf("url.Parse: %w", err)
	}

	svixClient = svix.New(token, &svix.SvixOptions{
		ServerUrl:  u,
		HTTPClient: httpClient,
	})

	app, err := svixClient.Application.GetOrCreate(&svix.ApplicationIn{
		Name: appName,
		Uid:  *svix.NullableString(appId),
	})
	if err != nil {
		return nil, "", fmt.Errorf("svix.Application.GetOrCreate: %w", err)
	}

	return svixClient, app.Id, nil
}
