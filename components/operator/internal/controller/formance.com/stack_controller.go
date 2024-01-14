package formance_com

import (
	"context"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	pkgError "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	CopiedSecretLabel = "formance.com/copied-secret"
	AnyValue          = "any"
	TrueValue         = "true"

	RewrittenSecretName               = "formance.com/referenced-by-name"
	OriginalSecretNamespaceAnnotation = "formance.com/original-secret-namespace"
	OriginalSecretNameAnnotation      = "formance.com/original-secret-name"
)

// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete;deletecollection
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=stacks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=stacks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=stacks/finalizers,verbs=update
// +kubebuilder:rbac:groups=formance.com,resources=opentelemetryconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=opentelemetryconfigurations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=opentelemetryconfigurations/finalizers,verbs=update
// +kubebuilder:rbac:groups=formance.com,resources=databaseconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=databaseconfigurations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=databaseconfigurations/finalizers,verbs=update
// +kubebuilder:rbac:groups=formance.com,resources=brokerconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=brokerconfigurations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=brokerconfigurations/finalizers,verbs=update
// +kubebuilder:rbac:groups=formance.com,resources=elasticsearchconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=elasticsearchconfigurations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=elasticsearchconfigurations/finalizers,verbs=update
// +kubebuilder:rbac:groups=formance.com,resources=registriesconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=registriesconfigurations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=registriesconfigurations/finalizers,verbs=update
// +kubebuilder:rbac:groups=formance.com,resources=temporalconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=temporalconfigurations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=temporalconfigurations/finalizers,verbs=update

// Reconciler reconciles a Stack object
type StackController struct{}

func (r *StackController) Reconcile(ctx Context, stack *v1beta1.Stack) error {

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

	if err := r.handleSecrets(ctx, stack); err != nil {
		return err
	}

	return nil
}

func (r *StackController) handleSecrets(ctx Context, stack *v1beta1.Stack) error {
	secrets, err := r.copySecrets(ctx, stack)
	if err != nil {
		return err
	}

	return r.cleanSecrets(ctx, stack, secrets)
}

func (r *StackController) copySecrets(ctx Context, stack *v1beta1.Stack) ([]corev1.Secret, error) {

	requirement, err := labels.NewRequirement(StackLabel, selection.In, []string{stack.Name, AnyValue})
	if err != nil {
		return nil, err
	}

	secretsToCopy := &corev1.SecretList{}
	if err := ctx.GetClient().List(ctx, secretsToCopy, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*requirement),
	}); err != nil {
		return nil, err
	}

	for _, secret := range secretsToCopy.Items {
		secretName, ok := secret.Annotations[RewrittenSecretName]
		if !ok {
			secretName = secret.Name
		}

		_, _, err = CreateOrUpdate[*corev1.Secret](ctx, types.NamespacedName{
			Namespace: stack.Name,
			Name:      secretName,
		},
			func(t *corev1.Secret) {
				t.Data = secret.Data
				t.StringData = secret.StringData
				t.Type = secret.Type
				t.Labels = map[string]string{
					CopiedSecretLabel: TrueValue,
				}
				t.Annotations = map[string]string{
					OriginalSecretNamespaceAnnotation: secret.Namespace,
					OriginalSecretNameAnnotation:      secret.Name,
				}
			},
			WithController[*corev1.Secret](ctx.GetScheme(), stack),
		)
	}

	return secretsToCopy.Items, nil
}

func (r *StackController) cleanSecrets(ctx Context, stack *v1beta1.Stack, copiedSecrets []corev1.Secret) error {
	requirement, err := labels.NewRequirement(CopiedSecretLabel, selection.Equals, []string{TrueValue})
	if err != nil {
		return err
	}

	existingSecrets := &corev1.SecretList{}
	if err := ctx.GetClient().List(ctx, existingSecrets, &client.ListOptions{
		Namespace:     stack.Name,
		LabelSelector: labels.NewSelector().Add(*requirement),
	}); err != nil {
		return err
	}

l:
	for _, existingSecret := range existingSecrets.Items {
		originalSecretNamespace := existingSecret.Annotations[OriginalSecretNamespaceAnnotation]
		originalSecretName := existingSecret.Annotations[OriginalSecretNameAnnotation]
		for _, copiedSecret := range copiedSecrets {
			if originalSecretNamespace == copiedSecret.Namespace && originalSecretName == copiedSecret.Name {
				continue l
			}
		}
		if err := ctx.GetClient().Delete(ctx, &existingSecret); err != nil {
			return pkgError.Wrap(err, "error deleting old secret")
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StackController) SetupWithManager(mgr Manager) (*builder.Builder, error) {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Stack{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(&corev1.Secret{}, handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, object client.Object) []reconcile.Request {
			if object.GetLabels()[StackLabel] == "any" {
				list := &v1beta1.StackList{}
				if err := mgr.GetClient().List(ctx, list); err != nil {
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
		})).
		Owns(&corev1.Secret{}).
		Owns(&corev1.Namespace{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForStack() *StackController {
	return &StackController{}
}
