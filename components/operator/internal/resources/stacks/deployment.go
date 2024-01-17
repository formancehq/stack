package stacks

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
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

func handleSecrets(ctx core.Context, stack *v1beta1.Stack) error {
	secrets, err := copySecrets(ctx, stack)
	if err != nil {
		return err
	}

	return cleanSecrets(ctx, stack, secrets)
}

func copySecrets(ctx core.Context, stack *v1beta1.Stack) ([]v1.Secret, error) {

	requirement, err := labels.NewRequirement(core.StackLabel, selection.In, []string{stack.Name, AnyValue})
	if err != nil {
		return nil, err
	}

	secretsToCopy := &v1.SecretList{}
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

		_, _, err = core.CreateOrUpdate[*v1.Secret](ctx, types.NamespacedName{
			Namespace: stack.Name,
			Name:      secretName,
		},
			func(t *v1.Secret) {
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
			core.WithController[*v1.Secret](ctx.GetScheme(), stack),
		)
	}

	return secretsToCopy.Items, nil
}

func cleanSecrets(ctx core.Context, stack *v1beta1.Stack, copiedSecrets []v1.Secret) error {
	requirement, err := labels.NewRequirement(CopiedSecretLabel, selection.Equals, []string{TrueValue})
	if err != nil {
		return err
	}

	existingSecrets := &v1.SecretList{}
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
			return errors.Wrap(err, "error deleting old secret")
		}
	}

	return nil
}
