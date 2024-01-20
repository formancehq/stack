package core

import (
	"fmt"
	"sort"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"golang.org/x/mod/semver"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetModuleName(ctx Context, module v1beta1.Module) string {
	kinds, _, err := ctx.GetScheme().ObjectKinds(module)
	if err != nil {
		panic(err)
	}
	return strings.ToLower(kinds[0].Kind)
}

func InstalledVersionName(ctx Context, module v1beta1.Module, version string) string {
	return fmt.Sprintf("%s-%s-%s", module.GetStack(), GetModuleName(ctx, module), version)
}

func ValidateInstalledVersion(ctx Context, module v1beta1.Module, version string) error {
	_, _, err := CreateOrUpdate[*v1beta1.VersionsHistory](ctx, types.NamespacedName{
		Name: InstalledVersionName(ctx, module, version),
	}, func(t *v1beta1.VersionsHistory) {
		t.Spec.Stack = module.GetStack()
		t.Spec.Module = GetModuleName(ctx, module)
		t.Spec.Version = version
	}, WithController[*v1beta1.VersionsHistory](ctx.GetScheme(), module))
	return err
}

func ActualVersion(ctx Context, module v1beta1.Module) (string, error) {
	list := &v1beta1.VersionsHistoryList{}
	if err := ctx.GetClient().List(ctx, list, client.MatchingFields{
		".spec.module": GetModuleName(ctx, module),
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
