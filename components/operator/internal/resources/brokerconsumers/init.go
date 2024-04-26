package brokerconsumers

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/brokers"
	v1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
)

func init() {
	core.Init(
		core.WithStackDependencyReconciler(Reconcile,
			core.WithOwn[*v1beta1.BrokerConsumer](&v1beta1.BrokerTopic{}, builder.MatchEveryOwner),
			core.WithOwn[*v1beta1.BrokerConsumer](&v1.Job{}),
			brokers.Watch[*v1beta1.BrokerConsumer](),
		),
	)
}
