package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/numary/go-libs/sharedlogging"
	webhooks "github.com/numary/webhooks/pkg"
	"github.com/numary/webhooks/pkg/storage"
)

var (
	ErrConfigNotFound    = errors.New("config not found")
	ErrConfigNotModified = errors.New("config not modified")
	ErrConfigNotDeleted  = errors.New("config not deleted")
)

func InsertOneConfig(ctx context.Context, cfg webhooks.ConfigUser, store storage.Store) (string, error) {
	id, err := store.InsertOneConfig(ctx, cfg)
	if err != nil {
		return "", fmt.Errorf("store.InsertOneConfig: %w", err)
	}

	sharedlogging.GetLogger(ctx).Debug("service.InsertOneConfig: id: ", id)
	return id, nil
}

func DeleteOneConfig(ctx context.Context, id string, store storage.Store) error {
	if err := findConfig(ctx, store, id); err != nil {
		return err
	}

	if deletedCount, err := store.DeleteOneConfig(ctx, id); err != nil {
		return fmt.Errorf("store.DeleteOneConfig: %w", err)
	} else if deletedCount == 0 {
		return ErrConfigNotDeleted
	}

	sharedlogging.GetLogger(ctx).Debug("service.DeleteOneConfig: id: ", id)
	return nil
}

func UpdateOneConfigActivation(ctx context.Context, active bool, id string, store storage.Store) error {
	if err := findConfig(ctx, store, id); err != nil {
		return err
	}

	if _, modifiedCount, _, _, err := store.UpdateOneConfigActivation(ctx, id, active); err != nil {
		return fmt.Errorf("store.UpdateOneConfigActivation: %w", err)
	} else if modifiedCount == 0 {
		return ErrConfigNotModified
	}

	return nil
}

func RotateOneConfigSecret(ctx context.Context, id, secret string, store storage.Store) error {
	if err := findConfig(ctx, store, id); err != nil {
		return err
	}

	if _, modifiedCount, _, _, err := store.UpdateOneConfigSecret(ctx, id, secret); err != nil {
		return fmt.Errorf("store.UpdateOneConfigSecret: %w", err)
	} else if modifiedCount == 0 {
		return ErrConfigNotModified
	}

	return nil
}

func findConfig(ctx context.Context, store storage.Store, id string) error {
	if cur, err := store.FindManyConfigs(ctx, map[string]any{webhooks.KeyID: id}); err != nil {
		return fmt.Errorf("store.FindManyConfigs: %w", err)
	} else if len(cur.Data) == 0 {
		return ErrConfigNotFound
	}

	return nil
}
