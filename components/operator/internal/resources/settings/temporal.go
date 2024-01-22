package settings

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
)

func FindTemporalConfiguration(ctx core.Context, stack *v1beta1.Stack) (*v1beta1.TemporalConfiguration, error) {
	temporalAddress, err := RequireString(ctx, stack.Name, "temporal.address")
	if err != nil {
		return nil, err
	}
	temporalNamespace, err := RequireString(ctx, stack.Name, "temporal.namespace")
	if err != nil {
		return nil, err
	}
	temporalTLSCrt, err := GetStringOrEmpty(ctx, stack.Name, "temporal.tls.crt")
	if err != nil {
		return nil, err
	}
	temporalTLSKey, err := GetStringOrEmpty(ctx, stack.Name, "temporal.tls.key")
	if err != nil {
		return nil, err
	}
	temporalTLSSecret, err := GetStringOrEmpty(ctx, stack.Name, "temporal.tls.secret")
	if err != nil {
		return nil, err
	}

	return &v1beta1.TemporalConfiguration{
		Address:   temporalAddress,
		Namespace: temporalNamespace,
		TLS: v1beta1.TemporalTLSConfig{
			CRT:        temporalTLSCrt,
			Key:        temporalTLSKey,
			SecretName: temporalTLSSecret,
		},
	}, nil
}
