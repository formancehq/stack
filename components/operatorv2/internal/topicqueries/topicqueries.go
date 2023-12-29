package topicqueries

import (
	"fmt"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/common"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	"github.com/formancehq/operator/v2/internal/utils"
	"k8s.io/apimachinery/pkg/types"
)

func Create(ctx reconcilers.Context, stack *v1beta1.Stack, service, queriedBy string) error {
	_, _, err := utils.CreateOrUpdate[*v1beta1.TopicQuery](ctx, types.NamespacedName{
		Name: common.GetObjectName(stack.Name, fmt.Sprintf("%s-%s", queriedBy, service)),
	}, func(t *v1beta1.TopicQuery) {
		t.Spec.QueriedBy = queriedBy
		t.Spec.Stack = stack.Name
		t.Spec.Service = service
	})
	return err
}
