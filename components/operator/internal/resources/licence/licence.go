package licence

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/resourcereferences"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func GetLicenceEnvVars(ctx core.Context, stack *v1beta1.Stack, ownerName string, owner v1beta1.Dependent) ([]v1.EnvVar, error) {
	ret := make([]v1.EnvVar, 0)

	platform := ctx.GetPlatform()

	if platform.LicenceSecret != "" {
		_, err := resourcereferences.Create(ctx, owner, ownerName+"-licence", platform.LicenceSecret, &v1.Secret{})
		if err != nil {
			return nil, err
		}
	} else {
		err := resourcereferences.Delete(ctx, owner, ownerName+"-licence")
		if err != nil {
			return nil, err
		}

		ret = append(ret, core.Env("LICENCE_ENABLED", "false"))
		return ret, nil
	}

	ns := &v1.Namespace{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: "kube-system",
	}, ns); err != nil {
		return nil, err
	}

	ret = append(ret, core.Env("LICENCE_ENABLED", "false"))

	ret = append(ret, core.EnvFromSecret("LICENCE_TOKEN", platform.LicenceSecret, "token"))
	ret = append(ret, core.EnvFromSecret("LICENCE_ISSUER", platform.LicenceSecret, "issuer"))
	ret = append(ret, core.Env("LICENCE_VALIDATE_TICK", "24h"))
	ret = append(ret, core.Env("LICENCE_CLUSTER_ID", string(ns.UID)))

	return ret, nil
}
