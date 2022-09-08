package storage

import (
	"context"

	"github.com/numary/go-libs/sharedapi"
	webhooks "github.com/numary/webhooks/pkg"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store interface {
	FindManyConfigs(ctx context.Context, filter map[string]any) (sharedapi.Cursor[webhooks.Config], error)
	InsertOneConfig(ctx context.Context, cfg webhooks.ConfigUser) (string, error)
	DeleteOneConfig(ctx context.Context, id string) (int64, error)
	UpdateOneConfigActivation(ctx context.Context, id string, active bool) (*webhooks.Config, int64, error)
	UpdateOneConfigSecret(ctx context.Context, id, secret string) (int64, error)
	InsertOneRequest(ctx context.Context, req webhooks.Request) (primitive.ObjectID, error)
	Close(ctx context.Context) error
}
