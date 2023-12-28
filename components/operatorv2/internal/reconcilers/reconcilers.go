package reconcilers

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Controller[T client.Object] interface {
	Reconcile(ctx context.Context, req T) error
	SetupWithManager(mgr ctrl.Manager, builder *builder.Builder) error
}

type Reconciler[T client.Object] struct {
	Client     client.Client
	Scheme     *runtime.Scheme
	Controller Controller[T]
}

func (r *Reconciler[T]) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	var t T

	log := log.FromContext(ctx, fmt.Sprintf("%T", t), req.NamespacedName)
	log.Info("Starting reconciliation")

	t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name: req.Name,
	}, t); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	cp := t.DeepCopyObject().(T)
	if err := r.Controller.Reconcile(ctx, t); err != nil {
		return ctrl.Result{}, err
	}

	if !equality.Semantic.DeepEqual(cp, t) {
		patch := client.MergeFrom(cp)
		if err := r.Client.Status().Patch(ctx, t, patch); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler[T]) SetupWithManager(mgr ctrl.Manager) error {
	var t T

	t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
	builder := ctrl.NewControllerManagedBy(mgr).For(t)
	if err := r.Controller.SetupWithManager(mgr, builder); err != nil {
		return err
	}

	return builder.Complete(r)
}

func New[T client.Object](client client.Client, scheme *runtime.Scheme, ctrl Controller[T]) *Reconciler[T] {
	return &Reconciler[T]{
		Client:     client,
		Scheme:     scheme,
		Controller: ctrl,
	}
}

func SetupReconcilers(mgr ctrl.Manager, reconcilers ...interface {
	SetupWithManager(mgr ctrl.Manager) error
}) error {
	for _, reconciler := range reconcilers {
		if err := reconciler.SetupWithManager(mgr); err != nil {
			return err
		}
	}
	return nil
}
