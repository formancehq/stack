package stacks

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
// +kubebuilder:rbac:groups=formance.com,resources=settings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=settings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=settings/finalizers,verbs=update

func Reconcile(ctx Context, stack *v1beta1.Stack) error {

	if stack.Spec.Disabled {
		ns := &corev1.Namespace{}
		ns.SetName(stack.Name)
		return client.IgnoreNotFound(ctx.GetClient().Delete(ctx, ns))
	}

	_, _, err := CreateOrUpdate(ctx, types.NamespacedName{
		Name: stack.Name,
	},
		WithController[*corev1.Namespace](ctx.GetScheme(), stack),
	)
	if err != nil {
		return err
	}

	if err := handleSecrets(ctx, stack); err != nil {
		return err
	}

	return nil
}

func init() {
	Init(
		WithSimpleIndex[*v1beta1.Stack](".spec.versionsFromFile", func(t *v1beta1.Stack) string {
			return t.Spec.VersionsFromFile
		}),
		WithStdReconciler(Reconcile,
			WithOwn(&corev1.Secret{}),
			WithOwn(&corev1.Namespace{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})),
			WithWatch[*corev1.Secret](func(ctx Context, object *corev1.Secret) []reconcile.Request {
				if object.GetLabels()[StackLabel] == "any" {
					list := &v1beta1.StackList{}
					if err := ctx.GetClient().List(ctx, list); err != nil {
						return []reconcile.Request{}
					}

					return MapObjectToReconcileRequests(Map(list.Items, ToPointer[v1beta1.Stack])...)
				}
				if object.GetLabels()[StackLabel] != "" {
					return []reconcile.Request{{
						NamespacedName: types.NamespacedName{
							Name: object.GetLabels()[StackLabel],
						},
					}}
				}

				return []reconcile.Request{}
			}),
		),
	)
}
