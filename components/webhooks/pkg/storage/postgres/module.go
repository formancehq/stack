package postgres

import (
	"github.com/formancehq/webhooks/pkg/storage"
	"go.uber.org/fx"
)

func NewModule(dsn string) fx.Option {
	return fx.Provide(func() (storage.Store, error) {
		return NewStore(dsn)
	})
}
