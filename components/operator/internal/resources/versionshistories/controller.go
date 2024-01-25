package versionshistories

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
)

// +kubebuilder:rbac:groups=formance.com,resources=versionshistories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=versionshistories/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=versionshistories/finalizers,verbs=update
// +kubebuilder:rbac:groups=formance.com,resources=versions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=versions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=versions/finalizers,verbs=update

func init() {
	Init(
		WithSimpleIndex[*v1beta1.VersionsHistory](".spec.module", func(t *v1beta1.VersionsHistory) string {
			return t.Spec.Module
		}),
	)
}
