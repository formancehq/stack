package registries

import (
	"fmt"
	"strings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
)

func TranslateImage(ctx core.Context, stackName, image string) (string, error) {
	registries, err := core.GetConfigurationObject[*v1beta1.RegistriesConfiguration](ctx, stackName)
	if err != nil {
		return "", err
	}

	if registries == nil {
		return image, nil
	}

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
	if config, ok := registries.Spec.Registries[registry]; ok && config.Endpoint != "" {
		return fmt.Sprintf("%s/%s:%s", config.Endpoint, path, parts[1]), nil
	}

	return image, nil
}
