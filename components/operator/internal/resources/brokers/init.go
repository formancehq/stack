package brokers

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	v1 "k8s.io/api/batch/v1"
)

func init() {
	core.Init(
		core.WithResourceReconciler(Reconcile,
			core.WithFinalizer[*v1beta1.Broker]("clear", deleteBroker),
			core.WithOwn[*v1beta1.Broker](&v1.Job{}),
		),
	)
}
