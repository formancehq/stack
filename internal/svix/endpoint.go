package svix

import (
	"fmt"

	"github.com/numary/webhooks/pkg/model"
	svixgo "github.com/svix/svix-webhooks/go"
)

func CreateEndpoint(endpointId string, cfg model.Config, svixClient *svixgo.Svix, svixAppId string) error {
	list, err := ListEventTypes(svixClient)
	if err != nil {
		return fmt.Errorf("svix.ListEventTypes: %w", err)
	}

	for _, newEventType := range cfg.EventTypes {
		alreadyCreated := false
		for _, eventType := range list.Data {
			if eventType.Name == newEventType {
				alreadyCreated = true
			}
		}
		if !alreadyCreated {
			eventTypeIn := svixgo.EventTypeIn{}
			var archived = false
			eventTypeIn.Archived = &archived
			eventTypeIn.Name = newEventType
			if _, err := svixClient.EventType.Create(&eventTypeIn); err != nil {
				return fmt.Errorf("svix.EventType.Create: %w", err)
			}
		}
	}

	endpointIn := &svixgo.EndpointIn{
		FilterTypes: cfg.EventTypes,
		Uid:         *svixgo.NullableString(endpointId),
		Url:         cfg.Endpoint,
		Version:     1,
	}
	opts := &svixgo.PostOptions{IdempotencyKey: &endpointId}
	if _, err := svixClient.Endpoint.CreateWithOptions(svixAppId, endpointIn, opts); err != nil {
		return fmt.Errorf("svix.Endpoint.CreateWithOptions: %w", err)
	}

	return nil
}

func ListEventTypes(svixClient *svixgo.Svix) (*svixgo.ListResponseEventTypeOut, error) {
	return svixClient.EventType.List(&svixgo.EventTypeListOptions{})
}

func DeleteEndpoint(endpointId string, svixClient *svixgo.Svix, svixAppId string) error {
	return svixClient.Endpoint.Delete(svixAppId, endpointId)
}

func ToggleEndpoint(endpointId string, updatedCfg model.ConfigInserted, svixClient *svixgo.Svix, svixAppId string) error {
	disabled := !updatedCfg.Active
	_, err := svixClient.Endpoint.Update(svixAppId, endpointId, &svixgo.EndpointUpdate{
		Disabled:    &disabled,
		FilterTypes: updatedCfg.EventTypes,
		Uid:         *svixgo.NullableString(updatedCfg.ID),
		Url:         updatedCfg.Endpoint,
		Version:     1,
	})
	if err != nil {
		return fmt.Errorf("svix.Endpoint.Update: %w", err)
	}

	return nil
}
