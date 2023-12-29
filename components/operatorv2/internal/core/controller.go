package core

import (
	"sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Controller[T client.Object] interface {
	Reconcile(ctx Context, req T) error
	SetupWithManager(mgr Manager) (*controllerruntime.Builder, error)
}
