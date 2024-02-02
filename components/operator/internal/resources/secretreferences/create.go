package secretreferences

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Sync(ctx core.Context, owner v1beta1.Dependent, name string, u *v1beta1.URI) (*v1beta1.SecretReference, error) {
	secretReferenceName := fmt.Sprintf("%s-%s", owner.GetName(), name)
	if secret := u.Query().Get("secret"); secret != "" {

		secretReference, _, err := core.CreateOrUpdate[*v1beta1.SecretReference](ctx, types.NamespacedName{
			Name: secretReferenceName,
		}, func(t *v1beta1.SecretReference) error {
			t.Spec.Stack = owner.GetStack()
			t.Spec.SecretName = secret

			return nil
		}, core.WithController[*v1beta1.SecretReference](ctx.GetScheme(), owner))
		if err != nil {
			return nil, err
		}

		if !secretReference.Status.Ready {
			return nil, core.NewPendingError()
		}

		return secretReference, nil
	} else {
		reference := &v1beta1.SecretReference{}
		reference.SetNamespace(owner.GetStack())
		reference.SetName(secretReferenceName)
		if err := ctx.GetClient().Delete(ctx, reference); client.IgnoreNotFound(err) != nil {
			return nil, err
		}
		return nil, nil
	}
}

func Annotate[T client.Object](key string, ref *v1beta1.SecretReference) core.ObjectMutator[T] {
	return func(t T) error {
		annotations := t.GetAnnotations()
		if ref == nil {
			if annotations == nil || len(annotations) == 0 {
				return nil
			}
			delete(annotations, key)
			return nil
		} else {
			if annotations == nil {
				annotations = map[string]string{}
			}
			annotations[key] = ref.Status.Hash
		}

		return nil
	}
}
