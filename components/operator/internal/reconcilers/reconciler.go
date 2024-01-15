package reconcilers

import (
	"context"
	"fmt"
	pkgError "github.com/pkg/errors"
	"reflect"

	"github.com/formancehq/operator/internal/core"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Reconciler[T client.Object] struct {
	Controller core.Controller[T]
	Manager    core.Manager
}

func (r *Reconciler[T]) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	var object T

	log := log.FromContext(ctx, fmt.Sprintf("%T", object), req.NamespacedName)
	log.Info("Starting reconciliation")

	object = reflect.New(reflect.TypeOf(object).Elem()).Interface().(T)
	if err := r.Manager.GetClient().Get(ctx, types.NamespacedName{
		Name: req.Name,
	}, object); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	setStatus := func(err error) {
		if s, ok := reflect.ValueOf(object).
			Elem().
			FieldByName("Status").
			Addr().
			Interface().(interface {
			SetReady(bool)
			SetError(string)
			EvalReadiness(generation int64)
		}); ok {
			if err != nil {
				s.SetReady(false)
				s.SetError(err.Error())
			} else {
				s.EvalReadiness(object.GetGeneration())
			}
		}
	}

	cp := object.DeepCopyObject().(T)
	patch := client.MergeFrom(cp)

	var reconcilerError error
	err := r.Controller.Reconcile(struct {
		context.Context
		core.Manager
	}{
		Context: ctx,
		Manager: r.Manager,
	}, object)
	if err != nil {
		setStatus(err)
		if !pkgError.Is(reconcilerError, core.ErrPending) &&
			!pkgError.Is(reconcilerError, core.ErrDeleted) {
			reconcilerError = err
		}
	} else {
		setStatus(nil)
	}

	if err := r.Manager.GetClient().Status().Patch(ctx, object, patch); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, reconcilerError
}

func (r *Reconciler[T]) SetupWithManager(mgr core.Manager) error {
	r.Manager = mgr
	builder, err := r.Controller.SetupWithManager(mgr)
	if err != nil {
		return err
	}
	return builder.Complete(r)
}

func New[T client.Object](ctrl core.Controller[T]) *Reconciler[T] {
	return &Reconciler[T]{
		Controller: ctrl,
	}
}
