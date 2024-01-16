package core

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"golang.org/x/mod/semver"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Controller[T client.Object] func(ctx Context, req T) error

func ForModule[T Module](underlyingController Controller[T]) Controller[T] {
	return func(ctx Context, t T) error {
		err := underlyingController(ctx, t)
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

		if semver.IsValid(t.GetVersion()) {
			return ValidateInstalledVersion(ctx, t)
		}

		return nil
	}
}

func ForStackDependency[T Dependent](ctrl Controller[T]) Controller[T] {
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
		}

		return ctrl(ctx, t)
	}
}
