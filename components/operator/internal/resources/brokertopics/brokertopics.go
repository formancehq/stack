package brokertopics

import (
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Find(ctx core.Context, stack *v1beta1.Stack, name string) (*v1beta1.BrokerTopic, error) {
	topicList := &v1beta1.BrokerTopicList{}
	if err := ctx.GetClient().List(ctx, topicList, client.MatchingFields{
		".spec.service": name,
		"stack":         stack.Name,
	}); err != nil {
		return nil, err
	}

	if len(topicList.Items) == 0 {
		return nil, nil
	}

	return &topicList.Items[0], nil
}
