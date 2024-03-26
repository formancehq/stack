package registries

import (
	"fmt"
	"strings"

	"github.com/formancehq/operator/internal/resources/settings"

	"github.com/formancehq/operator/internal/core"
)

func TranslateImage(ctx core.Context, stackName, image string) (string, error) {

	parts := strings.Split(image, ":")
	repository := parts[0]
	repositoryParts := strings.SplitN(repository, "/", 2)
	var (
		registry, path string
	)
	if len(repositoryParts) == 1 {
		registry = "docker.io"
		path = repository
	} else {
		registry = repositoryParts[0]
		path = repositoryParts[1]
	}

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

	return fmt.Sprintf("%s/%s:%s", registryEndpoint, imageOverride, parts[1]), nil
}
