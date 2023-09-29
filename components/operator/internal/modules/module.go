package modules

import (
	"context"
	"fmt"
	"sort"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Cron struct {
	Container Container
	Schedule  string
	Suspend   bool
}

type Version struct {
	Services    func(cfg ReconciliationConfig) Services
	Cron        func(cfg ReconciliationConfig) []Cron
	PreUpgrade  func(ctx context.Context, cfg ReconciliationConfig) error
	PostUpgrade func(ctx context.Context, cfg ReconciliationConfig) error
}

type Module interface {
	Name() string
	Versions() map[string]Version
}

type DependsOnAwareModule interface {
	Module
	DependsOn() []Module
}

type PostgresAwareModule interface {
	Module
	Postgres(cfg ReconciliationConfig) v1beta3.PostgresConfig
}

var modules = make([]Module, 0)

func Register(newModules ...Module) {
	modules = append(modules, newModules...)
}

func Get(name string) Module {
	for _, module := range modules {
		if module.Name() == name {
			return module
		}
	}
	return nil
}

func sortedVersions(module Module) []string {
	versionsKeys := collectionutils.Keys(module.Versions())
	sort.Strings(versionsKeys)

	return versionsKeys
}

func falseIfError(err error) (bool, error) {
	if err != nil {
		return false, err
	}
	return true, nil
}

type moduleReconciler struct {
	*StackReconciler
	module Module
}

func (r *moduleReconciler) installModule(ctx context.Context, registeredModules RegisteredModules) (bool, error) {

	log.FromContext(ctx).Info(fmt.Sprintf("Installing module %s", r.module.Name()))

	registeredModule := RegisteredModule{
		Module:   r.module,
		Services: map[string]RegisteredService{},
	}
	registeredModules[r.module.Name()] = registeredModule

	var (
		postgresConfig *v1beta3.PostgresConfig
		err            error
	)
	pam, ok := r.module.(PostgresAwareModule)
	if ok {
		postgresConfig, err = r.createDatabase(ctx, pam)
		if err != nil {
			return false, err
		}
	}

	var (
		isReady       = true
		chosenVersion Version
	)
	for _, version := range sortedVersions(r.module) {
		if !r.Versions.IsHigherOrEqual(r.module.Name(), version) {
			break
		}
		chosenVersion = r.module.Versions()[version]
		if chosenVersion.PreUpgrade == nil {
			continue
		}

		ready, err := r.runPreUpgradeMigration(ctx, r.module, version)
		if err != nil {
			return false, err
		}
		isReady = isReady && ready
	}

	// Stop install if we are not in light mode and all migrations are not passed
	if !isReady && !r.Configuration.Spec.LightMode {
		log.FromContext(ctx).Info(fmt.Sprintf("Stop install as module '%s' is not ready and stack is not in light mode", r.module.Name()))
		return false, nil
	}

	services := chosenVersion.Services(r.ReconciliationConfig)
	sort.Stable(services)

	me := &serviceErrors{}
	for _, service := range services {
		serviceName := r.module.Name()
		if service.Name != "" {
			serviceName = serviceName + "-" + service.Name
		}

		serviceReconciler := newServiceReconciler(r, *service, serviceName)
		err := serviceReconciler.reconcile(ctx, ServiceInstallConfiguration{
			ReconciliationConfig: r.ReconciliationConfig,
			RegisteredModules:    registeredModules,
			PostgresConfig:       postgresConfig,
		})
		if err != nil {
			me.setError(serviceName, err)
		}

		registeredModule.Services[serviceName] = RegisteredService{
			Port:    serviceReconciler.usedPort,
			Service: *service,
		}
	}
	if len(me.errors) > 0 {
		return false, me
	}

	return true, nil
}

func (r *moduleReconciler) finalizeModule(ctx context.Context, module Module) (bool, error) {
	versions := module.Versions()

	var selectedVersion Version
	for _, version := range sortedVersions(module) {
		if !r.Versions.IsHigherOrEqual(module.Name(), version) {
			break
		}
		selectedVersion = versions[version]
		if selectedVersion.PostUpgrade == nil {
			continue
		}

		migration := &v1beta3.Migration{}
		migrationName := fmt.Sprintf("%s-%s-post-upgrade", module.Name(), version)
		if err := r.namespacedResourceDeployer.client.Get(ctx, types.NamespacedName{
			Namespace: r.Stack.Name,
			Name:      migrationName,
		}, migration); err != nil {
			if !apierrors.IsNotFound(err) {
				return false, err
			}
			_, err := r.namespacedResourceDeployer.Migrations().CreateOrUpdate(ctx, migrationName, func(t *v1beta3.Migration) {
				t.Spec = v1beta3.MigrationSpec{
					Configuration:   r.Configuration.Name,
					Module:          module.Name(),
					TargetedVersion: version,
					Version:         r.Versions.Name,
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
		for _, cron := range selectedVersion.Cron(r.ReconciliationConfig) {
			_, err := r.namespacedResourceDeployer.CronJobs().CreateOrUpdate(ctx, cron.Container.Name, func(t *batchv1.CronJob) {
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

func (r *moduleReconciler) createDatabase(ctx context.Context, module PostgresAwareModule) (*v1beta3.PostgresConfig, error) {
	postgresConfig := module.Postgres(r.ReconciliationConfig)
	if err := CreatePostgresDatabase(ctx, postgresConfig.DSN(), r.Stack.GetServiceName(module.Name())); err != nil {
		return nil, err
	}

	return &postgresConfig, nil
}

func (r *moduleReconciler) runPreUpgradeMigration(ctx context.Context, module Module, version string) (bool, error) {
	migration := &v1beta3.Migration{}
	migrationName := fmt.Sprintf("%s-%s-pre-upgrade", module.Name(), version)
	if err := r.namespacedResourceDeployer.client.Get(ctx, types.NamespacedName{
		Namespace: r.Stack.Name,
		Name:      migrationName,
	}, migration); err != nil {
		if !apierrors.IsNotFound(err) {
			return false, err
		}
		_, err := r.namespacedResourceDeployer.Migrations().CreateOrUpdate(ctx, migrationName, func(t *v1beta3.Migration) {
			t.Spec = v1beta3.MigrationSpec{
				Configuration:   r.Configuration.Name,
				Module:          module.Name(),
				TargetedVersion: version,
				Version:         r.Versions.Name,
			}
		})
		return false, err
	}

	return migration.Status.Terminated, nil
}

func newModuleReconciler(stackReconciler *StackReconciler, module Module) *moduleReconciler {
	return &moduleReconciler{
		StackReconciler: stackReconciler,
		module:          module,
	}
}
