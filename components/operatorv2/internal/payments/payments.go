package payments

import (
	_ "embed"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	"github.com/formancehq/operator/v2/internal/utils"
)

//go:embed Caddyfile.payments.gotpl
var Caddyfile string

func GetIfEnabled(ctx reconcilers.Context, stackName string) (*v1beta1.Payments, error) {
	return utils.GetSingleStackDependencyObject[*v1beta1.PaymentsList, *v1beta1.Payments](ctx, stackName)
}
