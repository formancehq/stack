package ledger

import (
	"github.com/formancehq/go-libs/logging"
	"go.uber.org/fx"
)

type ModuleConfiguration struct {
	NSCacheConfiguration CacheConfiguration
}

func NewFXModule(configuration ModuleConfiguration) fx.Option {
	return fx.Options(
		fx.Provide(func(
			storageDriver StorageDriver,
			listener Listener,
			logger logging.Logger,
		) *Resolver {
			options := []Option{}
			if configuration.NSCacheConfiguration.MaxCount != 0 {
				options = append(options, WithCompiler(NewCachedCompiler(
					NewDefaultCompiler(),
					configuration.NSCacheConfiguration,
				)))
			}
			return NewResolver(storageDriver, listener, options...)
		}),
	)
}
