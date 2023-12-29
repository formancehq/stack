package auths

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	"github.com/formancehq/operator/v2/internal/utils"
)

func GetIfEnabled(ctx reconcilers.Context, stackName string) (*v1beta1.Auth, error) {
	return utils.GetSingleStackDependencyObject[*v1beta1.AuthList, *v1beta1.Auth](ctx, stackName)
}

func IsEnabled(ctx reconcilers.Context, stackName string) (bool, error) {
	return utils.HasSingleStackDependencyObject[*v1beta1.AuthList, *v1beta1.Auth](ctx, stackName)
}
