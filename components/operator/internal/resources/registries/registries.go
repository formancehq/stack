package registries

import (
	"fmt"
	"strings"

	"github.com/formancehq/operator/internal/resources/settings"

	"github.com/formancehq/operator/internal/core"
)

// ghcr.io/<organization>/<repository>:<version>
// public.ecr.aws/<organization>/jeffail/benthos:<version>
// docker.io/<organization|user>/<image>:<version>
func TranslateImage(ctx core.Context, stackName, image string) (string, error) {
	parts := strings.Split(image, ":")
	repository := parts[0]
	version := parts[1]

	organizationImage := strings.SplitN(repository, "/", 2)
	registry := organizationImage[0]
	path := organizationImage[1]

	imageOverride, err := settings.GetStringOrEmpty(ctx, stackName, "registries", registry, "images", path, "rewrite")
	if err != nil {
		return "", err
	}
	if imageOverride == "" {
		imageOverride = path
	}

	registryEndpoint, err := settings.GetStringOrEmpty(ctx, stackName, "registries", registry, "endpoint")
	if err != nil {
		return "", err
	}
	if registryEndpoint == "" {
		registryEndpoint = registry
	}

	return fmt.Sprintf("%s/%s:%s", registryEndpoint, imageOverride, version), nil
}
