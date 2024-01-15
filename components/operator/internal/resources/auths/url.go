package auths

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/gateways"
)

func URL(ctx core.Context, stackName string) (string, error) {
	gateway, err := core.GetIfEnabled[*v1beta1.Gateway](ctx, stackName)
	if err != nil {
		return "", err
	}

	if gateway != nil {
		return gateways.URL(gateway) + "/api/auth", nil
	}

	return "http://auth:8080", nil
}
