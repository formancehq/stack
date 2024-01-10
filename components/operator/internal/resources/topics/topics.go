package topics

import (
	"fmt"
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Find(ctx core.Context, stack *v1beta1.Stack, name string) (*v1beta1.Topic, error) {
	topicList := &v1beta1.TopicList{}
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

func CreateOrUpdate(ctx core.Context, stack *v1beta1.Stack, owner client.Object, service, name string) (*v1beta1.Topic, error) {
	ret, _, err := core.CreateOrUpdate[*v1beta1.Topic](ctx, types.NamespacedName{
		Name: fmt.Sprintf("%s-%s", stack.Name, name),
	},
		func(t *v1beta1.Topic) {
			t.Spec.Service = service
			t.Spec.Stack = stack.GetName()
		},
		core.WithController[*v1beta1.Topic](ctx.GetScheme(), owner),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
