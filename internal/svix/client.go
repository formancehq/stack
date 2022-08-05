package svix

import (
	"fmt"
	"net/url"

	"github.com/numary/webhooks-cloud/cmd/constants"
	"github.com/spf13/viper"
	svix "github.com/svix/svix-webhooks/go"
)

func New() (*svix.Svix, string, error) {
	token := viper.GetString(constants.SvixTokenFlag)
	organizationName := viper.GetString(constants.SvixOrganizationNameFlag)
	serverUrl := viper.GetString(constants.SvixServerUrlFlag)

	urlForServer, _ := url.Parse(serverUrl)
	opt := svix.SvixOptions{
		ServerUrl: urlForServer,
	}
	svixClient := svix.New(token, &opt)

	ApplicationListOptions := svix.ApplicationListOptions{}
	list, _ := svixClient.Application.List(&ApplicationListOptions)
	for _, s := range list.Data {
		if s.Id == organizationName {
			return svixClient, s.Id, nil
		}
	}

	optApp := svix.ApplicationIn{
		Name: organizationName,
	}
	app, err := svixClient.Application.Create(&optApp)
	if err != nil {
		return nil, "", fmt.Errorf("error creating svix application: %s", err)
	}

	return svixClient, app.Id, nil
}
