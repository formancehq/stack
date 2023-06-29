package modules

import (
	"fmt"
)

func NormalizeVersion(version string) string {
	if version == "" {
		version = "latest"
	}
	return version
}

func GetImage(component, version string) string {
	return fmt.Sprintf("ghcr.io/formancehq/%s:%s", component, NormalizeVersion(version))
}
