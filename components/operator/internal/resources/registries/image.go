package registries

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
)

func NormalizeVersion(version string) string {
	if version == "" {
		version = "latest"
	}
	return version
}

func GetImage(ctx core.Context, stack *v1beta1.Stack, name, version string) (string, error) {
	return TranslateImage(ctx, stack.Name,
		fmt.Sprintf("ghcr.io/formancehq/%s:%s", name, NormalizeVersion(version)))
}
