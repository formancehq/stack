package modules

import (
	"context"
	"fmt"
	"strings"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type StackReconcilerFactory struct {
	platform Platform
	client   client.Client
	scheme   *runtime.Scheme
}

func (sm *StackReconcilerFactory) NewDeployer(stack *v1beta3.Stack, configuration *v1beta3.Configuration, versions *v1beta3.Versions) *StackReconciler {
	return newStackReconciler(sm.client, sm.scheme, ReconciliationConfig{
		Stack:         stack,
		Configuration: configuration,
		Versions:      versions,
		Platform:      sm.platform,
	})
}

func NewsStackReconcilerFactory(client client.Client, scheme *runtime.Scheme, platform Platform) *StackReconcilerFactory {
	return &StackReconcilerFactory{
		platform: platform,
		client:   client,
		scheme:   scheme,
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
	podDeployer                PodDeployer
	portAllocator              PortAllocator
	namespacedResourceDeployer *ResourceDeployer

	ready collectionutils.Set[Module]
}

func newStackReconciler(client client.Client, scheme *runtime.Scheme, cfg ReconciliationConfig) *StackReconciler {

	resourceDeployer := NewDeployer(client, scheme, cfg.Stack, cfg.Configuration)

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

	registeredModules := RegisteredModules{}

	logger.Info("Prepare modules")
	err := r.prepareModules(ctx, registeredModules)
	if err != nil {
		return false, err
	}

	if finalizer, ok := r.podDeployer.(interface {
		finalize(context.Context) error
	}); ok {
		logger.Info("Finalize modules deployment")
		if err := finalizer.finalize(ctx); err != nil {
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
			logger.Info(fmt.Sprintf("Module '%s' marked as not ready", module.Name()))
		}
	}

	return allReady, nil
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

	isReady, err := newModuleReconciler(r, module).installModule(ctx, registeredModules)
	if err != nil {
		return err
	}
	if isReady {
		r.ready.Put(module)
	}

	return err
}

func (r *StackReconciler) deleteAllStackDeployments(ctx context.Context) error {
	return r.namespacedResourceDeployer.client.DeleteAllOf(ctx, &appsv1.Deployment{},
		client.InNamespace(r.Stack.Name),
		client.MatchingLabels{
			"stack": "true",
		},
	)
}

func (r *StackReconciler) checkDeployments(ctx context.Context) (bool, error) {
	logger := log.FromContext(ctx)
	deploymentsList := appsv1.DeploymentList{}
	if err := r.namespacedResourceDeployer.client.List(ctx, &deploymentsList,
		client.InNamespace(r.Stack.Name),
		client.MatchingLabels(map[string]string{
			stackLabel: "true",
		})); err != nil {
		return false, err
	}

	for _, deployment := range deploymentsList.Items {
		if deployment.Status.ObservedGeneration != deployment.Generation {
			logger.Info(fmt.Sprintf("Stop deployment as deployment '%s' is not ready (generation not matching)", deployment.Name))
			return false, nil
		}
		var moreRecentCondition appsv1.DeploymentCondition
		for _, condition := range deployment.Status.Conditions {
			if moreRecentCondition.Type == "" || condition.LastTransitionTime.After(moreRecentCondition.LastTransitionTime.Time) {
				moreRecentCondition = condition
			}
		}
		if moreRecentCondition.Type != appsv1.DeploymentAvailable {
			logger.Info(fmt.Sprintf("Stop deployment as deployment '%s' is not ready (last condition must be '%s', found '%s')", deployment.Name, appsv1.DeploymentAvailable, moreRecentCondition.Type))
			return false, nil
		}
		if moreRecentCondition.Status != "True" {
			logger.Info(fmt.Sprintf("Stop deployment as deployment '%s' is not ready ('%s' condition should be 'true')", deployment.Name, appsv1.DeploymentAvailable))
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
