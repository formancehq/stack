package gateways

import (
	_ "embed"
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	v1 "k8s.io/api/core/v1"
)

//go:embed Caddyfile.gotpl
var Caddyfile string

func EnvVarsIfEnabled(ctx core.Context, stackName string) ([]v1.EnvVar, error) {
	return EnvVarsIfEnabledWithPrefix(ctx, stackName, "")
}

func EnvVarsIfEnabledWithPrefix(ctx core.Context, stackName, prefix string) ([]v1.EnvVar, error) {
	gateway := &v1beta1.Gateway{}
	ok, err := core.GetIfExists(ctx, stackName, gateway)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return GetEnvVarsWithPrefix(gateway, prefix), nil
}

func GetEnvVars(gateway *v1beta1.Gateway) []v1.EnvVar {
	return GetEnvVarsWithPrefix(gateway, "")
}

func GetEnvVarsWithPrefix(gateway *v1beta1.Gateway, prefix string) []v1.EnvVar {
	ret := []v1.EnvVar{{
		Name:  fmt.Sprintf("%sSTACK_URL", prefix),
		Value: "http://gateway:8080",
	}}
	if gateway.Spec.Ingress != nil {
		ret = append(ret, v1.EnvVar{
			Name:  fmt.Sprintf("%sSTACK_PUBLIC_URL", prefix),
			Value: fmt.Sprintf("%s://%s", gateway.Spec.Ingress.Scheme, gateway.Spec.Ingress.Host),
		})
	}

	return ret
}
