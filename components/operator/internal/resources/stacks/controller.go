package stacks

import (
	"context"
	"reflect"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
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

func Reconcile(ctx Context, stack *v1beta1.Stack) error {
	_, _, err := CreateOrUpdate(ctx, types.NamespacedName{
		Name: stack.Name,
	},
		WithController[*corev1.Namespace](ctx.GetScheme(), stack),
	)
	if err != nil {
		return err
	}

	return nil
}

func Clean(ctx Context, t *v1beta1.Stack) error {
	stillModules := false
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
			"stack": t.Name,
		}); err != nil {
			return err
		}

		stillModules = stillModules || len(l.Items) > 0
		for _, item := range l.Items {
			if err := ctx.GetClient().Delete(ctx, &item); err != nil {
				return err
			}
		}
	}

	if stillModules {
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
					module, ok := v.(v1beta1.Module)
					if !ok {
						continue
					}

					b.Watches(module, handler.EnqueueRequestsFromMapFunc(func(watchContext context.Context, object client.Object) []reconcile.Request {
						return []reconcile.Request{{
							NamespacedName: types.NamespacedName{
								Name: object.(v1beta1.Module).GetStack(),
							},
						}}
					}))
				}

				return nil
			}),
			// notes(gfyrag): Some resources need to be properly dropped before the stack is dropped
			WithFinalizer[*v1beta1.Stack]("delete", Clean),
		),
	)
}
