package auths

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
)

func GetIfEnabled(ctx core.Context, stackName string) (*v1beta1.Auth, error) {
	ret, err := stacks.GetSingleStackDependencyObject(ctx, stackName, &v1beta1.AuthList{})
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}

	return ret.(*v1beta1.Auth), nil
}

func IsEnabled(ctx core.Context, stackName string) (bool, error) {
	return stacks.HasSingleStackDependencyObject(ctx, stackName, &v1beta1.AuthList{})
}
