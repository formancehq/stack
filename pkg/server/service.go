package server

import (
	"context"

	"github.com/formancehq/go-libs/sharedapi"
	"github.com/formancehq/go-libs/sharedlogging"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/pkg/errors"
)

var (
	ErrConfigNotFound    = errors.New("config not found")
	ErrConfigNotModified = errors.New("config not modified")
	ErrConfigNotDeleted  = errors.New("config not deleted")
)

func insertOneConfig(ctx context.Context, cfg webhooks.ConfigUser, store storage.Store) (string, error) {
	id, err := store.InsertOneConfig(ctx, cfg)
	if err != nil {
		return "", errors.Wrap(err, "store.insertOneConfig")
	}

	sharedlogging.GetLogger(ctx).Debug("insertOneConfig: id: ", id)
	return id, nil
}

func deleteOneConfig(ctx context.Context, id string, store storage.Store) error {
	if _, err := findConfig(ctx, store, id); err != nil {
		return errors.Wrap(err, "findConfig")
	}

	if deletedCount, err := store.DeleteOneConfig(ctx, id); err != nil {
		return errors.Wrap(err, "store.deleteOneConfig")
	} else if deletedCount == 0 {
		return ErrConfigNotDeleted
	}

	sharedlogging.GetLogger(ctx).Debug("deleteOneConfig: id: ", id)
	return nil
}

func updateOneConfigActivation(ctx context.Context, active bool, id string, store storage.Store) (sharedapi.Cursor[webhooks.Config], error) {
	if _, err := findConfig(ctx, store, id); err != nil {
		return sharedapi.Cursor[webhooks.Config]{}, errors.Wrap(err, "findConfig")
	}

	if _, modifiedCount, _, _, err := store.UpdateOneConfigActivation(ctx, id, active); err != nil {
		return sharedapi.Cursor[webhooks.Config]{}, errors.Wrap(err, "store.updateOneConfigActivation")
	} else if modifiedCount == 0 {
		return sharedapi.Cursor[webhooks.Config]{}, ErrConfigNotModified
	}

	cursor, err := findConfig(ctx, store, id)
	if err != nil {
		return sharedapi.Cursor[webhooks.Config]{}, errors.Wrap(err, "findConfig")
	}

	sharedlogging.GetLogger(ctx).Debugf("updateOneConfigActivation (%v): id: %s", active, id)
	return cursor, nil
}

func changeOneConfigSecret(ctx context.Context, id, secret string, store storage.Store) (sharedapi.Cursor[webhooks.Config], error) {
	if _, err := findConfig(ctx, store, id); err != nil {
		return sharedapi.Cursor[webhooks.Config]{}, errors.Wrap(err, "findConfig")
	}

	if _, modifiedCount, _, _, err := store.UpdateOneConfigSecret(ctx, id, secret); err != nil {
		return sharedapi.Cursor[webhooks.Config]{}, errors.Wrap(err, "store.UpdateOneConfigSecret")
	} else if modifiedCount == 0 {
		return sharedapi.Cursor[webhooks.Config]{}, ErrConfigNotModified
	}

	cursor, err := findConfig(ctx, store, id)
	if err != nil {
		return sharedapi.Cursor[webhooks.Config]{}, errors.Wrap(err, "findConfig")
	}

	sharedlogging.GetLogger(ctx).Debug("changeOneConfigSecret: id: ", id)
	return cursor, nil
}

func findConfig(ctx context.Context, store storage.Store, id string) (cur sharedapi.Cursor[webhooks.Config], err error) {
	if cur, err = store.FindManyConfigs(ctx, map[string]any{webhooks.KeyID: id}); err != nil {
		return sharedapi.Cursor[webhooks.Config]{}, errors.Wrap(err, "store.FindManyConfigs")
	} else if len(cur.Data) == 0 {
		return sharedapi.Cursor[webhooks.Config]{}, ErrConfigNotFound
	}

	return
}
