package reconcilers

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

type StackDependency[T core.Dependent] struct {
	core.Controller[T]
}

func (r *StackDependency[T]) Reconcile(ctx core.Context, t T) error {
	stack := &v1beta1.Stack{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: t.GetStack(),
	}, stack); err != nil {
		if errors.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	} else {
		if stack.Spec.Disabled {
			return nil
		}
	}

	return r.Controller.Reconcile(ctx, t)
}

func ForStackDependency[T core.Dependent](ctrl core.Controller[T]) *Reconciler[T] {
	return New[T](&StackDependency[T]{
		Controller: ctrl,
	})
}
