package secrets

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"net/url"
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

func Copy(ctx core.Context, owner v1beta1.Dependent, ns, name string) (*v1.Secret, error) {

	requirement, err := labels.NewRequirement(core.StackLabel, selection.In, []string{owner.GetStack(), AnyValue})
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

	_, _, err = core.CreateOrUpdate[*v1.Secret](ctx, types.NamespacedName{
		Namespace: ns,
		Name:      fmt.Sprintf("%s-%s", owner.GetName(), name),
	},
		func(t *v1.Secret) error {
			t.Data = foundSecrets[0].Data
			t.StringData = foundSecrets[0].StringData
			t.Type = foundSecrets[0].Type
			t.Labels = map[string]string{
				CopiedSecretLabel: TrueValue,
			}
			t.Annotations = map[string]string{
				OriginalSecretNamespaceAnnotation: foundSecrets[0].Namespace,
				OriginalSecretNameAnnotation:      foundSecrets[0].Name,
			}

			return nil
		},
		core.WithController[*v1.Secret](ctx.GetScheme(), owner),
	)

	return &foundSecrets[0], nil
}

func Clean(ctx core.Context, owner v1beta1.Dependent, ns string, ignoreSecrets ...*v1.Secret) error {
	requirement, err := labels.NewRequirement(CopiedSecretLabel, selection.Equals, []string{TrueValue})
	if err != nil {
		return err
	}

	existingSecrets := &v1.SecretList{}
	if err := ctx.GetClient().List(ctx, existingSecrets, &client.ListOptions{
		Namespace:     ns,
		LabelSelector: labels.NewSelector().Add(*requirement),
	}); err != nil {
		return errors.Wrap(err, "listing secrets")
	}

l:
	for _, existingSecret := range existingSecrets.Items {
		for _, ignoreSecret := range ignoreSecrets {
			if ignoreSecret.Name == existingSecret.Name {
				continue l
			}
			hasControllerReference, err := core.HasControllerReference(ctx, owner, &existingSecret)
			if err != nil {
				return errors.Wrap(err, "checking controller reference")
			}
			if !hasControllerReference {
				continue l
			}
		}
		if err := ctx.GetClient().Delete(ctx, &existingSecret); err != nil {
			return errors.Wrap(err, "error deleting old secret")
		}
	}

	return nil
}

func Sync(ctx core.Context, owner v1beta1.Dependent, ns string, names ...string) ([]*v1.Secret, error) {
	secrets := make([]*v1.Secret, 0)
	for _, name := range names {
		secret, err := Copy(ctx, owner, ns, name)
		if err != nil {
			return nil, errors.Wrapf(err, "copying secret '%s'", name)
		}
		secrets = append(secrets, secret)
	}

	if err := Clean(ctx, owner, ns, secrets...); err != nil {
		return nil, errors.Wrap(err, "cleaning old secret")
	}

	return secrets, nil
}

func SyncOne(ctx core.Context, owner v1beta1.Dependent, ns string, name string) (*v1.Secret, error) {
	secrets, err := Sync(ctx, owner, ns, name)
	if err != nil {
		return nil, err
	}

	return secrets[0], nil
}

func SyncFromURLs(ctx core.Context, owner v1beta1.Dependent, from, to *url.URL) error {
	existingSecret := from.Query().Get("secret")
	newSecret := to.Query().Get("secret")

	if existingSecret != "" && newSecret != existingSecret {
		secret := &v1.Secret{}
		secret.SetName(fmt.Sprintf("%s-%s", owner.GetName(), existingSecret))
		secret.SetNamespace(owner.GetStack())
		if err := ctx.GetClient().Delete(ctx, secret); err != nil {
			return errors.Wrapf(err, "deleting secret '%s'", secret.Name)
		}
	}

	if newSecret != "" {
		_, err := Copy(ctx, owner, owner.GetStack(), newSecret)
		if err != nil {
			return err
		}
	}

	return nil
}
