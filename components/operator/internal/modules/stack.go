package modules

import (
	"context"
	"fmt"
	"strings"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	PartOfConfigurationLabel    = "stacks.formance.cloud/partof-configuration"
	copiedSecretLabel           = "stack.formance.cloud/copied-secret"
	PartOfConfigurationAnyValue = "any"
	trueValue                   = "true"

	secretNameOnConfigurationAnnotation = "stacks.formance.cloud/referenced-by-name"
	originalSecretNamespaceAnnotation   = "stack.formance.cloud/original-secret-namespace"
	originalSecretNameAnnotation        = "stack.formance.cloud/original-secret-name"
)

type StackReconcilerFactory struct {
	platform Platform
	manager  manager.Manager
}

func (sm *StackReconcilerFactory) Platform() Platform {
	return sm.platform
}

func (sm *StackReconcilerFactory) NewDeployer(stack *v1beta3.Stack, configuration *v1beta3.Configuration, versions *v1beta3.Versions) *StackReconciler {
	return newStackReconciler(sm.manager, ReconciliationConfig{
		Stack:         stack,
		Configuration: configuration,
		Versions:      versions,
		Platform:      sm.platform,
	})
}

func NewsStackReconcilerFactory(mgr manager.Manager, platform Platform) *StackReconcilerFactory {
	return &StackReconcilerFactory{
		platform: platform,
		manager:  mgr,
	}
}

type ReconciliationConfig struct {
	Stack         *v1beta3.Stack
	Configuration *v1beta3.Configuration
	Versions      *v1beta3.Versions
	Platform      Platform
}

type StackReconciler struct {
	ReconciliationConfig
	JobRunner
	podDeployer                PodDeployer
	portAllocator              PortAllocator
	namespacedResourceDeployer *scopedResourceDeployer
	manager                    manager.Manager

	ready                collectionutils.Set[Module]
	secretsEventRecorder record.EventRecorder
}

func newStackReconciler(mgr manager.Manager, cfg ReconciliationConfig) *StackReconciler {

	resourceDeployer := NewScopedDeployer(mgr.GetClient(), mgr.GetScheme(), cfg.Stack, cfg.Stack)

	var (
		portAllocator PortAllocator = StaticPortAllocator(8080)
		podDeployer   PodDeployer   = NewDefaultPodDeployer(resourceDeployer)
	)

	if cfg.Configuration.Spec.LightMode {
		podDeployer = NewMonoPodDeployer(resourceDeployer, cfg.Stack.Name)
		portAllocator = NewPortRangeAllocator(10000)
	}
	return &StackReconciler{
		ReconciliationConfig:       cfg,
		namespacedResourceDeployer: resourceDeployer,
		podDeployer:                podDeployer,
		portAllocator:              portAllocator,
		ready:                      collectionutils.NewSet[Module](),
		JobRunner:                  NewJobRunner(mgr.GetClient(), mgr.GetScheme(), cfg.Stack, cfg.Stack, ""),
		manager:                    mgr,
		secretsEventRecorder:       mgr.GetEventRecorderFor("operator"),
	}
}

func (r *StackReconciler) Reconcile(ctx context.Context) (bool, error) {

	logger := log.FromContext(ctx)
	logger = logger.WithValues("stack", r.Stack.Name)

	// When Stack is Disabled, we want to remove all deployments
	if r.Stack.Spec.Disabled {
		logger.Info("Stack is disabled, remove all deployments")
		return falseIfError(r.deleteAllStackDeployments(ctx))
	}

	if r.Stack.ModeChanged(r.Configuration) {
		logger.Info("Stack mode has changed, remove all deployments")
		if err := r.deleteAllStackDeployments(ctx); err != nil {
			return false, err
		}
		r.Stack.Status.LightMode = r.Configuration.Spec.LightMode
	}

	if err := r.prepareSecrets(ctx); err != nil {
		return false, err
	}

	registeredModules := RegisteredModules{}

	logger.Info("Prepare modules")
	err := r.prepareModules(ctx, registeredModules)
	if err != nil {
		return false, err
	}

	if r.Configuration.Spec.LightMode {
		for _, module := range modules {
			if r.Stack.IsDisabled(module.Name()) {
				continue
			}
			if !r.ready.Contains(module) {
				logger.Info("Stop reconciliation because we're in light mode and module is not ready", "module", module.Name())
				return false, nil
			}
		}

		if err := r.podDeployer.(PodDeployerFinalizer).finalize(ctx); err != nil {
			return false, err
		}
	}

	ready, err := r.checkDeployments(ctx)
	if err != nil {
		return false, err
	}
	if !ready {
		logger.Info("Skip modules finalizing as all modules are not ready")
		return false, nil
	}

	logger.Info("Finalize modules")
	allReady := true
	for _, module := range modules {
		if r.Stack.IsDisabled(module.Name()) {
			continue
		}
		if !r.ready.Contains(module) {
			allReady = false
			logger.Info(fmt.Sprintf("Skip post install of modules '%s' as it is marked as not ready", module.Name()))
			continue
		}
		logger.Info(fmt.Sprintf("Post install module '%s'", module.Name()))
		completed, err := newModuleReconciler(r, module).finalizeModule(ctx, module)
		if err != nil {
			return false, err
		}
		allReady = allReady && completed
		if !completed {
			logger.Info(fmt.Sprintf("Module '%s' marked as not completed", module.Name()))
		}
	}

	return allReady, nil
}

func (r *StackReconciler) copySecrets(ctx context.Context) ([]corev1.Secret, error) {

	logger := log.FromContext(ctx)
	logger.Info("Copy secrets")

	requirement, err := labels.NewRequirement(PartOfConfigurationLabel, selection.In, []string{r.Configuration.Name, PartOfConfigurationAnyValue})
	if err != nil {
		return nil, err
	}

	secretsToCopy := &corev1.SecretList{}
	if err := r.manager.GetClient().List(ctx, secretsToCopy, &client.ListOptions{
		LabelSelector: labels.NewSelector().Add(*requirement),
	}); err != nil {
		return nil, err
	}

	for _, secret := range secretsToCopy.Items {
		secretName, ok := secret.Annotations[secretNameOnConfigurationAnnotation]
		if !ok {
			logger.Info("Secret name annotation not found, use secret name", "secret", secret.Name)
			secretName = secret.Name
		}

		_, operationResult, err := r.namespacedResourceDeployer.Secrets().CreateOrUpdate(ctx, secretName, func(t *corev1.Secret) {
			t.Data = secret.Data
			t.StringData = secret.StringData
			t.Type = secret.Type
			t.Labels = map[string]string{
				copiedSecretLabel: trueValue,
			}
			t.Annotations = map[string]string{
				originalSecretNamespaceAnnotation: secret.Namespace,
				originalSecretNameAnnotation:      secret.Name,
			}
		})
		if err != nil {
			return nil, err
		}

		switch operationResult {
		case controllerutil.OperationResultCreated:
			r.secretsEventRecorder.Eventf(r.Stack, "Normal", "Created secret",
				"Secret created from secret %s/%s with name '%s'", secret.Namespace, secret.Name, secretName)
		case controllerutil.OperationResultUpdated:
			r.secretsEventRecorder.Eventf(r.Stack, "Normal", "Updated secret",
				"Secret updated from secret %s/%s with name '%s'", secret.Namespace, secret.Name, secretName)
		}
	}

	return secretsToCopy.Items, nil
}

func (r *StackReconciler) cleanSecrets(ctx context.Context, copiedSecrets []corev1.Secret) error {
	logger := logging.FromContext(ctx)
	logger.Info("Clean secrets")

	requirement, err := labels.NewRequirement(copiedSecretLabel, selection.Equals, []string{trueValue})
	if err != nil {
		return err
	}

	existingSecrets := &corev1.SecretList{}
	if err := r.manager.GetClient().List(ctx, existingSecrets, &client.ListOptions{
		Namespace:     r.Stack.Name,
		LabelSelector: labels.NewSelector().Add(*requirement),
	}); err != nil {
		return err
	}

l:
	for _, existingSecret := range existingSecrets.Items {
		originalSecretNamespace := existingSecret.Annotations[originalSecretNamespaceAnnotation]
		originalSecretName := existingSecret.Annotations[originalSecretNameAnnotation]
		for _, copiedSecret := range copiedSecrets {
			if originalSecretNamespace == copiedSecret.Namespace && originalSecretName == copiedSecret.Name {
				continue l
			}
		}
		if err := r.manager.GetClient().Delete(ctx, &existingSecret); err != nil {
			return errors.Wrap(err, "error deleting old secret")
		}
		r.secretsEventRecorder.AnnotatedEventf(r.Stack, existingSecret.Annotations, "Normal",
			"Removed secret", "Secret %s removed", existingSecret.Name)
	}

	return nil
}

func (r *StackReconciler) prepareSecrets(ctx context.Context) error {
	secrets, err := r.copySecrets(ctx)
	if err != nil {
		return err
	}

	return r.cleanSecrets(ctx, secrets)
}

func (r *StackReconciler) prepareModules(ctx context.Context, registeredModules RegisteredModules) error {
	processed := collectionutils.NewSet[Module]()
	for _, module := range modules {
		if err := r.prepareModule(ctx, module, registeredModules, newGraphVisitor(), processed); err != nil {
			return err
		}
	}
	return nil
}

func (r *StackReconciler) prepareModule(ctx context.Context, module Module,
	registeredModules RegisteredModules, visitor *graphVisitor, processed collectionutils.Set[Module]) error {
	if err := visitor.visit(module); err != nil {
		return err
	}
	defer processed.Put(module)

	dependsOnAwareModule, ok := module.(DependsOnAwareModule)
	if ok {
		for _, dependency := range dependsOnAwareModule.DependsOn() {
			if err := r.prepareModule(ctx, dependency, registeredModules, visitor.copy(), processed); err != nil {
				return err
			}
		}
	}

	if processed.Contains(module) {
		return nil
	}

	if r.Stack.IsDisabled(module.Name()) {
		return r.scaleDownStackModule(ctx, module)
	}

	isReady, deploymentNames, err := newModuleReconciler(r, module).installModule(ctx, registeredModules)
	if err != nil {
		return err
	}
	if isReady {
		log.FromContext(ctx).Info("Mark module as ready", "module", module.Name())
		r.ready.Put(module)

		if err := r.cleanOldModuleDeployments(ctx, module, deploymentNames); err != nil {
			return err
		}
	}

	return nil
}

func (r *StackReconciler) deleteStackDeployment(ctx context.Context, deployment *appsv1.Deployment) error {
	return r.namespacedResourceDeployer.client.Delete(ctx, deployment)
}

func (r *StackReconciler) deleteAllStackDeployments(ctx context.Context) error {
	return r.namespacedResourceDeployer.client.DeleteAllOf(ctx, &appsv1.Deployment{},
		client.InNamespace(r.Stack.Name),
		client.MatchingLabels{
			"stack": "true",
		},
	)
}

func (r *StackReconciler) cleanOldModuleDeployments(ctx context.Context, module Module, deploymentNames map[string]struct{}) error {
	// Clean old deployments
	deploymentsList := appsv1.DeploymentList{}
	if err := r.namespacedResourceDeployer.client.List(ctx, &deploymentsList,
		client.InNamespace(r.Stack.Name),
		client.MatchingLabels(map[string]string{
			"app.kubernetes.io/name": module.Name(),
			stackLabel:               "true",
		})); err != nil {
		return err
	}

	for _, d := range deploymentsList.Items {
		if _, ok := deploymentNames[d.Name]; !ok {
			if err := r.deleteStackDeployment(ctx, &d); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *StackReconciler) checkDeployments(ctx context.Context) (bool, error) {
	deploymentsList := appsv1.DeploymentList{}
	if err := r.namespacedResourceDeployer.client.List(ctx, &deploymentsList,
		client.InNamespace(r.Stack.Name),
		client.MatchingLabels(map[string]string{
			stackLabel: "true",
		})); err != nil {
		return false, err
	}

	for _, deployment := range deploymentsList.Items {
		ok, err := ensureDeploymentSync(ctx, deployment)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	return true, nil
}

func (r *StackReconciler) scaleDownStackModule(ctx context.Context, module Module) error {
	return r.namespacedResourceDeployer.client.DeleteAllOf(ctx, &appsv1.Deployment{},
		client.InNamespace(r.Stack.Name),
		client.MatchingLabels{
			"app.kubernetes.io/name": module.Name(),
			"stack":                  "true",
		},
	)
}

type graphVisitor struct {
	visited collectionutils.Set[Module]
	stack   []Module
}

func (v *graphVisitor) visit(t Module) error {
	v.stack = append(v.stack, t)
	if v.visited.Contains(t) {
		return fmt.Errorf("circular dependencies detected: %s", strings.Join(collectionutils.Map(v.stack, Module.Name), " -> "))
	}
	v.visited.Put(t)

	return nil
}

func (v *graphVisitor) copy() *graphVisitor {
	return &graphVisitor{
		visited: collectionutils.CopyMap(v.visited),
		stack:   v.stack[:],
	}
}

func newGraphVisitor() *graphVisitor {
	return &graphVisitor{
		visited: make(map[Module]struct{}),
		stack:   []Module{},
	}
}
