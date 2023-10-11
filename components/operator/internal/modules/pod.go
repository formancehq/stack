package modules

import (
	"context"
	"fmt"
	"strconv"

	"sort"

	"github.com/formancehq/stack/libs/go-libs/pointer"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/formancehq/operator/internal/controllerutils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	monopodLabel  = "formance.com/monopod"
	stackLabel    = "stack"
	productLabel  = "app.kubernetes.io/name"
	stackReplicas = "stack-replicas"
)

type pod struct {
	name                 string
	moduleName           string
	volumes              []corev1.Volume
	initContainers       []corev1.Container
	containers           []corev1.Container
	disableRollingUpdate bool
	mono                 bool
	replicas             *int32
	annotations          map[string]string
}

type PodDeployer interface {
	deploy(ctx context.Context, pod pod) error
	shutdown(ctx context.Context, podName string) (bool, error)
}

type PodDeployerFinalizer interface {
	finalize(ctx context.Context) error
}

type defaultPodDeployer struct {
	deployer *scopedResourceDeployer
}

var _ PodDeployer = (*defaultPodDeployer)(nil)

func (d *defaultPodDeployer) deploy(ctx context.Context, pod pod) error {
	return controllerutils.JustError(d.deployer.
		Deployments().
		CreateOrUpdate(ctx, pod.name, func(t *appsv1.Deployment) {
			matchLabels := map[string]string{
				"app.kubernetes.io/name": pod.name,
			}
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
				stackLabel:   "true",
				productLabel: pod.moduleName,
			}
			replicas := pod.replicas
			if replicas == nil {
				if t.Annotations[stackReplicas] != "" {
					savedReplicas, err := strconv.ParseInt(t.Annotations[stackReplicas], 10, 32)
					if err != nil {
						panic(err)
					}
					replicas = pointer.For(int32(savedReplicas))
					delete(t.Annotations, stackReplicas)
				} else {
					replicas = t.Spec.Replicas
				}
			}

			if t.Annotations != nil {
				// notes(gfyrag): Usage of the stakater reloader operator until v0.17.0
				delete(t.Annotations, "reloader.stakater.com/auto")
			}
			t.Spec = appsv1.DeploymentSpec{
				Replicas: replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: matchLabels,
				},
				Strategy: strategy,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels:      matchLabels,
						Annotations: pod.annotations,
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

func (d *defaultPodDeployer) shutdown(ctx context.Context, podName string) (bool, error) {
	return scaleDownToZero(ctx, d.deployer, podName)
}

func NewDefaultPodDeployer(deployer *scopedResourceDeployer) *defaultPodDeployer {
	return &defaultPodDeployer{
		deployer: deployer,
	}
}

type monoPodDeployer struct {
	deployer *scopedResourceDeployer
	pod
}

var _ PodDeployer = (*monoPodDeployer)(nil)

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

func (d *monoPodDeployer) shutdown(ctx context.Context, _ string) (bool, error) {
	return scaleDownToZero(ctx, d.deployer, d.name)
}

func (d *monoPodDeployer) finalize(ctx context.Context) error {
	sort.SliceStable(d.pod.containers, func(i, j int) bool {
		return d.pod.containers[i].Name < d.pod.containers[j].Name
	})
	sort.SliceStable(d.pod.volumes, func(i, j int) bool {
		return d.pod.volumes[i].Name < d.pod.volumes[j].Name
	})
	return NewDefaultPodDeployer(d.deployer).deploy(ctx, d.pod)
}

func NewMonoPodDeployer(deployer *scopedResourceDeployer, stackName string) *monoPodDeployer {
	return &monoPodDeployer{
		deployer: deployer,
		pod: pod{
			name: stackName,
			mono: true,
		},
	}
}

func scaleDownToZero(ctx context.Context, deployer *scopedResourceDeployer, name string) (bool, error) {
	deployment, err := deployer.Deployments().Get(ctx, name)
	if apierrors.IsNotFound(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	ok, err := ensureDeploymentSync(ctx, *deployment)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	_, err = deployer.Deployments().Update(ctx, deployment, func(t *appsv1.Deployment) {
		t.Annotations[stackReplicas] = fmt.Sprint(*t.Spec.Replicas)
		t.Spec.Replicas = pointer.For(int32(0))
	})
	return false, err
}
