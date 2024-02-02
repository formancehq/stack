package secretreferences

import (
	"fmt"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//+kubebuilder:rbac:groups=formance.com,resources=secretreferences,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=formance.com,resources=secretreferences/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=formance.com,resources=secretreferences/finalizers,verbs=update

func init() {
	Init(
		WithStackDependencyReconciler[*v1beta1.SecretReference](Reconcile,
			WithOwn[*v1beta1.SecretReference](&v1.Secret{}),
			WithWatch[*v1beta1.SecretReference, *v1.Secret](func(ctx Context, object *v1.Secret) []reconcile.Request {
				ret := make([]reconcile.Request, 0)
				if object.GetLabels()[v1beta1.StackLabel] != "any" {
					for _, stack := range strings.Split(object.GetLabels()[v1beta1.StackLabel], ",") {
						ret = append(ret, BuildReconcileRequests(
							ctx,
							ctx.GetClient(),
							ctx.GetScheme(),
							&v1beta1.SecretReference{},
							client.MatchingFields{
								"stack": stack,
							},
						)...)
					}
				} else {
					ret = append(ret, BuildReconcileRequests(
						ctx,
						ctx.GetClient(),
						ctx.GetScheme(),
						&v1beta1.SecretReference{},
					)...)
				}

				return ret
			}),
		),
	)
}

const (
	CopiedSecretLabel = "formance.com/copied-secret"
	AnyValue          = "any"
	TrueValue         = "true"

	RewrittenSecretName               = "formance.com/referenced-by-name"
	OriginalSecretNamespaceAnnotation = "formance.com/original-secret-namespace"
	OriginalSecretNameAnnotation      = "formance.com/original-secret-name"
)

func Reconcile(ctx Context, stack *v1beta1.Stack, req *v1beta1.SecretReference) error {
	secret, err := findMatchingSecret(ctx, stack.Name, req.Spec.SecretName)
	if err != nil {
		return err
	}

	if req.Status.SyncedSecret != "" && req.Spec.SecretName != req.Status.SyncedSecret {
		oldSecret := &v1.Secret{}
		err := ctx.GetClient().Get(ctx, types.NamespacedName{
			Namespace: stack.Name,
			Name:      req.Status.SyncedSecret,
		}, oldSecret)
		if client.IgnoreNotFound(err) != nil {
			return err
		}

		if err == nil { // Can be not found, if the secret has been manually deleted
			patch := client.MergeFrom(oldSecret.DeepCopy())
			if err := controllerutil.RemoveOwnerReference(req, oldSecret, ctx.GetScheme()); err == nil {
				if err := ctx.GetClient().Patch(ctx, oldSecret, patch); err != nil {
					return nil
				}
			}
			if len(oldSecret.OwnerReferences) == 0 {
				if err := ctx.GetClient().Delete(ctx, oldSecret); err != nil {
					return err
				}
			}
		}
	}

	_, _, err = CreateOrUpdate[*v1.Secret](ctx, types.NamespacedName{
		Namespace: stack.Name,
		Name:      req.Spec.SecretName,
	},
		func(t *v1.Secret) error {
			t.Data = secret.Data

			return nil
		},
		// notes(gfyrag): Just set a owner reference is enough
		// we don't want a controller reference as a secret will be handled by potentially
		// multiple SecretReference object
		WithOwner[*v1.Secret](ctx.GetScheme(), req),
	)
	if err != nil {
		return err
	}

	req.Status.Hash = HashFromSecrets(secret)
	req.Status.SyncedSecret = req.Spec.SecretName

	return nil
}

func findMatchingSecret(ctx Context, stack, name string) (*v1.Secret, error) {
	requirement, err := labels.NewRequirement(v1beta1.StackLabel, selection.In, []string{stack, AnyValue})
	if err != nil {
		return nil, err
	}

	availableSecrets := &v1.SecretList{}
	if err := ctx.GetClient().List(ctx, availableSecrets, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*requirement),
	}); err != nil {
		return nil, errors.Wrap(err, "listing secrets")
	}

	foundSecrets := make([]v1.Secret, 0)
	for _, secret := range availableSecrets.Items {
		secretName, ok := secret.Annotations[RewrittenSecretName]
		if !ok {
			secretName = secret.Name
		}

		if secretName != name {
			continue
		}
		foundSecrets = append(foundSecrets, secret)
	}

	if len(foundSecrets) > 1 {
		return nil, fmt.Errorf("found more than one matching secret for '%s': %s", name, collectionutils.Map(foundSecrets, func(from v1.Secret) string {
			return from.Name
		}))
	}
	if len(foundSecrets) == 0 {
		return nil, fmt.Errorf("secret not found: %s", name)
	}

	return &foundSecrets[0], nil
}
