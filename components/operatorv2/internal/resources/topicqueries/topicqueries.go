package topicqueries

import (
	"fmt"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"k8s.io/apimachinery/pkg/types"
)

func Create(ctx core.Context, stack *v1beta1.Stack, service, queriedBy string) error {
	_, _, err := core.CreateOrUpdate[*v1beta1.TopicQuery](ctx, types.NamespacedName{
		Name: core.GetObjectName(stack.Name, fmt.Sprintf("%s-%s", queriedBy, service)),
	}, func(t *v1beta1.TopicQuery) {
		t.Spec.QueriedBy = queriedBy
		t.Spec.Stack = stack.Name
		t.Spec.Service = service
	})
	return err
}
