package gateways

import (
	_ "embed"
	"fmt"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:embed Caddyfile.gotpl
var Caddyfile string

func GetURLSAsEnvVarsIfEnabled(ctx reconcilers.Context, stackName string) ([]v1.EnvVar, error) {
	gateway, err := GetIfEnabled(ctx, stackName)
	if err != nil {
		return nil, err
	}
	if gateway == nil {
		return nil, nil
	}
	ret := []v1.EnvVar{{
		Name:  "STACK_URL",
		Value: "http://gateway:8080",
	}}
	if gateway.Spec.Ingress != nil {
		ret = append(ret, v1.EnvVar{
			Name:  "STACK_PUBLIC_URL",
			Value: fmt.Sprintf("%s://%s", gateway.Spec.Ingress.Scheme, gateway.Spec.Ingress.Host),
		})
	}

	return ret, nil
}

func GetIfEnabled(ctx reconcilers.Context, stackName string) (*v1beta1.Gateway, error) {
	gatewayList := &v1beta1.GatewayList{}
	if err := ctx.GetClient().List(ctx, gatewayList, client.MatchingFields{
		".spec.stack": stackName,
	}); err != nil {
		return nil, err
	}

	switch len(gatewayList.Items) {
	case 0:
		return nil, nil
	case 1:
		return &gatewayList.Items[0], nil
	default:
		return nil, errors.New("found multiple gateway")
	}
}
