package core

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type ObjectController[T client.Object] func(ctx Context, req T) error

func ForObjectController[T v1beta1.Object](controller ObjectController[T]) ObjectController[T] {
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
			if !IsApplicationError(err) {
				reconcilerError = err
			}
		} else {
			setStatus(nil)
		}

		return reconcilerError
	}
}

type StackDependentObjectController[T v1beta1.Dependent] func(ctx Context, stack *v1beta1.Stack, req T) error

func ForStackDependency[T v1beta1.Dependent](ctrl StackDependentObjectController[T]) ObjectController[T] {
	return func(ctx Context, t T) error {
		stack := &v1beta1.Stack{}
		if err := ctx.GetClient().Get(ctx, types.NamespacedName{
			Name: t.GetStack(),
		}, stack); err != nil {
			if apierrors.IsNotFound(err) {
				return NewStackNotFoundError()
			} else {
				return err
			}
		} else {
			if stack.GetAnnotations()[v1beta1.SkipLabel] == "true" {
				return nil
			}
		}

		return ctrl(ctx, stack, t)
	}
}

func ForResource[T v1beta1.Dependent](ctrl StackDependentObjectController[T]) StackDependentObjectController[T] {
	return func(ctx Context, stack *v1beta1.Stack, req T) error {
		// notes(gfyrag): We attach the database object to the stack as owner
		// this way, even if the controller is removed, the Database object will not be removed until
		// the stack is removed. This allows us to remove a module and re-add it later if we wants.
		hasOwnerReferenceOnStack, err := HasOwnerReference(ctx, stack, req)
		if err != nil {
			return err
		}
		if !hasOwnerReferenceOnStack {
			patch := client.MergeFrom(req.DeepCopyObject().(T))

			gvk, err := apiutil.GVKForObject(stack, ctx.GetScheme())
			if err != nil {
				return err
			}

			ownerReferences := req.GetOwnerReferences()
			ownerReferences = append(ownerReferences, metav1.OwnerReference{
				APIVersion: gvk.GroupVersion().String(),
				Kind:       gvk.Kind,
				UID:        stack.GetUID(),
				Name:       stack.GetName(),
			})
			req.SetOwnerReferences(ownerReferences)

			if err := ctx.GetClient().Patch(ctx, req, patch); err != nil {
				return err
			}
		}

		return ctrl(ctx, stack, req)
	}
}

type ModuleController[T v1beta1.Module] func(ctx Context, stack *v1beta1.Stack, req T, version string) error

func ForModule[T v1beta1.Module](underlyingController ModuleController[T]) StackDependentObjectController[T] {
	return func(ctx Context, stack *v1beta1.Stack, t T) error {

		moduleVersion, err := GetModuleVersion(ctx, stack, t)
		if err != nil {
			return err
		}

		hasOwnerReference, err := HasOwnerReference(ctx, stack, t)
		if err != nil {
			return err
		}

		if !hasOwnerReference {
			patch := client.MergeFrom(t.DeepCopyObject().(T))
			if err := controllerutil.SetOwnerReference(stack, t, ctx.GetScheme()); err != nil {
				return err
			}
			if err := ctx.GetClient().Patch(ctx, t, patch); err != nil {
				return errors.Wrap(err, "patching object to add owner reference on stack")
			}
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
				return NewPendingError()
			}
		}

		return nil
	}
}
