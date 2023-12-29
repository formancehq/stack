package auths

import (
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetIfEnabled(ctx reconcilers.Context, stackName string) (*v1beta1.Auth, error) {
	authList := &v1beta1.AuthList{}
	if err := ctx.GetClient().List(ctx, authList, client.MatchingFields{
		".spec.stack": stackName,
	}); err != nil {
		return nil, err
	}

	switch len(authList.Items) {
	case 0:
		return nil, nil
	case 1:
		return &authList.Items[0], nil
	default:
		return nil, errors.New("found multiple auth")
	}
}

func IsEnabled(ctx reconcilers.Context, stackName string) (bool, error) {
	list := &v1beta1.AuthList{}
	if err := ctx.GetClient().List(ctx, list, client.MatchingFields{
		".spec.stack": stackName,
	}); err != nil {
		return false, err
	}

	return len(list.Items) > 0, nil
}
