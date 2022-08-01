package storage

import (
	"context"

	"github.com/numary/go-libs/sharedapi"
	"github.com/numary/webhooks-cloud/pkg/model"
)

type Store interface {
	FindAllConfigs(ctx context.Context) (sharedapi.Cursor[model.ConfigInserted], error)
	InsertOneConfig(ctx context.Context, config model.Config) (string, error)
	DropConfigsCollection(ctx context.Context) error
	Close(ctx context.Context) error
}
