package payments

import (
	_ "embed"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"github.com/formancehq/operator/v2/internal/stacks"
)

//go:embed Caddyfile.payments.gotpl
var Caddyfile string

func GetIfEnabled(ctx core.Context, stackName string) (*v1beta1.Payments, error) {
	return stacks.GetSingleStackDependencyObject[*v1beta1.PaymentsList, *v1beta1.Payments](ctx, stackName)
}
