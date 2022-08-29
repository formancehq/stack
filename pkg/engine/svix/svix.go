package svix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/constants"
	"github.com/numary/webhooks/pkg/engine"
	"github.com/numary/webhooks/pkg/model"
	"github.com/spf13/viper"
	svixgo "github.com/svix/svix-webhooks/go"
)

type Engine struct {
	Client *svixgo.Svix
	AppId  string
}

var _ engine.Engine = &Engine{}

func NewEngine(httpClient *http.Client) (Engine, error) {
	token := viper.GetString(constants.SvixTokenFlag)
	appName := viper.GetString(constants.SvixAppNameFlag)
	appId := viper.GetString(constants.SvixAppIdFlag)
	serverUrl := viper.GetString(constants.SvixServerUrlFlag)

	u, err := url.Parse(serverUrl)
	if err != nil {
		return Engine{}, fmt.Errorf("url.Parse: %w", err)
	}

	client := svixgo.New(token, &svixgo.SvixOptions{
		ServerUrl:  u,
		HTTPClient: httpClient,
	})

	app, err := client.Application.GetOrCreate(&svixgo.ApplicationIn{
		Name: appName,
		Uid:  *svixgo.NullableString(appId),
	})
	if err != nil {
		return Engine{}, fmt.Errorf("svix.Application.GetOrCreate: %w", err)
	}

	return Engine{
		Client: client,
		AppId:  app.Id,
	}, nil
}

func (a Engine) InsertOneConfig(ctx context.Context, id string, cfg model.Config) error {
	if err := makeSureEventTypesFromCfgAreCreated(ctx, cfg, a); err != nil {
		return fmt.Errorf("makeSureEventTypesFromCfgAreCreated: %w", err)
	}

	endpointIn := &svixgo.EndpointIn{
		FilterTypes: cfg.EventTypes,
		Secret:      *svixgo.NullableString("whsec_" + cfg.Secret),
		Uid:         *svixgo.NullableString(id),
		Url:         cfg.Endpoint,
		Version:     1,
	}
	opts := &svixgo.PostOptions{IdempotencyKey: &id}
	if out, err := a.Client.Endpoint.CreateWithOptions(a.AppId, endpointIn, opts); err != nil {
		return fmt.Errorf("svix.Svix.Endpoint.CreateWithOptions: %w", err)
	} else {
		dumpOut := spew.Sdump(out)
		sharedlogging.GetLogger(ctx).Debug("svix.Svix.Endpoint.CreateWithOptions: ", dumpOut)
	}

	return nil
}

func makeSureEventTypesFromCfgAreCreated(ctx context.Context, cfg model.Config, svixApp Engine) error {
	includeArchived, withContent := true, true
	eventTypeListOptions := svixgo.EventTypeListOptions{
		IncludeArchived: &includeArchived,
		WithContent:     &withContent,
	}
	list, err := svixApp.Client.EventType.List(&eventTypeListOptions) //nolint:contextcheck
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
			archived := false
			eventTypeIn := svixgo.EventTypeIn{
				Archived: &archived,
				Name:     newEventType,
			}
			if out, err := svixApp.Client.EventType.Create(&eventTypeIn); err != nil { //nolint:contextcheck
				return fmt.Errorf("svix.Svix.EventType.Create: %w", err)
			} else {
				dumpOut := spew.Sdump(out)
				sharedlogging.GetLogger(ctx).Debug("svix.Svix.EventType.Create: ", dumpOut)
			}
		}
	}

	return nil
}

func (a Engine) DeleteOneConfig(ctx context.Context, id string) error {
	if err := a.Client.Endpoint.Delete(a.AppId, id); err != nil { //nolint:contextcheck
		return fmt.Errorf("svix.Svix.Endpoint.Delete: %w", err)
	}

	sharedlogging.GetLogger(ctx).Debug("svix.Svix.Endpoint.Delete: %s", id)
	return nil
}

func (a Engine) UpdateOneConfig(ctx context.Context, id string, cfg *model.ConfigInserted) error {
	disabled := !cfg.Active
	endpointUpdate := svixgo.EndpointUpdate{
		Disabled:    &disabled,
		FilterTypes: cfg.EventTypes,
		Uid:         *svixgo.NullableString(cfg.ID),
		Url:         cfg.Endpoint,
		Version:     1,
	}
	if out, err := a.Client.Endpoint.Update(a.AppId, id, &endpointUpdate); err != nil { //nolint:contextcheck
		return fmt.Errorf("svix.Svix.Endpoint.Update: %w", err)
	} else {
		sharedlogging.GetLogger(ctx).Debug("svix.Svix.Endpoint.Update: ", spew.Sdump(out))
	}

	return nil
}

func (a Engine) RotateOneConfigSecret(ctx context.Context, id, secret string) error {
	endpointSecretRotateIn := &svixgo.EndpointSecretRotateIn{
		Key: *svixgo.NullableString("whsec_" + secret),
	}
	if err := a.Client.Endpoint.RotateSecret(a.AppId, id, endpointSecretRotateIn); err != nil { //nolint:contextcheck
		return fmt.Errorf("svix.Svix.Endpoint.RotateSecret: %w", err)
	} else {
		sharedlogging.GetLogger(ctx).Debug("svix.Svix.Endpoint.RotateSecret: OK")
	}

	return nil
}

func (a Engine) ProcessKafkaMessage(ctx context.Context, eventType string, msgValue []byte) error {
	id := uuid.NewString()
	var p map[string]any
	if err := json.Unmarshal(msgValue, &p); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	messageIn := &svixgo.MessageIn{
		EventId:   *svixgo.NullableString(id),
		EventType: eventType,
		Payload:   p,
	}

	options := &svixgo.PostOptions{IdempotencyKey: &id}
	dumpIn := spew.Sdump(
		"svix appId: ", a.AppId,
		"svix.MessageIn: ", messageIn,
		"svix.PostOptions: ", options)

	if out, err := a.Client.Message.CreateWithOptions( //nolint:contextcheck
		a.AppId, messageIn, options); err != nil {
		return fmt.Errorf("svix.Svix.Message.CreateWithOptions: %w: dumpIn: %s",
			err, dumpIn)
	} else {
		fmt.Printf("\nNEW WEBHOOK SENT\n%+v\n", out)
	}

	return nil
}
