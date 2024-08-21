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
// version: "v2.0.0-rc.35-scratch@sha256:4a29620448a90f3ae50d2e375c993b86ef141ead4b6ac1edd1674e9ff6b933f8"
// docker.io/<organization|user>/<image>:v2.0.0-rc.35-scratch@sha256:4a29620448a90f3ae50d2e375c993b86ef141ead4b6ac1edd1674e9ff6b933f8
type imageOrigin struct {
	Registry string
	Image    string
	Version  string
}

func (o imageOrigin) String() string {
	return fmt.Sprintf("%s/%s:%s", o.Registry, o.Image, o.Version)
}

//go:generate mockgen -source ./registries.go -destination ./registries_generated.go -package registries . ImageSettingsOverrider
type ImageSettingsOverrider interface {
	OverrideWithSetting(*imageOrigin, string) error
}

type defaultImageSettingOverrider struct {
	ctx core.Context
}

func NewImageSettingsOverrider(ctx core.Context) ImageSettingsOverrider {
	return &defaultImageSettingOverrider{ctx: ctx}
}

func (is *defaultImageSettingOverrider) OverrideWithSetting(o *imageOrigin, stackName string) (err error) {
	imageOverride, err := settings.GetStringOrEmpty(is.ctx, stackName, "registries", o.Registry, "images", o.Image, "rewrite")
	if err != nil {
		return err
	}
	if imageOverride != "" {
		o.Image = imageOverride
	}

	registryEndpoint, err := settings.GetStringOrEmpty(is.ctx, stackName, "registries", o.Registry, "endpoint")
	if err != nil {
		return err
	}
	if registryEndpoint != "" {
		o.Registry = registryEndpoint
	}

	return nil
}

func TranslateImage(
	stackName string,
	settingsOverrider ImageSettingsOverrider,
	image string,
) (string, error) {
	repository, version, found := strings.Cut(image, ":")
	if !found {
		return "", fmt.Errorf("invalid image format: %s", image)
	}

	organizationImage := strings.SplitN(repository, "/", 2)
	origin := &imageOrigin{
		Registry: organizationImage[0],
		Image:    organizationImage[1],
		Version:  version,
	}

	if err := settingsOverrider.OverrideWithSetting(origin, stackName); err != nil {
		return "", err
	}

	return origin.String(), nil
}
