package resourcereferences

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func Create(ctx core.Context, owner v1beta1.Dependent, name, resourceName string, object client.Object) (*v1beta1.ResourceReference, error) {

	gvk, err := apiutil.GVKForObject(object, ctx.GetScheme())
	if err != nil {
		return nil, err
	}

	resourceReferenceName := fmt.Sprintf("%s-%s", owner.GetName(), name)
	resourceReference, _, err := core.CreateOrUpdate[*v1beta1.ResourceReference](ctx, types.NamespacedName{
		Name: resourceReferenceName,
	}, func(t *v1beta1.ResourceReference) error {
		t.Spec.Stack = owner.GetStack()
		t.Spec.Name = resourceName
		t.Spec.GroupVersionKind = &metav1.GroupVersionKind{
			Group:   gvk.Group,
			Version: gvk.Version,
			Kind:    gvk.Kind,
		}

		return nil
	}, core.WithController[*v1beta1.ResourceReference](ctx.GetScheme(), owner))
	if err != nil {
		return nil, err
	}

	if !resourceReference.Status.Ready {
		return nil, core.NewPendingError()
	}

	return resourceReference, nil
}

func Delete(ctx core.Context, owner v1beta1.Dependent, name string) error {
	resourceReferenceName := fmt.Sprintf("%s-%s", owner.GetName(), name)
	reference := &v1beta1.ResourceReference{}
	reference.SetNamespace(owner.GetStack())
	reference.SetName(resourceReferenceName)
	if err := ctx.GetClient().Delete(ctx, reference); client.IgnoreNotFound(err) != nil {
		return err
	}
	return nil
}
