package modules

import (
	"fmt"
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

func GetImage(component, version string) string {
	return fmt.Sprintf("ghcr.io/formancehq/%s:%s", component, NormalizeVersion(version))
}

func GetPullPolicy(imageName string) corev1.PullPolicy {
	imageVersion := strings.Split(imageName, ":")[1]
	pullPolicy := corev1.PullIfNotPresent
	if !semver.IsValid(imageVersion) {
		pullPolicy = corev1.PullAlways
	}
	return pullPolicy
}
