package storage

import (
	"context"

	"github.com/numary/go-libs/sharedapi"
	webhooks "github.com/numary/webhooks/pkg"
)

type Store interface {
	FindManyConfigs(ctx context.Context, filter map[string]any) (sharedapi.Cursor[webhooks.Config], error)
	InsertOneConfig(ctx context.Context, cfg webhooks.ConfigUser) (string, error)
	DeleteOneConfig(ctx context.Context, id string) (int64, error)
	UpdateOneConfigActivation(ctx context.Context, id string, active bool) (matchedCount, modifiedCount, upsertedCount int64, upsertedID any, err error)
	UpdateOneConfigSecret(ctx context.Context, id, secret string) (matchedCount, modifiedCount, upsertedCount int64, upsertedID any, err error)
	FindManyAttempts(ctx context.Context, filter map[string]any) (sharedapi.Cursor[webhooks.Attempt], error)
	FindDistinctWebhookIDs(ctx context.Context, filter map[string]any) ([]string, error)
	UpdateManyAttemptsStatus(ctx context.Context, webhookID string, status string) (matchedCount, modifiedCount, upsertedCount int64, upsertedID any, err error)
	InsertOneAttempt(ctx context.Context, att webhooks.Attempt) (insertedID string, err error)
	Close(ctx context.Context) error
}
