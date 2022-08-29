package engine

import (
	"context"

	"github.com/numary/webhooks/pkg/model"
)

type Engine interface {
	InsertOneConfig(ctx context.Context, id string, cfg model.Config) error
	DeleteOneConfig(ctx context.Context, id string) error
	UpdateOneConfig(ctx context.Context, id string, cfg *model.ConfigInserted) error
	RotateOneConfigSecret(ctx context.Context, id, secret string) error
	ProcessKafkaMessage(ctx context.Context, eventType string, msgValue []byte) error
}
