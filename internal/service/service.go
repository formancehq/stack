package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/internal/model"
	"github.com/numary/webhooks/internal/storage"
	"github.com/numary/webhooks/internal/svix"
	svixgo "github.com/svix/svix-webhooks/go"
)

var ErrConfigNotFound = errors.New("config not found")

func InsertOneConfig(cfg model.Config, ctx context.Context, store storage.Store, svixClient *svixgo.Svix, svixAppId string) (string, error) {
	var id string
	var err error
	if id, err = store.InsertOneConfig(ctx, cfg); err != nil {
		return "", fmt.Errorf("store.InsertOneConfig: %w", err)
	}

	if err := svix.CreateEndpoint(ctx, id, cfg, svixClient, svixAppId); err != nil {
		if _, err = store.DeleteOneConfig(ctx, id); err != nil {
			return "", fmt.Errorf("store.DeleteOneConfig: %w", err)
		}
		return "", fmt.Errorf("svix.CreateEndpoint: %w", err)
	}

	sharedlogging.GetLogger(ctx).Debug("service.InsertOneConfig: id: ", id)
	return id, nil
}

func DeleteOneConfig(ctx context.Context, id string, store storage.Store, svixClient *svixgo.Svix, svixAppId string) error {
	if deletedCount, err := store.DeleteOneConfig(ctx, id); err != nil {
		return fmt.Errorf("store.DeleteOneConfig: %w", err)
	} else if deletedCount == 0 {
		return ErrConfigNotFound
	}

	if err := svix.DeleteOneEndpoint(id, svixClient, svixAppId); err != nil {
		return fmt.Errorf("svix.DeleteOneEndpoint: %w", err)
	}

	sharedlogging.GetLogger(ctx).Debug("service.DeleteOneConfig: id: ", id)
	return nil
}

func ActivateOneConfig(active bool, ctx context.Context, id string, store storage.Store, svixClient *svixgo.Svix, svixAppId string) error {
	updatedCfg, modifiedCount, err := store.UpdateOneConfigActive(ctx, id, active)
	if err != nil {
		return fmt.Errorf("store.UpdateOneConfigActive: %w", err)
	} else if modifiedCount == 0 {
		return ErrConfigNotFound
	}

	if err := svix.UpdateOneEndpoint(ctx, id, updatedCfg, svixClient, svixAppId); err != nil {
		return fmt.Errorf("svix.UpdateOneEndpoint: %w", err)
	}

	return nil
}

func RotateOneConfigSecret(ctx context.Context, id string, secret string, store storage.Store, svixClient *svixgo.Svix, svixAppId string) error {
	if err := svix.RotateOneEndpointSecret(ctx, id, secret, svixClient, svixAppId); err != nil {
		return fmt.Errorf("svix.RotateOneEndpointSecret: %w", err)
	}

	modifiedCount, err := store.UpdateOneConfigSecret(ctx, id, secret)
	if err != nil {
		return fmt.Errorf("store.UpdateOneConfigSecret: %w", err)
	} else if modifiedCount == 0 {
		return ErrConfigNotFound
	}

	return nil
}
