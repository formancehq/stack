package deployments

import (
	"github.com/formancehq/stack/libs/go-libs/pointer"
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func WithMatchingLabels(name string) func(deployment *v1.Deployment) {
	return func(deployment *v1.Deployment) {
		matchLabels := map[string]string{
			"app.kubernetes.io/name": name,
		}
		if deployment.Spec.Selector == nil {
			deployment.Spec.Selector = &v12.LabelSelector{}
		}
		deployment.Spec.Selector.MatchLabels = matchLabels
		deployment.Spec.Template.Labels = matchLabels
	}
}

func WithReplicas(replicas int32) func(t *v1.Deployment) {
	return func(t *v1.Deployment) {
		t.Spec.Replicas = pointer.For(replicas)
	}
}

func WithContainers(containers ...v13.Container) func(r *v1.Deployment) {
	return func(r *v1.Deployment) {
		r.Spec.Template.Spec.Containers = containers
	}
}

func WithInitContainers(containers ...v13.Container) func(r *v1.Deployment) {
	return func(r *v1.Deployment) {
		r.Spec.Template.Spec.InitContainers = containers
	}
}

func WithVolumes(volumes ...v13.Volume) func(t *v1.Deployment) {
	return func(t *v1.Deployment) {
		t.Spec.Template.Spec.Volumes = volumes
	}
}
