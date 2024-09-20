package auths

import (
	"sort"

	. "github.com/formancehq/go-libs/collectionutils"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func createConfiguration(ctx Context, stack *v1beta1.Stack, auth *v1beta1.Auth, items []*v1beta1.AuthClient) (*corev1.ConfigMap, error) {

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	yamlData, err := yaml.Marshal(struct {
		Clients any `yaml:"clients"`
	}{
		Clients: Map(items, func(from *v1beta1.AuthClient) any {
			return from.Spec
		}),
	})
	if err != nil {
		return nil, err
	}

	cm, _, err := CreateOrUpdate[*corev1.ConfigMap](ctx, types.NamespacedName{
		Namespace: stack.Name,
		Name:      "auth-configuration",
	}, func(t *corev1.ConfigMap) error {
		t.Data = map[string]string{
			"config.yaml": string(yamlData),
		}

		return nil
	}, WithController[*corev1.ConfigMap](ctx.GetScheme(), auth))
	if err != nil {
		return nil, err
	}

	return cm, nil
}
