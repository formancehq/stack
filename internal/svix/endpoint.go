package svix

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/pkg/model"
	svixgo "github.com/svix/svix-webhooks/go"
)

func CreateEndpoint(ctx context.Context, endpointId string, cfg model.Config, svixClient *svixgo.Svix, svixAppId string) error {
	if err := makeSureEventTypesFromCfgAreCreated(ctx, svixClient, cfg); err != nil {
		return fmt.Errorf("makeSureEventTypesFromCfgAreCreated: %w", err)
	}

	endpointIn := &svixgo.EndpointIn{
		FilterTypes: cfg.EventTypes,
		Secret:      *svixgo.NullableString("whsec_" + cfg.Secret),
		Uid:         *svixgo.NullableString(endpointId),
		Url:         cfg.Endpoint,
		Version:     1,
	}
	opts := &svixgo.PostOptions{IdempotencyKey: &endpointId}
	if out, err := svixClient.Endpoint.CreateWithOptions(svixAppId, endpointIn, opts); err != nil {
		return fmt.Errorf("svix.Svix.Endpoint.CreateWithOptions: %w", err)
	} else {
		dumpOut := spew.Sdump(out)
		sharedlogging.GetLogger(ctx).Debug("svix.Svix.Endpoint.CreateWithOptions: ", dumpOut)
	}

	return nil
}

func makeSureEventTypesFromCfgAreCreated(ctx context.Context, svixClient *svixgo.Svix, cfg model.Config) error {
	includeArchived, withContent := true, true
	eventTypeListOptions := svixgo.EventTypeListOptions{
		IncludeArchived: &includeArchived,
		WithContent:     &withContent,
	}
	list, err := svixClient.EventType.List(&eventTypeListOptions)
	if err != nil {
		return fmt.Errorf("svix.Svix.EventType.List: %w", err)
	}

	for _, newEventType := range cfg.EventTypes {
		alreadyCreated := false
		for _, eventType := range list.Data {
			if eventType.Name == newEventType {
				alreadyCreated = true
			}
		}
		if !alreadyCreated {
			var archived = false
			eventTypeIn := svixgo.EventTypeIn{
				Archived: &archived,
				Name:     newEventType,
			}
			if out, err := svixClient.EventType.Create(&eventTypeIn); err != nil {
				return fmt.Errorf("svix.Svix.EventType.Create: %w", err)
			} else {
				dumpOut := spew.Sdump(out)
				sharedlogging.GetLogger(ctx).Debug("svix.Svix.EventType.Create: ", dumpOut)
			}
		}
	}

	return nil
}

func DeleteOneEndpoint(endpointId string, svixClient *svixgo.Svix, svixAppId string) error {
	return svixClient.Endpoint.Delete(svixAppId, endpointId)
}

func UpdateOneEndpoint(ctx context.Context, endpointId string, updatedCfg model.ConfigInserted, svixClient *svixgo.Svix, svixAppId string) error {
	disabled := !updatedCfg.Active
	endpointUpdate := svixgo.EndpointUpdate{
		Disabled:    &disabled,
		FilterTypes: updatedCfg.EventTypes,
		Uid:         *svixgo.NullableString(updatedCfg.ID),
		Url:         updatedCfg.Endpoint,
		Version:     1,
	}
	if out, err := svixClient.Endpoint.Update(svixAppId, endpointId, &endpointUpdate); err != nil {
		return fmt.Errorf("svix.Svix.Endpoint.Update: %w", err)
	} else {
		sharedlogging.GetLogger(ctx).Debug("svix.Svix.Endpoint.Update: ", spew.Sdump(out))
	}

	return nil
}

func RotateOneEndpointSecret(ctx context.Context, endpointId, secret string, svixClient *svixgo.Svix, svixAppId string) error {
	endpointSecretRotateIn := &svixgo.EndpointSecretRotateIn{
		Key: *svixgo.NullableString("whsec_" + secret),
	}
	if err := svixClient.Endpoint.RotateSecret(svixAppId, endpointId, endpointSecretRotateIn); err != nil {
		return fmt.Errorf("svix.Svix.Endpoint.RotateSecret: %w", err)
	} else {
		sharedlogging.GetLogger(ctx).Debug("svix.Svix.Endpoint.RotateSecret: OK")
	}

	return nil
}
