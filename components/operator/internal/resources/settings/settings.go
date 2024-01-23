package settings

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"strings"
)

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
			ConfigurationProperties: v1beta1.ConfigurationProperties{
				Stacks: stacks,
			},
			Key:   key,
			Value: value,
		},
	}
}

func CreateOrUpdate(ctx core.Context, name, key string, value any, stacks ...string) (*v1beta1.Settings, error) {
	settings, _, err := core.CreateOrUpdate[*v1beta1.Settings](ctx, types.NamespacedName{
		Name: name,
	}, func(t *v1beta1.Settings) {
		t.Spec.Key = key
		t.Spec.Value = fmt.Sprint(value)
		t.Spec.Stacks = stacks
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating settings '%s' with key '%s'", name, key)
	}

	return settings, nil
}
