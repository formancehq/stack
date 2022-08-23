package svix

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/numary/webhooks/constants"
	"github.com/spf13/viper"
	svix "github.com/svix/svix-webhooks/go"
)

type App struct {
	Client *svix.Svix
	AppId  string
}

func New(httpClient *http.Client) (App, error) {
	token := viper.GetString(constants.SvixTokenFlag)
	appName := viper.GetString(constants.SvixAppNameFlag)
	appId := viper.GetString(constants.SvixAppIdFlag)
	serverUrl := viper.GetString(constants.SvixServerUrlFlag)

	u, err := url.Parse(serverUrl)
	if err != nil {
		return App{}, fmt.Errorf("url.Parse: %w", err)
	}

	client := svix.New(token, &svix.SvixOptions{
		ServerUrl:  u,
		HTTPClient: httpClient,
	})

	app, err := client.Application.GetOrCreate(&svix.ApplicationIn{
		Name: appName,
		Uid:  *svix.NullableString(appId),
	})
	if err != nil {
		return App{}, fmt.Errorf("svix.Application.GetOrCreate: %w", err)
	}

	return App{
		Client: client,
		AppId:  app.Id,
	}, nil
}
