package applications

import (
	"fmt"

	"github.com/formancehq/operator/internal/resources/licence"
	"github.com/formancehq/operator/internal/resources/settings"
	v1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/stoewer/go-strcase"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

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

type runAs struct {
	User  *int64 `json:"user"`
	Group *int64 `json:"group"`
}

func configureSecurityContext(container *corev1.Container, runAs *runAs) {
	if container.SecurityContext == nil {
		container.SecurityContext = &corev1.SecurityContext{}
	}
	if container.SecurityContext.Capabilities == nil {
		container.SecurityContext.Capabilities = &corev1.Capabilities{}
	}
	if container.SecurityContext.Capabilities.Drop == nil {
		container.SecurityContext.Capabilities.Drop = []corev1.Capability{"all"}
	}
	if container.SecurityContext.Privileged == nil {
		container.SecurityContext.Privileged = pointer.For(false)
	}
	if container.SecurityContext.ReadOnlyRootFilesystem == nil {
		container.SecurityContext.ReadOnlyRootFilesystem = pointer.For(true)
	}
	if container.SecurityContext.AllowPrivilegeEscalation == nil {
		container.SecurityContext.AllowPrivilegeEscalation = pointer.For(false)
	}
	if container.SecurityContext.RunAsUser == nil && runAs != nil && runAs.User != nil {
		container.SecurityContext.RunAsUser = runAs.User
	}
	if container.SecurityContext.RunAsGroup == nil && runAs != nil && runAs.Group != nil {
		container.SecurityContext.RunAsGroup = runAs.Group
	}
	if container.SecurityContext.RunAsNonRoot == nil {
		container.SecurityContext.RunAsNonRoot = pointer.For(runAs.User != nil || runAs.Group != nil)
	}
}

type Application struct {
	stateful      bool
	isEE          bool
	owner         v1beta1.Dependent
	deploymentTpl *appsv1.Deployment
}

func (a Application) Stateful() Application {
	return a.WithStateful(true)
}

func (a Application) WithStateful(v bool) Application {
	a.stateful = v
	return a
}

func (a Application) IsEE() Application {
	return a.WithEE(true)
}

func (a Application) WithEE(v bool) Application {
	a.isEE = v
	return a
}

func (a Application) Install(ctx core.Context) error {
	deploymentLabels := map[string]string{
		"app.kubernetes.io/name": a.deploymentTpl.Name,
	}

	err := a.handleDeployment(ctx, deploymentLabels)
	if err != nil {
		return err
	}

	return a.handlePDB(ctx, deploymentLabels)
}

func (a Application) handleDeployment(ctx core.Context, deploymentLabels map[string]string) error {
	condition := v1beta1.Condition{
		Type:               "DeploymentReady",
		ObservedGeneration: a.owner.GetGeneration(),
		LastTransitionTime: metav1.Now(),
		Reason:             strcase.UpperCamelCase(a.deploymentTpl.Name),
	}
	defer func() {
		a.owner.GetConditions().AppendOrReplace(condition, v1beta1.AndConditions(
			v1beta1.ConditionTypeMatch("DeploymentReady"),
			v1beta1.ConditionReasonMatch(strcase.UpperCamelCase(a.deploymentTpl.Name)),
		))
	}()

	mutators := make([]core.ObjectMutator[*appsv1.Deployment], 0)
	mutators = append(mutators, func(deployment *appsv1.Deployment) error {
		a.deploymentTpl.Spec.DeepCopyInto(&deployment.Spec)
		deployment.SetName(a.deploymentTpl.Name)
		deployment.SetNamespace(a.owner.GetStack())

		// Configure matching labels
		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: deploymentLabels,
		}
		deployment.Spec.Template.Labels = deploymentLabels

		// Configure security context
		deployment.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{
			RunAsNonRoot: pointer.For(true),
			SeccompProfile: &corev1.SeccompProfile{
				Type: corev1.SeccompProfileTypeRuntimeDefault,
			},
		}

		for ind, container := range deployment.Spec.Template.Spec.InitContainers {
			resourceRequirements, err := settings.GetResourceRequirements(ctx, a.owner.GetStack(),
				"deployments", deployment.Name, "init-containers", container.Name, "resource-requirements")
			if err != nil {
				return err
			}
			container.Resources = mergeResourceRequirements(container.Resources, *resourceRequirements)

			runAs, err := settings.GetAs[runAs](ctx, a.owner.GetStack(),
				"deployments", deployment.Name, "init-containers", container.Name, "run-as")
			if err != nil {
				return err
			}

			configureSecurityContext(&container, runAs)
			deployment.Spec.Template.Spec.InitContainers[ind] = container
		}
		for ind, container := range deployment.Spec.Template.Spec.Containers {
			resourceRequirements, err := settings.GetResourceRequirements(ctx, a.owner.GetStack(),
				"deployments", deployment.Name, "containers", container.Name, "resource-requirements")
			if err != nil {
				return err
			}
			container.Resources = mergeResourceRequirements(container.Resources, *resourceRequirements)

			runAs, err := settings.GetAs[runAs](ctx, a.owner.GetStack(),
				"deployments", deployment.Name, "containers", container.Name, "run-as")
			if err != nil {
				return err
			}

			configureSecurityContext(&container, runAs)
			deployment.Spec.Template.Spec.Containers[ind] = container
		}

		return nil
	})

	if !a.stateful {
		mutators = append(mutators, func(t *appsv1.Deployment) error {
			replicas, err := settings.GetInt32(ctx, a.owner.GetStack(), "deployments", a.deploymentTpl.Name, "replicas")
			if err != nil {
				return err
			}
			t.Spec.Replicas = replicas
			return nil
		})
	}

	if a.isEE {
		licenceSecretResourceRef, licenceEnv, err := licence.GetLicenceEnvVars(ctx, a.deploymentTpl.Name, a.owner)
		if err != nil {
			return err
		}
		if len(licenceEnv) > 0 {
			mutators = append(mutators, func(t *appsv1.Deployment) error {
				for i, container := range t.Spec.Template.Spec.InitContainers {
					container.Env = append(container.Env, licenceEnv...)
					t.Spec.Template.Spec.InitContainers[i] = container
				}
				for i, container := range t.Spec.Template.Spec.Containers {
					container.Env = append(container.Env, licenceEnv...)
					t.Spec.Template.Spec.Containers[i] = container
				}
				if licenceSecretResourceRef != nil {
					if t.Spec.Template.Annotations == nil {
						t.Spec.Template.Annotations = map[string]string{}
					}
					t.Spec.Template.Annotations["licence-secret-hash"] = licenceSecretResourceRef.Status.Hash
				}

				return nil
			})
		}
	}

	mutators = append(mutators, core.WithController[*appsv1.Deployment](ctx.GetScheme(), a.owner))

	isJsonLogging, err := settings.GetBoolOrFalse(ctx, a.owner.GetStack(), "logging", "json")
	if err != nil {
		return err
	}
	if isJsonLogging {
		mutators = append(mutators, func(t *appsv1.Deployment) error {
			v := corev1.EnvVar{
				Name:  "JSON_FORMATTING_LOGGER",
				Value: "true",
			}
			for i, container := range t.Spec.Template.Spec.InitContainers {
				container.Env = append(container.Env, v)
				t.Spec.Template.Spec.InitContainers[i] = container
			}
			for i, container := range t.Spec.Template.Spec.Containers {
				container.Env = append(container.Env, v)
				t.Spec.Template.Spec.Containers[i] = container
			}

			return nil
		})
	}

	deployment, _, err := core.CreateOrUpdate[*appsv1.Deployment](ctx, types.NamespacedName{
		Namespace: a.owner.GetStack(),
		Name:      a.deploymentTpl.Name,
	}, mutators...)
	if err != nil {
		condition.Message = err.Error()
		condition.Status = metav1.ConditionFalse
		return err
	}

	ready, message := checkStatus(deployment)
	condition.Message = message
	if !ready {
		condition.Status = metav1.ConditionFalse
	} else {
		condition.Status = metav1.ConditionTrue
	}

	return nil
}

type podDisruptionBudgetConfiguration struct {
	MinAvailable   string `json:"minAvailable"`
	MaxUnavailable string `json:"maxUnavailable"`
}

func (a Application) handlePDB(ctx core.Context, deploymentLabels map[string]string) error {
	podDisruptionBudgetCondition := v1beta1.NewCondition("PodDisruptionBudget", a.owner.GetGeneration()).
		SetReason(strcase.UpperCamelCase(a.deploymentTpl.Name))
	defer func() {
		a.owner.GetConditions().AppendOrReplace(*podDisruptionBudgetCondition, v1beta1.AndConditions(
			v1beta1.ConditionTypeMatch("PodDisruptionBudget"),
			v1beta1.ConditionReasonMatch(strcase.UpperCamelCase(a.deploymentTpl.Name)),
		))
	}()
	if !a.stateful {

		pdb, err := settings.GetAs[podDisruptionBudgetConfiguration](ctx, a.owner.GetStack(), "deployments", a.deploymentTpl.Name, "pod-disruption-budget")
		if err != nil {
			return err
		}

		if pdb.MinAvailable != "" || pdb.MaxUnavailable != "" {
			podDisruptionBudgetConfiguredCondition := v1beta1.NewCondition("PodDisruptionBudgetConfigured", a.owner.GetGeneration()).
				SetReason(strcase.UpperCamelCase(a.deploymentTpl.Name))

			defer func() {
				a.owner.GetConditions().AppendOrReplace(*podDisruptionBudgetConfiguredCondition, v1beta1.AndConditions(
					v1beta1.ConditionTypeMatch("PodDisruptionBudgetConfigured"),
					v1beta1.ConditionReasonMatch(strcase.UpperCamelCase(a.deploymentTpl.Name)),
				))
			}()

			_, _, err = core.CreateOrUpdate(ctx, types.NamespacedName{
				Namespace: a.owner.GetStack(),
				Name:      a.deploymentTpl.Name,
			}, func(t *v1.PodDisruptionBudget) error {
				if pdb.MinAvailable != "" {
					t.Spec.MinAvailable = pointer.For(intstr.Parse(pdb.MinAvailable))
				}
				if pdb.MaxUnavailable != "" {
					t.Spec.MaxUnavailable = pointer.For(intstr.Parse(pdb.MaxUnavailable))
				}
				t.Spec.Selector = &metav1.LabelSelector{
					MatchLabels: deploymentLabels,
				}
				return nil
			},
				core.WithController[*v1.PodDisruptionBudget](ctx.GetScheme(), a.owner),
			)
			if err != nil {
				podDisruptionBudgetConfiguredCondition.SetStatus(metav1.ConditionFalse).SetMessage(err.Error())
				return err
			}
		} else {
			if err := a.deletePDBIfExists(ctx); err != nil {
				return err
			}
			podDisruptionBudgetCondition.SetMessage("no PDB found")
		}
	} else {
		if err := a.deletePDBIfExists(ctx); err != nil {
			return err
		}

		podDisruptionBudgetCondition.SetMessage("application defined as stateful")
	}

	return nil
}

func (a Application) deletePDBIfExists(ctx core.Context) error {
	pdb := &v1.PodDisruptionBudget{}
	pdb.SetName(a.deploymentTpl.Name)
	pdb.SetNamespace(a.owner.GetStack())

	return client.IgnoreNotFound(ctx.GetClient().Delete(ctx, pdb))
}

func New(owner v1beta1.Dependent, deploymentTpl *appsv1.Deployment) *Application {
	return &Application{
		owner:         owner,
		deploymentTpl: deploymentTpl,
	}
}
