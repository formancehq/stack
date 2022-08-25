package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/pkg/model"
	"github.com/numary/webhooks/pkg/storage"
	"github.com/numary/webhooks/pkg/webhooks"
)

var (
	ErrConfigNotFound    = errors.New("config not found")
	ErrConfigNotModified = errors.New("config not modified")
)

func InsertOneConfig(ctx context.Context, cfg model.Config, store storage.Store, engine webhooks.Engine) (string, error) {
	var id string
	var err error
	if id, err = store.InsertOneConfig(ctx, cfg); err != nil {
		return "", fmt.Errorf("store.InsertOneConfig: %w", err)
	}

	if err := engine.InsertOneConfig(ctx, id, cfg); err != nil {
		if _, err = store.DeleteOneConfig(ctx, id); err != nil {
			return "", fmt.Errorf("store.DeleteOneConfig: %w", err)
		}
		return "", fmt.Errorf("engine.InsertOneConfig: %w", err)
	}

	sharedlogging.GetLogger(ctx).Debug("service.InsertOneConfig: id: ", id)
	return id, nil
}

func DeleteOneConfig(ctx context.Context, id string, store storage.Store, engine webhooks.Engine) error {
	var cfg model.Config
	if cur, err := store.FindManyConfigs(ctx, map[string]any{"_id": id}); err != nil {
		return fmt.Errorf("sotre.FindManyConfigs: %w", err)
	} else if len(cur.Data) == 0 {
		return ErrConfigNotFound
	} else {
		cfg = cur.Data[0].Config
	}

	if err := engine.DeleteOneConfig(ctx, id); err != nil {
		return fmt.Errorf("engine.DeleteOneConfig: %w", err)
	}

	if deletedCount, err := store.DeleteOneConfig(ctx, id); err != nil {
		if err := engine.InsertOneConfig(ctx, id, cfg); err != nil {
			return fmt.Errorf("engine.InsertOneConfig: %w", err)
		}
		return fmt.Errorf("store.DeleteOneConfig: %w", err)
	} else if deletedCount == 0 {
		return ErrConfigNotFound
	}

	sharedlogging.GetLogger(ctx).Debug("service.DeleteOneConfig: id: ", id)
	return nil
}

func ActivateOneConfig(ctx context.Context, active bool, id string, store storage.Store, engine webhooks.Engine) error {
	if cur, err := store.FindManyConfigs(ctx, map[string]any{"_id": id}); err != nil {
		return fmt.Errorf("sotre.FindManyConfigs: %w", err)
	} else if len(cur.Data) == 0 {
		return ErrConfigNotFound
	}

	updatedCfg, modifiedCount, err := store.UpdateOneConfigActivation(ctx, id, active)
	if err != nil {
		return fmt.Errorf("store.UpdateOneConfigActivation: %w", err)
	} else if updatedCfg == nil {
		return ErrConfigNotFound
	} else if modifiedCount == 0 {
		return ErrConfigNotModified
	}

	if err := engine.UpdateOneConfig(ctx, id, updatedCfg); err != nil {
		if _, _, err := store.UpdateOneConfigActivation(ctx, id, !active); err != nil {
			return fmt.Errorf("store.UpdateOneConfigActivation: %w", err)
		}
		return fmt.Errorf("engine.UpdateOneConfig: %w", err)
	}

	return nil
}

func RotateOneConfigSecret(ctx context.Context, id, secret string, store storage.Store, engine webhooks.Engine) error {
	var currentSecret string
	if cur, err := store.FindManyConfigs(ctx, map[string]any{"_id": id}); err != nil {
		return fmt.Errorf("sotre.FindManyConfigs: %w", err)
	} else if len(cur.Data) == 0 {
		return ErrConfigNotFound
	} else {
		currentSecret = cur.Data[0].Secret
	}

	if err := engine.RotateOneConfigSecret(ctx, id, secret); err != nil {
		return fmt.Errorf("engine.RotateOneConfigSecret: %w", err)
	}

	if modifiedCount, err := store.UpdateOneConfigSecret(ctx, id, secret); err != nil {
		if err := engine.RotateOneConfigSecret(ctx, id, currentSecret); err != nil {
			return fmt.Errorf("engine.RotateOneConfigSecret: %w", err)
		}
		return fmt.Errorf("store.UpdateOneConfigSecret: %w", err)
	} else if modifiedCount == 0 {
		return ErrConfigNotFound
	}

	return nil
}
