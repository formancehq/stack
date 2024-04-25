package brokers

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	v1 "k8s.io/api/batch/v1"
)

//+kubebuilder:rbac:groups=formance.com,resources=brokers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=brokers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=brokers/finalizers,verbs=update

func init() {
	core.Init(
		core.WithResourceReconciler(Reconcile,
			core.WithFinalizer[*v1beta1.Broker]("clear", deleteBroker),
			core.WithOwn[*v1beta1.Broker](&v1.Job{}),
		),
	)
}
