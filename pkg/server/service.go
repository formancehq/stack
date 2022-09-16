package server

import (
	"context"
	"errors"
	"fmt"

	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/numary/go-libs/sharedlogging"
)

var (
	ErrConfigNotFound    = errors.New("config not found")
	ErrConfigNotModified = errors.New("config not modified")
	ErrConfigNotDeleted  = errors.New("config not deleted")
)

func insertOneConfig(ctx context.Context, cfg webhooks.ConfigUser, store storage.Store) (string, error) {
	id, err := store.InsertOneConfig(ctx, cfg)
	if err != nil {
		return "", fmt.Errorf("store.insertOneConfig: %w", err)
	}

	sharedlogging.GetLogger(ctx).Debug("insertOneConfig: id: ", id)
	return id, nil
}

func deleteOneConfig(ctx context.Context, id string, store storage.Store) error {
	if err := findConfig(ctx, store, id); err != nil {
		return err
	}

	if deletedCount, err := store.DeleteOneConfig(ctx, id); err != nil {
		return fmt.Errorf("store.deleteOneConfig: %w", err)
	} else if deletedCount == 0 {
		return ErrConfigNotDeleted
	}

	sharedlogging.GetLogger(ctx).Debug("deleteOneConfig: id: ", id)
	return nil
}

func updateOneConfigActivation(ctx context.Context, active bool, id string, store storage.Store) error {
	if err := findConfig(ctx, store, id); err != nil {
		return err
	}

	if _, modifiedCount, _, _, err := store.UpdateOneConfigActivation(ctx, id, active); err != nil {
		return fmt.Errorf("store.updateOneConfigActivation: %w", err)
	} else if modifiedCount == 0 {
		return ErrConfigNotModified
	}

	sharedlogging.GetLogger(ctx).Debug("updateOneConfigActivation (%v): id: ", active, id)
	return nil
}

func changeOneConfigSecret(ctx context.Context, id, secret string, store storage.Store) error {
	if err := findConfig(ctx, store, id); err != nil {
		return err
	}

	if _, modifiedCount, _, _, err := store.UpdateOneConfigSecret(ctx, id, secret); err != nil {
		return fmt.Errorf("store.UpdateOneConfigSecret: %w", err)
	} else if modifiedCount == 0 {
		return ErrConfigNotModified
	}

	sharedlogging.GetLogger(ctx).Debug("changeOneConfigSecret: id: ", id)
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
