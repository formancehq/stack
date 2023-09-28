package modules

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type StackDeployer struct {
	roundTripper http.RoundTripper
}

func NewStackDeployer(roundTripper http.RoundTripper) *StackDeployer {
	return &StackDeployer{
		roundTripper: roundTripper,
	}
}

func (sd *StackDeployer) HandleStack(ctx Context, deployer *ResourceDeployer) (bool, error) {

	logger := log.FromContext(ctx)
	logger = logger.WithValues("stack", ctx.Stack.Name)

	var (
		portAllocator PortAllocator = StaticPortAllocator(8080)
		podDeployer   PodDeployer   = NewDefaultPodDeployer(deployer)
	)

	if ctx.Configuration.Spec.LightMode {
		podDeployer = NewMonoPodDeployer(deployer, ctx.Stack.Name)
		portAllocator = NewPortRangeAllocator(10000)
	}

	type servicesWithContext struct {
		services       Services
		prepareContext ModuleContext
	}

	allServices := make(map[string]servicesWithContext)
	moduleNames := make([]string, 0)
	// When Service in Stack is Disabled, we want to remove the deployment
	// TODO: It's possible to remove more than one deployment or another resource
	for moduleName := range modules {
		if ctx.Stack.Spec.Services.IsDisabled(moduleName) {
			if err := deployer.client.DeleteAllOf(ctx, &appsv1.Deployment{},
				client.InNamespace(ctx.Stack.Name),
				client.MatchingLabels{
					"app.kubernetes.io/name": moduleName,
					"stack":                  "true",
				},
			); err != nil {
				return false, err
			}
			continue
		}
		moduleNames = append(moduleNames, moduleName)
	}
	// When Stack is Disabled, we want to remove all deployments
	if ctx.Stack.Spec.Disabled {
		logger.Info("Stack is disabled, remove all deployments")
		if err := deployer.client.DeleteAllOf(ctx, &appsv1.Deployment{},
			client.InNamespace(ctx.Stack.Name),
			client.MatchingLabels{
				"stack": "true",
			},
		); err != nil {
			return false, err
		}
		return true, nil
	}
	logger.Info(fmt.Sprintf("List of Modules activated: %s", strings.Join(moduleNames, ",")))

	// Always process service in order to keep things idempotent
	sort.Strings(moduleNames)

	allReady := true
	for _, moduleName := range moduleNames {
		logger.Info(fmt.Sprintf("Pre install module '%s'", moduleName))
		module := modules[moduleName]
		prepareContext, services, ready, err := preInstallModule(ctx, module, deployer, portAllocator, moduleName)
		if err != nil {
			return false, err
		}
		allReady = allReady && ready
		if !ready {
			logger.Info(fmt.Sprintf("Module '%s' marked as not ready", moduleName))
			continue
		}

		allServices[moduleName] = servicesWithContext{
			services:       services,
			prepareContext: *prepareContext,
		}
	}
	if !allReady {
		return false, nil
	}

	registeredModules := RegisteredModules{}
	for moduleName, servicesWithContext := range allServices {
		registeredModules[moduleName] = RegisteredModule{
			Module:   modules[moduleName],
			Services: servicesWithContext.services,
		}
	}

	for _, moduleName := range moduleNames {
		logger.Info(fmt.Sprintf("Install module '%s'", moduleName))
		holder := allServices[moduleName]
		if err := holder.services.install(ServiceInstallContext{
			ModuleContext:     holder.prepareContext,
			RegisteredModules: registeredModules,
			PodDeployer:       podDeployer,
		}, deployer, moduleName); err != nil {
			return false, err
		}
	}

	if finalizer, ok := podDeployer.(interface {
		finalize(context.Context) error
	}); ok {
		if err := finalizer.finalize(ctx); err != nil {
			return false, err
		}
	}

	var matchingLabels client.MatchingLabels
	switch {
	case ctx.Stack.Status.LightMode && !ctx.Configuration.Spec.LightMode:
		matchingLabels = map[string]string{
			monopodLabel: "true",
		}
	case !ctx.Stack.Status.LightMode && ctx.Configuration.Spec.LightMode:
		matchingLabels = map[string]string{
			monopodLabel: "false",
		}
	}
	if matchingLabels != nil {
		logger.Info("Delete orphan deployments")

		if err := deployer.client.DeleteAllOf(ctx, &appsv1.Deployment{},
			client.InNamespace(ctx.Stack.Name),
			matchingLabels,
		); err != nil {
			return false, err
		}
	}

	ctx.Stack.Status.LightMode = ctx.Configuration.Spec.LightMode

	deploymentsList := appsv1.DeploymentList{}
	if err := deployer.client.List(ctx, &deploymentsList,
		client.InNamespace(ctx.Stack.Name),
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

	allReady = true
	for _, moduleName := range moduleNames {
		logger.Info(fmt.Sprintf("Post install module '%s'", moduleName))
		module := modules[moduleName]
		ready, err := postInstall(ctx, module, deployer, moduleName)
		if err != nil {
			return false, err
		}
		allReady = allReady && ready
		if !ready {
			logger.Info(fmt.Sprintf("Module '%s' marked as not ready", moduleName))
		}
	}

	return allReady, nil
}

type Services []*Service

func (services Services) Len() int {
	return len(services)
}

func (services Services) Less(i, j int) bool {
	return strings.Compare(services[i].Name, services[j].Name) < 0
}

func (services Services) Swap(i, j int) {
	services[i], services[j] = services[j], services[i]
}

func (services Services) preInstall(ctx ModuleContext, moduleName string) {
	for _, service := range services {
		serviceName := moduleName
		if service.Name != "" {
			serviceName = serviceName + "-" + service.Name
		}

		service.Prepare(ctx, serviceName)
	}
}

func (services Services) install(ctx ServiceInstallContext, deployer *ResourceDeployer, moduleName string) error {
	me := &serviceErrors{}
	for _, service := range services {
		serviceName := moduleName
		if service.Name != "" {
			serviceName = serviceName + "-" + service.Name
		}

		err := service.install(ctx, deployer, serviceName)
		if err != nil {
			me.setError(serviceName, err)
		}
	}
	if len(me.errors) > 0 {
		return me
	}
	return nil
}

// CreatePostgresDatabase Ugly hack to allow mocking
var CreatePostgresDatabase = func(ctx context.Context, dsn, dbName string) error {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, dbName))
	if err != nil {
		pgErr := &pgconn.PgError{}
		if !errors.As(err, &pgErr) || pgErr.Code != "42P04" { // Database already exists error
			return err
		}
	}

	return nil
}

type Cron struct {
	Container Container
	Schedule  string
	Suspend   bool
}

type Version struct {
	PreUpgrade  func(ctx Context) error
	PostUpgrade func(ctx PostInstallContext) error
	Services    func(ctx ModuleContext) Services
	Cron        func(ctx Context) []Cron
}

type Module interface {
	Versions() map[string]Version
}

type PostgresAwareModule interface {
	Postgres(ctx Context) v1beta3.PostgresConfig
}

func preInstallModule(ctx Context, module Module, deployer *ResourceDeployer, portAllocator PortAllocator, moduleName string) (*ModuleContext, Services, bool, error) {
	moduleContext := ModuleContext{
		Module:        moduleName,
		Context:       ctx,
		PortAllocator: portAllocator,
	}
	if postgresAwareModule, ok := module.(PostgresAwareModule); ok {
		postgresConfig := postgresAwareModule.Postgres(ctx)
		if err := CreatePostgresDatabase(ctx, postgresConfig.DSN(), ctx.Stack.GetServiceName(moduleName)); err != nil {
			return nil, nil, false, err
		}
		moduleContext.Postgres = &postgresConfig
	}

	versions := module.Versions()

	versionsKeys := collectionutils.Keys(versions)
	sort.Strings(versionsKeys)

	var chosenVersion Version
	for _, version := range versionsKeys {
		if !moduleContext.HasVersionHigherOrEqual(version) {
			break
		}
		chosenVersion = versions[version]
		if chosenVersion.PreUpgrade == nil {
			continue
		}

		migration := &v1beta3.Migration{}
		migrationName := fmt.Sprintf("%s-%s-pre-upgrade", moduleName, version)
		if err := deployer.client.Get(ctx, types.NamespacedName{
			Namespace: ctx.Stack.Name,
			Name:      migrationName,
		}, migration); err != nil {
			if !apierrors.IsNotFound(err) {
				return nil, nil, false, err
			}
			_, err := deployer.Migrations().CreateOrUpdate(ctx, migrationName, func(t *v1beta3.Migration) {
				t.Spec = v1beta3.MigrationSpec{
					Configuration:   ctx.Configuration.Name,
					Module:          moduleName,
					TargetedVersion: version,
					Version:         ctx.Versions.Name,
				}
			})
			if err != nil {
				return nil, nil, false, err
			}
			return nil, nil, false, nil
		}
		if !migration.Status.Terminated {
			return nil, nil, false, nil
		}
	}

	services := chosenVersion.Services(moduleContext)
	sort.Stable(services)
	services.preInstall(moduleContext, moduleName)

	return &moduleContext, services, true, nil
}

func postInstall(ctx Context, module Module, deployer *ResourceDeployer, moduleName string) (bool, error) {

	versions := module.Versions()

	versionsKeys := collectionutils.Keys(versions)
	sort.Strings(versionsKeys)

	var selectedVersion Version
	for _, version := range versionsKeys {
		if !ctx.Versions.IsHigherOrEqual(moduleName, version) {
			break
		}
		selectedVersion = versions[version]
		if selectedVersion.PostUpgrade == nil {
			continue
		}

		migration := &v1beta3.Migration{}
		migrationName := fmt.Sprintf("%s-%s-post-upgrade", moduleName, version)
		if err := deployer.client.Get(ctx, types.NamespacedName{
			Namespace: ctx.Stack.Name,
			Name:      migrationName,
		}, migration); err != nil {
			if !apierrors.IsNotFound(err) {
				return false, err
			}
			_, err := deployer.Migrations().CreateOrUpdate(ctx, migrationName, func(t *v1beta3.Migration) {
				t.Spec = v1beta3.MigrationSpec{
					Configuration:   ctx.Configuration.Name,
					Module:          moduleName,
					TargetedVersion: version,
					Version:         ctx.Versions.Name,
					PostUpgrade:     true,
				}
			})
			if err != nil {
				return false, err
			}
			return false, nil
		}
		if !migration.Status.Terminated {
			return false, nil
		}
	}

	if selectedVersion.Cron != nil {
		for _, cron := range selectedVersion.Cron(ctx) {
			_, err := deployer.CronJobs().CreateOrUpdate(ctx, cron.Container.Name, func(t *batchv1.CronJob) {
				t.Spec = batchv1.CronJobSpec{
					Suspend:  &cron.Suspend,
					Schedule: cron.Schedule,
					JobTemplate: batchv1.JobTemplateSpec{
						Spec: batchv1.JobSpec{
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									RestartPolicy: corev1.RestartPolicyNever,
									Containers: []corev1.Container{{
										Name:    cron.Container.Name,
										Image:   cron.Container.Image,
										Command: cron.Container.Command,
										Args:    cron.Container.Args,
										Env:     cron.Container.Env.ToCoreEnv(),
									}},
								},
							},
						},
					},
				}
			})
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

var modules = map[string]Module{}

func Register(name string, module Module) {
	modules[name] = module
}

func Get(name string) Module {
	return modules[name]
}
