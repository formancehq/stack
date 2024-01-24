package core

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	pkgError "github.com/pkg/errors"
	"golang.org/x/mod/semver"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Controller[T client.Object] func(ctx Context, req T) error

func ForModule[T v1beta1.Module](underlyingController func(ctx Context, stack *v1beta1.Stack, t T, version string) error) func(ctx Context, stack *v1beta1.Stack, t T) error {
	return func(ctx Context, stack *v1beta1.Stack, t T) error {

		moduleVersion, err := GetModuleVersion(ctx, stack, t)
		if err != nil {
			return err
		}

		err = underlyingController(ctx, stack, t, moduleVersion)
		if err != nil {
			return err
		}

		for _, condition := range t.GetConditions() {
			if condition.ObservedGeneration != t.GetGeneration() {
				continue
			}

			if condition.Status != metav1.ConditionTrue {
				return ErrPending
			}
		}

		if semver.IsValid(moduleVersion) {
			return ValidateInstalledVersion(ctx, t, moduleVersion)
		}

		return nil
	}
}

func ForStackDependency[T v1beta1.Dependent](ctrl func(ctx Context, stack *v1beta1.Stack, t T) error) func(ctx Context, t T) error {
	return func(ctx Context, t T) error {
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
			if stack.GetLabels()[SkipLabel] == "true" {
				return nil
			}
		}

		return ctrl(ctx, stack, t)
	}
}

func ForReadier[T v1beta1.Object](controller Controller[T]) Controller[T] {
	return func(ctx Context, object T) error {
		setStatus := func(err error) {
			if err != nil {
				object.SetReady(false)
				object.SetError(err.Error())
			} else {
				object.SetReady(true)
				object.SetError("")
			}
		}

		var reconcilerError error
		err := controller(ctx, object)
		if err != nil {
			setStatus(err)
			if !pkgError.Is(err, ErrPending) &&
				!pkgError.Is(err, ErrDeleted) {
				reconcilerError = err
			}
		} else {
			setStatus(nil)
		}

		return reconcilerError
	}
}
