package settings

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
)

func FindTemporalConfiguration(ctx core.Context, stack *v1beta1.Stack) (*v1beta1.TemporalConfiguration, error) {

	temporalDSN, err := RequireString(ctx, stack.Name, "temporal.dsn")
	if err != nil {
		return nil, err
	}

	temporalURI, err := url.Parse(temporalDSN)
	if err != nil {
		return nil, err
	}

	if temporalURI.Scheme != "temporal" {
		return nil, fmt.Errorf("invalid temporal uri: %s", temporalDSN)
	}

	if temporalURI.Path == "" {
		return nil, fmt.Errorf("invalid temporal uri: %s", temporalDSN)
	}

	if !strings.HasPrefix(temporalURI.Path, "/") {
		return nil, fmt.Errorf("invalid temporal uri: %s", temporalDSN)
	}

	temporalTLSCrt, err := GetStringOrEmpty(ctx, stack.Name, "temporal.tls.crt")
	if err != nil {
		return nil, err
	}

	temporalTLSKey, err := GetStringOrEmpty(ctx, stack.Name, "temporal.tls.key")
	if err != nil {
		return nil, err
	}

	return &v1beta1.TemporalConfiguration{
		Address:   temporalURI.Host,
		Namespace: temporalURI.Path[1:],
		TLS: v1beta1.TemporalTLSConfig{
			CRT:        temporalTLSCrt,
			Key:        temporalTLSKey,
			SecretName: temporalURI.Query().Get("secret"),
		},
	}, nil
}
