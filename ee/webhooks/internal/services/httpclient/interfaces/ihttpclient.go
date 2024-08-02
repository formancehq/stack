package interfaces

import (
	"context"

	"github.com/formancehq/webhooks/internal/models"
)

type IHTTPClient interface {
	Call(context.Context, *models.Hook, *models.Attempt, bool) (int, error)
}
