package topics

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TopicExists(ctx core.Context, stack *v1beta1.Stack, name string) (bool, error) {
	topicList := &v1beta1.TopicList{}
	if err := ctx.GetClient().List(ctx, topicList, client.MatchingFields{
		".spec.service": name,
		".spec.stack":   stack.Name,
	}); err != nil {
		return false, err
	}
	return len(topicList.Items) > 0, nil
}
