package core

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"golang.org/x/mod/semver"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sort"
	"strings"
)

func GetModuleName(module Module) string {
	return strings.ToLower(module.GetObjectKind().GroupVersionKind().Kind)
}

func ValidateInstalledVersion(ctx Context, module Module) error {
	_, _, err := CreateOrUpdate[*v1beta1.VersionsHistory](ctx, types.NamespacedName{
		Name: fmt.Sprintf("%s-%s-%s", module.GetStack(), GetModuleName(module), module.GetVersion()),
	}, func(t *v1beta1.VersionsHistory) {
		t.Spec.Stack = module.GetStack()
		t.Spec.Module = GetModuleName(module)
		t.Spec.Version = module.GetVersion()
	}, WithController[*v1beta1.VersionsHistory](ctx.GetScheme(), module))
	return err
}

func ActualVersion(ctx Context, module Module) (string, error) {
	list := &v1beta1.VersionsHistoryList{}
	if err := ctx.GetClient().List(ctx, list, client.MatchingFields{
		".spec.module": GetModuleName(module),
	}); err != nil {
		return "", err
	}

	items := list.Items
	if len(items) == 0 {
		return "", nil
	}

	sort.Slice(items, func(i, j int) bool {
		return semver.Compare(items[i].Spec.Version, items[j].Spec.Version) > 0
	})

	return items[0].Spec.Version, nil
}
