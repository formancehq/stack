package auths

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/stacks"
)

func GetIfEnabled(ctx core.Context, stackName string) (*v1beta1.Auth, error) {
	return stacks.GetSingleStackDependencyObject[*v1beta1.AuthList, *v1beta1.Auth](ctx, stackName)
}

func IsEnabled(ctx core.Context, stackName string) (bool, error) {
	return stacks.HasSingleStackDependencyObject[*v1beta1.AuthList, *v1beta1.Auth](ctx, stackName)
}
