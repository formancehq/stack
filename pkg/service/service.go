package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/pkg/model"
	"github.com/numary/webhooks/pkg/storage"
	"github.com/numary/webhooks/pkg/svix"
)

var (
	ErrConfigNotFound    = errors.New("config not found")
	ErrConfigNotModified = errors.New("config not modified")
)

func InsertOneConfig(cfg model.Config, ctx context.Context, store storage.Store, svixApp svix.App) (string, error) {
	var id string
	var err error
	if id, err = store.InsertOneConfig(ctx, cfg); err != nil {
		return "", fmt.Errorf("store.InsertOneConfig: %w", err)
	}

	if err := svix.CreateEndpoint(ctx, id, cfg, svixApp); err != nil {
		if _, err = store.DeleteOneConfig(ctx, id); err != nil {
			return "", fmt.Errorf("store.DeleteOneConfig: %w", err)
		}
		return "", fmt.Errorf("svix.CreateEndpoint: %w", err)
	}

	sharedlogging.GetLogger(ctx).Debug("service.InsertOneConfig: id: ", id)
	return id, nil
}

func DeleteOneConfig(ctx context.Context, id string, store storage.Store, svixApp svix.App) error {
	if deletedCount, err := store.DeleteOneConfig(ctx, id); err != nil {
		return fmt.Errorf("store.DeleteOneConfig: %w", err)
	} else if deletedCount == 0 {
		return ErrConfigNotFound
	}

	if err := svix.DeleteOneEndpoint(id, svixApp); err != nil {
		return fmt.Errorf("svix.DeleteOneEndpoint: %w", err)
	}

	sharedlogging.GetLogger(ctx).Debug("service.DeleteOneConfig: id: ", id)
	return nil
}

func ActivateOneConfig(active bool, ctx context.Context, id string, store storage.Store, svixApp svix.App) error {
	updatedCfg, modifiedCount, err := store.UpdateOneConfigActivation(ctx, id, active)
	if err != nil {
		return fmt.Errorf("store.UpdateOneConfigActivation: %w", err)
	} else if updatedCfg == nil {
		return ErrConfigNotFound
	} else if modifiedCount == 0 {
		return ErrConfigNotModified
	}

	if err := svix.UpdateOneEndpoint(ctx, id, updatedCfg, svixApp); err != nil {
		return fmt.Errorf("svix.UpdateOneEndpoint: %w", err)
	}

	return nil
}

func RotateOneConfigSecret(ctx context.Context, id string, secret string, store storage.Store, svixApp svix.App) error {
	if err := svix.RotateOneEndpointSecret(ctx, id, secret, svixApp); err != nil {
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
