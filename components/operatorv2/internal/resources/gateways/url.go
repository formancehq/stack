package gateways

import (
	"fmt"
	"github.com/formancehq/operator/v2/api/v1beta1"
)

func URL(gateway *v1beta1.Gateway) string {
	if gateway.Spec.Ingress != nil {
		return fmt.Sprintf("%s://%s", gateway.Spec.Ingress.Scheme, gateway.Spec.Ingress.Host)
	} else {
		return fmt.Sprintf("http://gateway:8080")
	}
}
