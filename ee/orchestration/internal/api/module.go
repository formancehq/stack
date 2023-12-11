package api

import (
	"go.uber.org/fx"
)

func TagVersion() fx.Annotation {
	return fx.ResultTags(`group:"apiVersions"`)
}

func NewModule() fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotate(NewRouter, fx.ParamTags(``, ``, ``, `group:"apiVersions"`))),
		fx.Provide(NewDefaultBackend),
	)
}
