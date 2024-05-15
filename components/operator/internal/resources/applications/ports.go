package applications

import v1 "k8s.io/api/core/v1"

func StandardHTTPPort() v1.ContainerPort {
	return v1.ContainerPort{
		Name:          "http",
		ContainerPort: 8080,
	}
}
