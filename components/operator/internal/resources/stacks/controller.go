package stacks

import (
	"context"
	"fmt"
	"reflect"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=stacks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=stacks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=stacks/finalizers,verbs=update
// +kubebuilder:rbac:groups=formance.com,resources=versions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=versions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=versions/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack) error {
	_, _, err := CreateOrUpdate[*corev1.Namespace](ctx, types.NamespacedName{
		Name: stack.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

func Clean(ctx Context, t *v1beta1.Stack) error {
	logger := log.FromContext(ctx)
	logger = logger.WithValues("stack", t.Name)
	logger.Info("Clean stack")

	if err := deleteModules(ctx, t, logger); err != nil {
		return err
	}

	if err := deleteResources(ctx, t, logger); err != nil {
		return err
	}

	logger.Info("All dependencies removed, remove namespace")
	ns := &corev1.Namespace{}
	ns.SetName(t.Name)
	if err := ctx.GetClient().Delete(ctx, ns); err != nil {
		return err
	}

	logger.Info("Stack cleaned.")

	return nil
}

func deleteModules(ctx Context, stack *v1beta1.Stack, logger logr.Logger) error {
	stillExistingModules := false
	for _, rtype := range ctx.GetScheme().AllKnownTypes() {
		v := reflect.New(rtype).Interface()
		module, ok := v.(v1beta1.Module)
		if !ok {
			continue
		}

		gvk, err := apiutil.GVKForObject(module, ctx.GetScheme())
		if err != nil {
			return err
		}

		l := &unstructured.UnstructuredList{}
		l.SetGroupVersionKind(gvk)
		if err := ctx.GetClient().List(ctx, l, client.MatchingFields{
			"stack": stack.Name,
		}); err != nil {
			return err
		}

		stillExistingModules = stillExistingModules || len(l.Items) > 0

		for _, item := range l.Items {
			if item.GetDeletionTimestamp().IsZero() {
				logger.Info(fmt.Sprintf("Delete module %s", item.GetName()))
				if err := ctx.GetClient().Delete(ctx, &item); client.IgnoreNotFound(err) != nil {
					return err
				}
			}
		}
	}

	if stillExistingModules {
		logger.Info("Still pending modules")
		return NewPendingError()
	}

	return nil
}

func deleteResources(ctx Context, stack *v1beta1.Stack, logger logr.Logger) error {
	stillExistingResources := false
	for _, rtype := range ctx.GetScheme().AllKnownTypes() {
		v := reflect.New(rtype).Interface()
		resource, ok := v.(v1beta1.Resource)
		if !ok {
			continue
		}
		gvk, err := apiutil.GVKForObject(resource, ctx.GetScheme())
		if err != nil {
			return err
		}

		l := &unstructured.UnstructuredList{}
		l.SetGroupVersionKind(gvk)
		if err := ctx.GetClient().List(ctx, l, client.MatchingFields{
			"stack": stack.Name,
		}); err != nil {
			return err
		}

		stillExistingResources = stillExistingResources || len(l.Items) > 0

		for _, item := range l.Items {
			if err := ctx.GetClient().Delete(ctx, &item); client.IgnoreNotFound(err) != nil {
				return err
			}
		}
	}

	if stillExistingResources {
		logger.Info("Still pending resources")
		return NewPendingError()
	}

	return nil
}

func init() {
	Init(
		WithSimpleIndex[*v1beta1.Stack](".spec.versionsFromFile", func(t *v1beta1.Stack) string {
			return t.Spec.VersionsFromFile
		}),
		WithStdReconciler(Reconcile,
			WithOwn[*v1beta1.Stack](&corev1.Namespace{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})),
			WithRaw[*v1beta1.Stack](func(ctx Context, b *builder.Builder) error {
				for _, rtype := range ctx.GetScheme().AllKnownTypes() {
					v := reflect.New(rtype).Interface()

					switch v.(type) {
					case v1beta1.Module:
						b.Watches(v.(v1beta1.Module), handler.EnqueueRequestsFromMapFunc(func(watchContext context.Context, object client.Object) []reconcile.Request {
							return []reconcile.Request{{
								NamespacedName: types.NamespacedName{
									Name: object.(v1beta1.Module).GetStack(),
								},
							}}
						}))
					case v1beta1.Resource:
						b.Watches(v.(v1beta1.Resource), handler.EnqueueRequestsFromMapFunc(func(watchContext context.Context, object client.Object) []reconcile.Request {
							return []reconcile.Request{{
								NamespacedName: types.NamespacedName{
									Name: object.(v1beta1.Resource).GetStack(),
								},
							}}
						}))

					}
				}

				return nil
			}),
			// notes(gfyrag): Some resources need to be properly dropped before the stack is dropped
			WithFinalizer[*v1beta1.Stack]("delete", Clean),
		),
	)
}
