package brokerconsumers

import (
	"fmt"
	. "github.com/formancehq/stack/libs/go-libs/collectionutils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sort"
	"strings"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Create(ctx core.Context, owner interface {
	client.Object
	GetStack() string
	SetCondition(condition v1beta1.Condition)
}, services ...string) (*v1beta1.BrokerConsumer, error) {
	queriedBy := strings.ToLower(owner.GetObjectKind().GroupVersionKind().Kind)

	sort.Strings(services)

	brokerConsumer, _, err := core.CreateOrUpdate[*v1beta1.BrokerConsumer](ctx, types.NamespacedName{
		Name: fmt.Sprintf("%s-%s", owner.GetName(),
			strings.ToLower(owner.GetObjectKind().GroupVersionKind().Kind),
		),
	},
		func(t *v1beta1.BrokerConsumer) error {
			t.Spec.QueriedBy = queriedBy
			t.Spec.Stack = owner.GetStack()
			t.Spec.Services = services

			return nil
		},
		core.WithController[*v1beta1.BrokerConsumer](ctx.GetScheme(), owner),
	)
	if err != nil {
		return nil, err
	}

	return brokerConsumer, nil
}

func CreateOrUpdateOnAllServices(ctx core.Context, consumer interface {
	client.Object
	GetStack() string
	SetCondition(condition v1beta1.Condition)
}) (*v1beta1.BrokerConsumer, error) {
	services, err := core.ListEventPublishers(ctx, consumer.GetStack())
	if err != nil {
		return nil, err
	}

	filteredServices := Filter(services, func(u unstructured.Unstructured) bool {
		return u.GetKind() != consumer.GetObjectKind().GroupVersionKind().Kind
	})

	return Create(ctx, consumer, Map(filteredServices, func(from unstructured.Unstructured) string {
		return strings.ToLower(from.GetKind())
	})...)
}
