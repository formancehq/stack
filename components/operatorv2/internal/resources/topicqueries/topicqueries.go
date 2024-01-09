package topicqueries

import (
	"fmt"
	v1beta1 "github.com/formancehq/operator/v2/api/formance.com/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func Create(ctx core.Context, stack *v1beta1.Stack, service string, owner client.Object) error {
	queriedBy := strings.ToLower(owner.GetObjectKind().GroupVersionKind().Kind)
	_, _, err := core.CreateOrUpdate[*v1beta1.TopicQuery](ctx, types.NamespacedName{
		Name: core.GetObjectName(stack.Name, fmt.Sprintf("%s-%s", queriedBy, service)),
	},
		func(t *v1beta1.TopicQuery) {
			t.Spec.QueriedBy = queriedBy
			t.Spec.Stack = stack.Name
			t.Spec.Service = service
		},
		core.WithController[*v1beta1.TopicQuery](ctx.GetScheme(), owner),
	)
	return err
}
