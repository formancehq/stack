package brokertopicconsumers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/stoewer/go-strcase"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Create(ctx core.Context, service string, owner interface {
	client.Object
	GetStack() string
	SetCondition(condition v1beta1.Condition)
}) (*v1beta1.BrokerTopicConsumer, error) {
	queriedBy := strings.ToLower(owner.GetObjectKind().GroupVersionKind().Kind)

	condition := v1beta1.Condition{
		Type:               "BrokerTopicConsumerCreated",
		ObservedGeneration: owner.GetGeneration(),
		LastTransitionTime: metav1.Now(),
		Reason:             strcase.UpperCamelCase(service),
	}
	defer func() {
		owner.SetCondition(condition)
	}()

	brokerTopicConsumer, _, err := core.CreateOrUpdate[*v1beta1.BrokerTopicConsumer](ctx, types.NamespacedName{
		Name: core.GetObjectName(owner.GetStack(), fmt.Sprintf("%s-%s", queriedBy, service)),
	},
		func(t *v1beta1.BrokerTopicConsumer) error {
			t.Spec.QueriedBy = queriedBy
			t.Spec.Stack = owner.GetStack()
			t.Spec.Service = service

			return nil
		},
		core.WithController[*v1beta1.BrokerTopicConsumer](ctx.GetScheme(), owner),
	)
	if err != nil {
		condition.Message = err.Error()
		condition.Status = metav1.ConditionFalse
		return nil, err
	}
	if !brokerTopicConsumer.Status.Ready {
		condition.Message = "broker topic consumer ready creation pending"
		condition.Status = metav1.ConditionFalse
	} else {
		condition.Message = "broker topic consumer is ok"
		condition.Status = metav1.ConditionTrue
	}

	return brokerTopicConsumer, nil
}

type Consumers []*v1beta1.BrokerTopicConsumer

func (c Consumers) Ready() bool {
	for _, consumer := range c {
		if !consumer.Status.Ready {
			return false
		}
	}
	return true
}

func CreateOrUpdateOnAllServices(ctx core.Context, consumer interface {
	client.Object
	GetStack() string
	SetCondition(condition v1beta1.Condition)
}) (Consumers, error) {
	services, err := core.ListEventPublishers(ctx, consumer.GetStack())
	if err != nil {
		return nil, err
	}

	ret := make([]*v1beta1.BrokerTopicConsumer, 0)
	for _, service := range services {
		brokerTopicConsumer, err := Create(ctx, strings.ToLower(service.GetKind()), consumer)
		if err != nil {
			return nil, err
		}
		ret = append(ret, brokerTopicConsumer)
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Name < ret[j].Name
	})

	return ret, nil
}
