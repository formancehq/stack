package reconcilers

import (
	"github.com/formancehq/operator/internal/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Module[T core.Module] struct {
	core.Controller[T]
}

func (r *Module[T]) Reconcile(ctx core.Context, t T) error {

	err := r.Controller.Reconcile(ctx, t)
	if err != nil {
		return err
	}

	for _, condition := range t.GetConditions() {
		if condition.ObservedGeneration != t.GetGeneration() {
			continue
		}

		if condition.Status != metav1.ConditionTrue {
			return core.ErrPending
		}
	}

	patch := client.MergeFrom(t.DeepCopyObject().(T))
	if updated := core.ValidateInstalledVersion(t); updated {
		if err := ctx.GetClient().Patch(ctx, t, patch); err != nil {
			return err
		}
	}

	return nil
}

func ForModule[T core.Module](ctrl core.Controller[T]) *Reconciler[T] {
	return New[T](&StackDependency[T]{
		Controller: &Module[T]{
			Controller: ctrl,
		},
	})
}
