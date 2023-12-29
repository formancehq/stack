package reconcilers

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler[T client.Object] struct {
	Controller Controller[T]
	Manager    Manager
}

func (r *Reconciler[T]) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	var t T

	log := log.FromContext(ctx, fmt.Sprintf("%T", t), req.NamespacedName)
	log.Info("Starting reconciliation")

	t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
	if err := r.Manager.GetClient().Get(ctx, types.NamespacedName{
		Name: req.Name,
	}, t); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	cp := t.DeepCopyObject().(T)
	if err := r.Controller.Reconcile(struct {
		context.Context
		Manager
	}{
		Context: ctx,
		Manager: r.Manager,
	}, t); err != nil {
		return ctrl.Result{}, err
	}

	if !equality.Semantic.DeepEqual(cp, t) {
		patch := client.MergeFrom(cp)
		if err := r.Manager.GetClient().Status().Patch(ctx, t, patch); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler[T]) SetupWithManager(mgr Manager) error {
	r.Manager = mgr
	builder, err := r.Controller.SetupWithManager(mgr)
	if err != nil {
		return err
	}
	return builder.Complete(r)
}

func New[T client.Object](ctrl Controller[T]) *Reconciler[T] {
	return &Reconciler[T]{
		Controller: ctrl,
	}
}
