package payments

import (
	_ "embed"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:embed Caddyfile.payments.gotpl
var Caddyfile string

func GetIfEnabled(ctx reconcilers.Context, stackName string) (*v1beta1.Payments, error) {
	PaymentsList := &v1beta1.PaymentsList{}
	if err := ctx.GetClient().List(ctx, PaymentsList, client.MatchingFields{
		".spec.stack": stackName,
	}); err != nil {
		return nil, err
	}

	switch len(PaymentsList.Items) {
	case 0:
		return nil, nil
	case 1:
		return &PaymentsList.Items[0], nil
	default:
		return nil, errors.New("found multiple Payments")
	}
}
