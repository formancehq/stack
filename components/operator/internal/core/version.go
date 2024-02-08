package core

import (
	"golang.org/x/mod/semver"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"k8s.io/apimachinery/pkg/types"
)

func GetModuleVersion(ctx Context, stack *v1beta1.Stack, module v1beta1.Module) (string, error) {
	if module.GetVersion() != "" {
		return module.GetVersion(), nil
	}
	if stack.Spec.Version != "" {
		return stack.Spec.Version, nil
	}
	if stack.Spec.VersionsFromFile != "" {
		versions := &v1beta1.Versions{}
		if err := ctx.GetClient().Get(ctx, types.NamespacedName{
			Name: stack.Spec.VersionsFromFile,
		}, versions); err != nil {
			return "", err
		}

		kinds, _, err := ctx.GetScheme().ObjectKinds(module)
		if err != nil {
			return "", err
		}
		kind := strings.ToLower(kinds[0].Kind)

		version, ok := versions.Spec[kind]
		if ok && version != "" {
			return version, nil
		}
	}

	return "latest", nil
}

func IsGreaterOrEqual(version string, than string) bool {
	if !semver.IsValid(than) {
		return !semver.IsValid(version) // Any semver version is considered lower
	}
	if !semver.IsValid(version) {
		return true
	}
	return semver.Compare(version, than) >= 0
}

func IsLower(version string, than string) bool {
	if !semver.IsValid(than) {
		return semver.IsValid(version) // Any semver version is considered higher
	}
	if !semver.IsValid(version) {
		return false
	}
	return semver.Compare(version, than) < 0
}
