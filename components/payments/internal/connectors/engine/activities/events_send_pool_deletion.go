package activities

import (
	"context"

	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
)

func (a Activities) EventsSendPoolDeletion(ctx context.Context, id uuid.UUID) error {
	return a.events.Publish(ctx, a.events.NewEventDeletePool(id))
}

var EventsSendPoolDeletionActivity = Activities{}.EventsSendPoolDeletion

func EventsSendPoolDeletion(ctx workflow.Context, id uuid.UUID) error {
	return executeActivity(ctx, EventsSendPoolDeletionActivity, nil, id)
}
