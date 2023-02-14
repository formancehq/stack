package modules

import (
	"fmt"
)

func GetImage(component, version string) string {
	if version == "" {
		version = "latest"
	}
	return fmt.Sprintf("ghcr.io/formancehq/%s:%s", component, version)
}
