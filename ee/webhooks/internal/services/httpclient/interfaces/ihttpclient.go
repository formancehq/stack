package interfaces

import (
	"context"

	"github.com/formancehq/webhooks/internal/commons"
)

type IHTTPClient interface {
	Call(context.Context, *commons.Hook, *commons.Attempt, bool) (int, error)
}
