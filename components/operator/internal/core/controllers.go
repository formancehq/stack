package core

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ObjectController[T client.Object] func(ctx Context, reconcilerOptions *ReconcilerOptions[T], req T) error

func ForObjectController[T v1beta1.Object](controller ObjectController[T]) ObjectController[T] {
	return func(ctx Context, reconcilerOptions *ReconcilerOptions[T], object T) error {
		setStatus := func(err error) {
			if err != nil {
				object.SetReady(false)
				object.SetError(err.Error())
			} else {
				object.SetReady(true)
				object.SetError("Up to date")
			}
		}

		var reconcilerError error
		err := controller(ctx, reconcilerOptions, object)
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

type StackDependentObjectController[T v1beta1.Dependent] func(ctx Context, stack *v1beta1.Stack, reconcilerOptions *ReconcilerOptions[T], req T) error

func ForStackDependency[T v1beta1.Dependent](ctrl StackDependentObjectController[T], allowDeleted bool) ObjectController[T] {
	return func(ctx Context, reconcilerOptions *ReconcilerOptions[T], t T) error {
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
				patch := client.MergeFrom(t.DeepCopyObject().(T))
				annotations := t.GetAnnotations()
				if annotations == nil {
					annotations = map[string]string{}
				}
				annotations[v1beta1.SkippedLabel] = "true"
				t.SetAnnotations(annotations)
				return ctx.GetClient().Patch(ctx, t, patch)
			}
		}
		annotations := t.GetAnnotations()
		if annotations != nil {
			patch := client.MergeFrom(t.DeepCopyObject().(T))
			_, ok := annotations[v1beta1.SkippedLabel]
			if ok {
				delete(annotations, v1beta1.SkippedLabel)
				t.SetAnnotations(annotations)
				if err := ctx.GetClient().Patch(ctx, t, patch); err != nil {
					return err
				}
			}
		}

		if !allowDeleted {
			if !stack.GetDeletionTimestamp().IsZero() {
				return NewStackNotFoundError()
			}
		}

		return ctrl(ctx, stack, reconcilerOptions, t)
	}
}

type ModuleController[T v1beta1.Module] func(ctx Context, stack *v1beta1.Stack, reconcilerOptions *ReconcilerOptions[T], req T, version string) error

func ForModule[T v1beta1.Module](underlyingController ModuleController[T]) StackDependentObjectController[T] {
	return func(ctx Context, stack *v1beta1.Stack, reconcilerOptions *ReconcilerOptions[T], t T) error {

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
			log.FromContext(ctx).Info("Add owner reference on stack")
		}

		if stack.Spec.Disabled {
			// notes(gfyrag): When disabling a stack, we remove all owned objects for modules.
			// Owned objects must be controlled by the module.
			// if not, they will not be automatically removed on stack removal.
			// resources objects (like Database and BrokerTopic) are not removed since we could re-enable the stack later.
			if err := removeAllModulesOwnedObjects(ctx, t, reconcilerOptions.Owns); err != nil {
				return err
			}
		} else {
			err = underlyingController(ctx, stack, reconcilerOptions, t, moduleVersion)
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
		}

		return nil
	}
}

func removeAllModulesOwnedObjects(ctx Context, owner client.Object, owns map[client.Object][]builder.OwnsOption) error {
	for object := range owns {
		if _, ok := object.(v1beta1.Resource); ok {
			// Resources must not be deleted
			continue
		}

		gvk, err := apiutil.GVKForObject(object, ctx.GetScheme())
		if err != nil {
			return err
		}

		list := &unstructured.UnstructuredList{}
		list.SetGroupVersionKind(gvk)
		if err := ctx.GetClient().List(ctx, list); err != nil {
			return err
		}

		for _, item := range list.Items {
			hasControllerReference, err := HasControllerReference(ctx, owner, &item)
			if err != nil {
				return err
			}
			if hasControllerReference {
				if err := ctx.GetClient().Delete(ctx, &item); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
