package httpapis

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Create(ctx core.Context, stack *v1beta1.Stack, owner client.Object, objectName string, options ...func(spec *v1beta1.HTTPAPISpec)) error {
	_, _, err := core.CreateOrUpdate[*v1beta1.HTTPAPI](ctx, types.NamespacedName{
		Name: core.GetObjectName(stack.Name, objectName),
	},
		func(t *v1beta1.HTTPAPI) {
			t.Spec = v1beta1.HTTPAPISpec{
				StackDependency: v1beta1.StackDependency{
					Stack: stack.Name,
				},
				Name: objectName,
			}
			for _, option := range options {
				option(&t.Spec)
			}
		},
		core.WithController[*v1beta1.HTTPAPI](ctx.GetScheme(), owner),
	)
	return err
}

func Secured() func(httpapi *v1beta1.HTTPAPISpec) {
	return func(httpapiSpec *v1beta1.HTTPAPISpec) {
		httpapiSpec.Secured = true
	}
}
