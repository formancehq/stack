package core

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/iancoleman/strcase"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func LowerCamelCaseName(ctx Context, ob client.Object) string {
	kinds, _, err := ctx.GetScheme().ObjectKinds(ob)
	if err != nil {
		panic(err)
	}
	return strcase.ToLowerCamel(kinds[0].Kind)
}

func InstalledVersionName(ctx Context, module v1beta1.Module, version string) string {
	return fmt.Sprintf("%s-%s-%s", module.GetStack(), LowerCamelCaseName(ctx, module), version)
}

func ValidateInstalledVersion(ctx Context, module v1beta1.Module, version string) error {
	_, _, err := CreateOrUpdate[*v1beta1.VersionsHistory](ctx, types.NamespacedName{
		Name: InstalledVersionName(ctx, module, version),
	}, func(t *v1beta1.VersionsHistory) error {
		t.Spec.Stack = module.GetStack()
		t.Spec.Module = LowerCamelCaseName(ctx, module)
		t.Spec.Version = version

		return nil
	}, WithController[*v1beta1.VersionsHistory](ctx.GetScheme(), module))
	return err
}
