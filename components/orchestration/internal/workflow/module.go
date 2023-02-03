package workflow

import (
	"github.com/uptrace/bun"
	"go.temporal.io/sdk/client"
	"go.uber.org/fx"
)

func NewModule(taskQueue string) fx.Option {
	return fx.Provide(func(db *bun.DB, temporalClient client.Client) *Manager {
		return NewManager(db, temporalClient, taskQueue)
	})
}
