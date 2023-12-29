package reconcilers

import (
	"context"
	"github.com/formancehq/operator/v2/internal/controller/shared"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ContextualManager interface {
	context.Context
	Manager
}

type Controller[T client.Object] interface {
	Reconcile(ctx ContextualManager, req T) error
	SetupWithManager(mgr Manager) (*ctrl.Builder, error)
}

type reconciler interface {
	SetupWithManager(mgr Manager) error
}

var reconcilers []reconciler

func Register(newReconcilers ...reconciler) {
	reconcilers = append(reconcilers, newReconcilers...)
}

func Setup(mgr ctrl.Manager, platform shared.Platform) error {
	wrappedMgr := newDefaultManager(mgr, platform)
	for _, reconciler := range reconcilers {
		if err := reconciler.SetupWithManager(wrappedMgr); err != nil {
			return err
		}
	}
	return nil
}
