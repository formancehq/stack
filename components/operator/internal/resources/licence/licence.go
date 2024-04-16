package licence

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	v1 "k8s.io/api/core/v1"
)

func GetLicenceEnvVars(ctx core.Context, stack *v1beta1.Stack) ([]v1.EnvVar, error) {
	ret := make([]v1.EnvVar, 0)

	ret = append(ret, core.Env("LICENCE_ENABLED", "true"))

	platform := ctx.GetPlatform()
	ret = append(ret, core.Env("LICENCE_TOKEN", platform.Licence.Token))
	ret = append(ret, core.Env("LICENCE_ISSUER", platform.Licence.Issuer))
	ret = append(ret, core.Env("LICENCE_CLUSTER_ID", platform.Licence.ClusterID))
	ret = append(ret, core.Env("LICENCE_VALIDATE_TICK", "24h"))

	return ret, nil
}
