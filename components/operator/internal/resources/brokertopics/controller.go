package brokertopics

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"k8s.io/apimachinery/pkg/types"
)

func Reconcile(ctx core.Context, stack *v1beta1.Stack, topic *v1beta1.BrokerTopic) error {

	if len(topic.GetOwnerReferences()) == 1 { // Remains only the stack
		return ctx.GetClient().Delete(ctx, topic)
	}

	broker, _, err := core.CreateOrUpdate[*v1beta1.Broker](ctx, types.NamespacedName{
		Name: stack.Name,
	}, func(t *v1beta1.Broker) error {
		t.Spec.Stack = stack.Name
		return nil
	}, core.WithController[*v1beta1.Broker](ctx.GetScheme(), stack))
	if err != nil {
		return err
	}

	if !broker.Status.Ready {
		return core.NewApplicationError().WithMessage("broker not ready")
	}
	topic.Status.Ready = broker.Status.Ready

	return nil
}
