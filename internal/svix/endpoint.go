package svix

import (
	"fmt"
	"time"

	"github.com/numary/go-libs/sharedlogging"
	svixgo "github.com/svix/svix-webhooks/go"
)

func CreateEndpoint(svixClient *svixgo.Svix, svixAppId, url string) error {
	out, err := svixClient.Endpoint.Create(svixAppId, &svixgo.EndpointIn{
		Url:     url,
		Version: 1,
	})
	if err != nil {
		return fmt.Errorf("svix.Endpoint.Create: appId: %s: url: %s: %w",
			svixAppId, url, err)
	}

	sharedlogging.Infof("svix endpoint created: url: %s id: %s createdAt: %s",
		url, out.Id, out.CreatedAt.Format(time.RFC3339))
	return nil
}

func DeleteAllEndpoints(svixClient *svixgo.Svix, svixAppId string) error {
	endpointList, err := svixClient.Endpoint.List(svixAppId, &svixgo.EndpointListOptions{})
	if err != nil {
		return fmt.Errorf("svix.Endpoint.List: appId: %s: %w", svixAppId, err)
	}

	for _, endpoint := range endpointList.Data {
		if err := svixClient.Endpoint.Delete(svixAppId, endpoint.Id); err != nil {
			return fmt.Errorf("svix.Endpoint.Delete: app: %s: endpointId: %s: %w",
				svixAppId, endpoint.Id, err)
		}
		sharedlogging.Infof("svix endpoint deleted: url: %s id: %s createdAt: %s",
			endpoint.Url, endpoint.Id, endpoint.CreatedAt.Format(time.RFC3339))
	}

	return nil
}
