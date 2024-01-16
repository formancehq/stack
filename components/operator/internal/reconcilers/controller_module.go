package reconcilers

import (
	"github.com/formancehq/operator/internal/core"
	"golang.org/x/mod/semver"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	if semver.IsValid(t.GetVersion()) {
		return core.ValidateInstalledVersion(ctx, t)
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
