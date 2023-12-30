package payments

import (
	_ "embed"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/resources/stacks"
)

//go:embed Caddyfile.payments.gotpl
var Caddyfile string

func GetIfEnabled(ctx core.Context, stackName string) (*v1beta1.Payments, error) {
	ret, err := stacks.GetSingleStackDependencyObject(ctx, stackName, &v1beta1.PaymentsList{})
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}

	return ret.(*v1beta1.Payments), nil
}
