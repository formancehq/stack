package connectors

import (
	"context"

	"github.com/google/uuid"
)

type webhookIDKey struct{}

var _webhookIDKey = webhookIDKey{}

func ContextWithWebhookID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, _webhookIDKey, id)
}

func WebhookIDFromContext(ctx context.Context) uuid.UUID {
	return ctx.Value(_webhookIDKey).(uuid.UUID)
}
