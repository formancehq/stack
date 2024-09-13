package storage

import (
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

func Module(connectionOptions bunconnect.ConnectionOptions, configEncryptionKey string) fx.Option {
	return fx.Options(
		bunconnect.Module(connectionOptions),
		fx.Provide(func(db *bun.DB) Storage {
			return newStorage(db, configEncryptionKey)
		}),
	)
}
