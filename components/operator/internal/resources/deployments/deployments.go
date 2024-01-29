package deployments

import (
	"fmt"
	"strconv"

	"github.com/formancehq/operator/internal/resources/settings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/stoewer/go-strcase"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func WithMatchingLabels(name string) func(deployment *appsv1.Deployment) error {
	return func(deployment *appsv1.Deployment) error {
		matchLabels := map[string]string{
			"app.kubernetes.io/name": name,
		}
		if deployment.Spec.Selector == nil {
			deployment.Spec.Selector = &metav1.LabelSelector{}
		}
		deployment.Spec.Selector.MatchLabels = matchLabels
		deployment.Spec.Template.Labels = matchLabels

		return nil
	}
}

func WithReplicas(replicas int32) func(t *appsv1.Deployment) error {
	return func(t *appsv1.Deployment) error {
		t.Spec.Replicas = pointer.For(replicas)

		return nil
	}
}

func WithContainers(containers ...corev1.Container) func(r *appsv1.Deployment) error {
	return func(r *appsv1.Deployment) error {
		r.Spec.Template.Spec.Containers = containers

		return nil
	}
}

func WithInitContainers(containers ...corev1.Container) func(r *appsv1.Deployment) error {
	return func(r *appsv1.Deployment) error {
		r.Spec.Template.Spec.InitContainers = containers

		return nil
	}
}

func WithVolumes(volumes ...corev1.Volume) func(t *appsv1.Deployment) error {
	return func(t *appsv1.Deployment) error {
		t.Spec.Template.Spec.Volumes = volumes

		return nil
	}
}

func CreateOrUpdate(ctx core.Context, stack *v1beta1.Stack, owner interface {
	client.Object
	GetStack() string
	SetCondition(condition v1beta1.Condition)
}, name string, mutators ...core.ObjectMutator[*appsv1.Deployment]) (*appsv1.Deployment, error) {

	condition := v1beta1.Condition{
		Type:               "DeploymentReady",
		ObservedGeneration: owner.GetGeneration(),
		LastTransitionTime: metav1.Now(),
		Reason:             strcase.UpperCamelCase(name),
	}
	defer func() {
		owner.SetCondition(condition)
	}()

	mutators = append(mutators, core.WithController[*appsv1.Deployment](ctx.GetScheme(), owner))
	mutators = append(mutators, func(t *appsv1.Deployment) error {
		for ind, container := range t.Spec.Template.Spec.InitContainers {
			resourceRequirements, err := settings.GetResourceRequirements(ctx, owner.GetStack(),
				"deployments", t.Name, "init-containers", container.Name, "resource-requirements")
			if err != nil {
				return err
			}
			container.Resources = mergeResourceRequirements(container.Resources, *resourceRequirements)
			t.Spec.Template.Spec.InitContainers[ind] = container
		}
		for ind, container := range t.Spec.Template.Spec.Containers {
			resourceRequirements, err := settings.GetResourceRequirements(ctx, owner.GetStack(),
				"deployments", t.Name, "containers", container.Name, "resource-requirements")
			if err != nil {
				return err
			}
			container.Resources = mergeResourceRequirements(container.Resources, *resourceRequirements)
			t.Spec.Template.Spec.Containers[ind] = container
		}
		if stack.Spec.Disabled {
			if t.Spec.Replicas != nil {
				// Store the number of replicas to be able to restore it
				// if the stack is re-enabled
				t.Annotations["replicas"] = fmt.Sprint(*t.Spec.Replicas)
			}
			t.Spec.Replicas = pointer.For(int32(0))
		} else {
			// Restore the number of replicas previously stored if the stack was disabled
			if replicasStr := t.Annotations["replicas"]; replicasStr != "" {
				replicas, err := strconv.ParseInt(replicasStr, 10, 32)
				if err != nil {
					return err
				}
				t.Spec.Replicas = pointer.For(int32(replicas))
			}
		}

		return nil
	})

	deployment, _, err := core.CreateOrUpdate[*appsv1.Deployment](ctx, types.NamespacedName{
		Namespace: owner.GetStack(),
		Name:      name,
	}, mutators...)
	if err != nil {
		condition.Message = err.Error()
		condition.Status = metav1.ConditionFalse
		return nil, err
	}

	ready, message := checkStatus(deployment)
	condition.Message = message
	if !ready {
		condition.Status = metav1.ConditionFalse
	} else {
		condition.Status = metav1.ConditionTrue
	}

	return deployment, nil
}

func checkStatus(deployment *appsv1.Deployment) (bool, string) {
	if deployment.Status.ObservedGeneration != deployment.Generation {
		return false, fmt.Sprintf("Generation not matching, generation: %d, observed: %d)",
			deployment.Generation, deployment.Status.ObservedGeneration)
	}
	if deployment.Spec.Replicas != nil && deployment.Status.UpdatedReplicas < *deployment.Spec.Replicas {
		return false, fmt.Sprintf("waiting for deployment %q rollout to finish: %d out of %d new replicas have been updated",
			deployment.Name, deployment.Status.UpdatedReplicas, *deployment.Spec.Replicas)
	}
	if deployment.Status.Replicas > deployment.Status.UpdatedReplicas {
		return false, fmt.Sprintf("waiting for deployment %q rollout to finish: %d old replicas are pending termination",
			deployment.Name, deployment.Status.Replicas-deployment.Status.UpdatedReplicas)
	}
	if deployment.Status.AvailableReplicas < deployment.Status.UpdatedReplicas {
		return false, fmt.Sprintf("waiting for deployment %q rollout to finish: %d of %d updated replicas are available",
			deployment.Name, deployment.Status.AvailableReplicas, deployment.Status.UpdatedReplicas)
	}

	return true, "deployment is ready"
}

func mergeResourceRequirements(dest, src corev1.ResourceRequirements) corev1.ResourceRequirements {
	if dest.Limits == nil {
		dest.Limits = src.Limits
	}
	if dest.Requests == nil {
		dest.Requests = src.Requests
	}
	if dest.Claims == nil {
		dest.Claims = src.Claims
	}
	return dest
}
