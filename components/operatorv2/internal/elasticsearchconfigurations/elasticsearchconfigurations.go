package elasticsearchconfigurations

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/stacks"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Require(ctx core.Context, stackName string) (*v1beta1.ElasticSearchConfiguration, error) {

	elasticSearchConfiguration, err := Get(ctx, stackName)
	if err != nil {
		return nil, err
	}
	if elasticSearchConfiguration == nil {
		return nil, stacks.ErrNotFound
	}

	return elasticSearchConfiguration, nil
}

func Get(ctx core.Context, stackName string) (*v1beta1.ElasticSearchConfiguration, error) {

	stackSelectorRequirement, err := labels.NewRequirement("formance.com/stack", selection.In, []string{"any", stackName})
	if err != nil {
		return nil, err
	}

	elasticSearchConfigurationList := &v1beta1.ElasticSearchConfigurationList{}
	if err := ctx.GetClient().List(ctx, elasticSearchConfigurationList, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*stackSelectorRequirement),
	}); err != nil {
		return nil, err
	}

	switch len(elasticSearchConfigurationList.Items) {
	case 0:
		return nil, nil
	case 1:
		return &elasticSearchConfigurationList.Items[0], nil
	default:
		return nil, errors.New("found multiple elasticsearch config")
	}
}
