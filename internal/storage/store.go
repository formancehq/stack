package storage

import (
	"context"

	"github.com/numary/go-libs/sharedapi"
	"github.com/numary/webhooks/internal/model"
)

type Store interface {
	FindAllConfigs(ctx context.Context) (sharedapi.Cursor[model.ConfigInserted], error)
	InsertOneConfig(ctx context.Context, cfg model.Config) (string, error)
	DeleteOneConfig(ctx context.Context, id string) (int64, error)
	UpdateOneConfigActive(ctx context.Context, id string, active bool) (model.ConfigInserted, int64, error)
	UpdateOneConfigSecret(ctx context.Context, id, secret string) (int64, error)
	FindEventType(ctx context.Context, eventType string) (bool, error)
	Close(ctx context.Context) error
}
