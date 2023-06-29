package modules

import (
	"context"

	"github.com/formancehq/operator/internal/collectionutils"
	"github.com/formancehq/operator/internal/controllerutils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	monopodLabel = "formance.com/monopod"
	stackLabel   = "stack"
)

type pod struct {
	name                 string
	volumes              []corev1.Volume
	initContainers       []corev1.Container
	containers           []corev1.Container
	disableRollingUpdate bool
	mono                 bool
}

type PodDeployer interface {
	deploy(ctx context.Context, pod pod) error
}

type defaultPodDeployer struct {
	deployer *ResourceDeployer
}

func (d *defaultPodDeployer) deploy(ctx context.Context, pod pod) error {
	return controllerutils.JustError(d.deployer.
		Deployments().
		CreateOrUpdate(ctx, pod.name, func(t *appsv1.Deployment) {
			matchLabels := collectionutils.CreateMap("app.kubernetes.io/name", pod.name)
			strategy := appsv1.DeploymentStrategy{}
			if pod.disableRollingUpdate {
				strategy.Type = appsv1.RecreateDeploymentStrategyType
			}
			t.Labels = map[string]string{
				monopodLabel: func() string {
					if pod.mono {
						return "true"
					}
					return "false"
				}(),
				stackLabel: "true",
			}
			t.Spec = appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: matchLabels,
				},
				Strategy: strategy,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: matchLabels,
					},
					Spec: corev1.PodSpec{
						Volumes:        pod.volumes,
						InitContainers: pod.initContainers,
						Containers:     pod.containers,
					},
				},
			}
		}),
	)
}

func NewDefaultPodDeployer(deployer *ResourceDeployer) *defaultPodDeployer {
	return &defaultPodDeployer{
		deployer: deployer,
	}
}

type monoPodDeployer struct {
	deployer *ResourceDeployer
	pod
}

func (d *monoPodDeployer) deploy(ctx context.Context, pod pod) error {
	if pod.disableRollingUpdate {
		logf.FromContext(ctx).Info("cannot disable rolling update in monopods, use with caution")
	}
	for _, volume := range pod.volumes {
		volume.Name = pod.name + "-" + volume.Name
		d.volumes = append(d.volumes, volume)
	}
	for _, container := range pod.containers {
		for ind := range container.VolumeMounts {
			container.VolumeMounts[ind].Name = pod.name + "-" + container.VolumeMounts[ind].Name
		}
		d.containers = append(d.containers, container)
	}
	for _, container := range pod.initContainers {
		for ind := range container.VolumeMounts {
			container.VolumeMounts[ind].Name = pod.name + "-" + container.VolumeMounts[ind].Name
		}
		d.initContainers = append(d.initContainers, container)
	}

	return nil
}

func (d *monoPodDeployer) finalize(ctx context.Context) error {
	return NewDefaultPodDeployer(d.deployer).deploy(ctx, d.pod)
}

func NewMonoPodDeployer(deployer *ResourceDeployer, stackName string) *monoPodDeployer {
	return &monoPodDeployer{
		deployer: deployer,
		pod: pod{
			name: stackName,
			mono: true,
		},
	}
}
