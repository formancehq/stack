package registries

import (
	"fmt"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"
	"strings"

	"golang.org/x/mod/semver"
	corev1 "k8s.io/api/core/v1"
)

func NormalizeVersion(version string) string {
	if version == "" {
		version = "latest"
	}
	return version
}

func GetImage(ctx core.Context, stack *v1beta1.Stack, component, version string) (string, error) {
	return TranslateImage(ctx, stack.Name,
		fmt.Sprintf("ghcr.io/formancehq/%s:%s", component, NormalizeVersion(core.GetVersion(stack, version))))
}

func GetPullPolicy(imageName string) corev1.PullPolicy {
	imageVersion := strings.Split(imageName, ":")[1]
	pullPolicy := corev1.PullIfNotPresent
	if !semver.IsValid(imageVersion) {
		pullPolicy = corev1.PullAlways
	}
	return pullPolicy
}
