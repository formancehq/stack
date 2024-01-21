package versionshistories

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
)

func init() {
	Init(
		WithSimpleIndex[*v1beta1.VersionsHistory](".spec.module", func(t *v1beta1.VersionsHistory) string {
			return t.Spec.Module
		}),
	)
}
