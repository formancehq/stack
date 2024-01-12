package deployments

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/stoewer/go-strcase"
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func WithMatchingLabels(name string) func(deployment *v1.Deployment) {
	return func(deployment *v1.Deployment) {
		matchLabels := map[string]string{
			"app.kubernetes.io/name": name,
		}
		if deployment.Spec.Selector == nil {
			deployment.Spec.Selector = &metav1.LabelSelector{}
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

func Create(ctx core.Context, owner interface {
	client.Object
	GetStack() string
	SetCondition(condition v1beta1.Condition)
}, name string, mutators ...core.ObjectMutator[*v1.Deployment]) (*v1.Deployment, error) {

	condition := v1beta1.Condition{
		Type:               "DeploymentReady",
		ObservedGeneration: owner.GetGeneration(),
		LastTransitionTime: metav1.Now(),
		Reason:             strcase.UpperCamelCase(name),
	}
	defer func() {
		owner.SetCondition(condition)
	}()

	mutators = append(mutators, core.WithController[*v1.Deployment](ctx.GetScheme(), owner))

	deployment, _, err := core.CreateOrUpdate[*v1.Deployment](ctx, types.NamespacedName{
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

func checkStatus(deployment *v1.Deployment) (bool, string) {
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
