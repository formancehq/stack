package stacks

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

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

func init() {
	Init(
		WithSimpleIndex[*v1beta1.Stack](".spec.versionsFromFile", func(t *v1beta1.Stack) string {
			return t.Spec.VersionsFromFile
		}),
		WithStdReconciler(Reconcile,
			WithOwn[*v1beta1.Stack](&corev1.Secret{}),
			WithOwn[*v1beta1.Stack](&corev1.Namespace{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})),
		),
	)
}
