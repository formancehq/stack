package settings

import (
	"fmt"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// +kubebuilder:rbac:groups=formance.com,resources=settings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=formance.com,resources=settings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=formance.com,resources=settings/finalizers,verbs=update

func init() {
	core.Init(
		core.WithSimpleIndex[*v1beta1.Settings]("keylen", func(t *v1beta1.Settings) string {
			return fmt.Sprint(len(strings.Split(t.Spec.Key, ".")))
		}),
	)
}

func New(name, key, value string, stacks ...string) *v1beta1.Settings {
	return &v1beta1.Settings{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		Spec: v1beta1.SettingsSpec{
			Stacks: stacks,
			Key:    key,
			Value:  value,
		},
	}
}

func CreateOrUpdate(ctx core.Context, name, key string, value any, stacks ...string) (*v1beta1.Settings, error) {
	settings, _, err := core.CreateOrUpdate[*v1beta1.Settings](ctx, types.NamespacedName{
		Name: name,
	}, func(t *v1beta1.Settings) error {
		t.Spec.Key = key
		t.Spec.Value = fmt.Sprint(value)
		t.Spec.Stacks = stacks

		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating settings '%s' with key '%s'", name, key)
	}

	return settings, nil
}
