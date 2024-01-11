package reconcilers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/formancehq/operator/internal/core"
	pkgError "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/equality"
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
			SetStatus(status bool, error string)
		}); ok {
			if err == nil {
				s.SetStatus(true, "")
			} else {
				s.SetStatus(false, err.Error())
			}
		}
	}

	cp := object.DeepCopyObject().(T)
	reconciliationError := r.Controller.Reconcile(struct {
		context.Context
		core.Manager
	}{
		Context: ctx,
		Manager: r.Manager,
	}, object)
	if reconciliationError != nil {
		setStatus(reconciliationError)
	} else {
		setStatus(nil)
	}

	if !equality.Semantic.DeepEqual(cp, object) {
		patch := client.MergeFrom(cp)
		if err := r.Manager.GetClient().Status().Patch(ctx, object, patch); err != nil {
			return ctrl.Result{}, err
		}
	}

	if !pkgError.Is(reconciliationError, core.ErrPending) &&
		!pkgError.Is(reconciliationError, core.ErrDeleted) {
		return ctrl.Result{}, reconciliationError
	}

	return ctrl.Result{}, nil
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
