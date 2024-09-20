package storage

import (
	"github.com/formancehq/go-libs/bun/bunconnect"
	"github.com/formancehq/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

func Module(cmd *cobra.Command, connectionOptions bunconnect.ConnectionOptions, configEncryptionKey string) fx.Option {
	return fx.Options(
		bunconnect.Module(connectionOptions, service.IsDebug(cmd)),
		fx.Provide(func(db *bun.DB) Storage {
			return newStorage(db, configEncryptionKey)
		}),
	)
}
