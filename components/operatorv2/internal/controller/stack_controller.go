package controller

import (
	"context"
	pkgError "github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"

	"github.com/formancehq/operator/v2/api/v1beta1"
	. "github.com/formancehq/operator/v2/internal/controller/internal"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ApplyToStacksLabel = "formance.com/apply-to-stacks"
	CopiedSecretLabel  = "formance.com/copied-secret"
	AnyValue           = "any"
	TrueValue          = "true"

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
// +kubebuilder:rbac:groups=formance.com,resources=configurations,verbs=get;list;watch
// +kubebuilder:rbac:groups=formance.com,resources=versions,verbs=get;list;watch
// +kubebuilder:rbac:groups=formance.com,resources=opentelemetryconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=opentelemetryconfigurations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=opentelemetryconfigurations/finalizers,verbs=update
// +kubebuilder:rbac:groups=formance.com,resources=databaseconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=databaseconfigurations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=databaseconfigurations/finalizers,verbs=update
// +kubebuilder:rbac:groups=formance.com,resources=brokerconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=brokerconfigurations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=brokerconfigurations/finalizers,verbs=update

// Reconciler reconciles a Stack object
type StackController struct {
	Client client.Client
	Scheme *runtime.Scheme
}

func (r *StackController) Reconcile(ctx context.Context, stack *v1beta1.Stack) error {

	_, _, err := CreateOrUpdate(ctx, r.Client, types.NamespacedName{
		Name: stack.Name,
	}, WithController[*corev1.Namespace](r.Scheme, stack))
	if err != nil {
		return err
	}

	if err := r.handleSecrets(ctx, stack); err != nil {
		return err
	}

	return nil
}

func (r *StackController) handleSecrets(ctx context.Context, stack *v1beta1.Stack) error {
	secrets, err := r.copySecrets(ctx, stack)
	if err != nil {
		return err
	}

	return r.cleanSecrets(ctx, stack, secrets)
}

func (r *StackController) copySecrets(ctx context.Context, stack *v1beta1.Stack) ([]corev1.Secret, error) {

	requirement, err := labels.NewRequirement(ApplyToStacksLabel, selection.In, []string{stack.Name, AnyValue})
	if err != nil {
		return nil, err
	}

	secretsToCopy := &corev1.SecretList{}
	if err := r.Client.List(ctx, secretsToCopy, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*requirement),
	}); err != nil {
		return nil, err
	}

	for _, secret := range secretsToCopy.Items {
		secretName, ok := secret.Annotations[RewrittenSecretName]
		if !ok {
			secretName = secret.Name
		}

		_, _, err = CreateOrUpdate[*corev1.Secret](ctx, r.Client, types.NamespacedName{
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
			WithController[*corev1.Secret](r.Scheme, stack),
		)
	}

	return secretsToCopy.Items, nil
}

func (r *StackController) cleanSecrets(ctx context.Context, stack *v1beta1.Stack, copiedSecrets []corev1.Secret) error {
	requirement, err := labels.NewRequirement(CopiedSecretLabel, selection.Equals, []string{TrueValue})
	if err != nil {
		return err
	}

	existingSecrets := &corev1.SecretList{}
	if err := r.Client.List(ctx, existingSecrets, &client.ListOptions{
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
		if err := r.Client.Delete(ctx, &existingSecret); err != nil {
			return pkgError.Wrap(err, "error deleting old secret")
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StackController) SetupWithManager(mgr ctrl.Manager) (*builder.Builder, error) {
	indexer := mgr.GetFieldIndexer()
	if err := indexer.IndexField(context.Background(), &v1beta1.OpenTelemetryConfiguration{}, ".spec.stack", func(rawObj client.Object) []string {
		return []string{rawObj.(*v1beta1.OpenTelemetryConfiguration).Spec.Stack}
	}); err != nil {
		return nil, err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Stack{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&corev1.Namespace{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})), nil
}

func ForStack(client client.Client, scheme *runtime.Scheme) *StackController {
	return &StackController{
		Client: client,
		Scheme: scheme,
	}
}
