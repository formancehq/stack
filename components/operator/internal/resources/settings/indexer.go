package settings

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
)

func init() {
	core.Init(
		core.WithSimpleIndex[*v1beta1.Settings]("key", func(t *v1beta1.Settings) string {
			return t.Spec.Key
		}),
	)
}
