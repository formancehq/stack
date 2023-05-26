package modules

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	v1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func HandleStack(ctx Context, deployer *ResourceDeployer) error {

	if ctx.Stack.Spec.Disabled {
		if err := deployer.client.DeleteAllOf(ctx, &v1.Deployment{},
			client.InNamespace(ctx.Stack.Name),
			client.MatchingLabels{
				"stack": "true",
			},
		); err != nil {
			return err
		}
		return nil
	}

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
		prepareContext PrepareContext
	}

	allServices := make(map[string]servicesWithContext)
	moduleNames := make([]string, 0)
	for moduleName := range modules {
		moduleNames = append(moduleNames, moduleName)
	}
	// Always process service in order to keep things idempotent
	sort.Strings(moduleNames)

	for _, moduleName := range moduleNames {
		module := modules[moduleName]
		prepareContext, services, err := module.prepare(ctx, portAllocator, moduleName)
		if err != nil {
			return err
		}

		allServices[moduleName] = servicesWithContext{
			services:       services,
			prepareContext: *prepareContext,
		}
	}

	registeredModules := RegisteredModules{}
	for moduleName, servicesWithContext := range allServices {
		registeredModules[moduleName] = RegisteredModule{
			Module:   modules[moduleName],
			Services: servicesWithContext.services,
		}
	}

	for _, moduleName := range moduleNames {
		holder := allServices[moduleName]
		if err := holder.services.install(ServiceInstallContext{
			PrepareContext:    holder.prepareContext,
			RegisteredModules: registeredModules,
			PodDeployer:       podDeployer,
		}, deployer, moduleName); err != nil {
			return err
		}
	}

	if finalizer, ok := podDeployer.(interface {
		finalize(context.Context) error
	}); ok {
		if err := finalizer.finalize(ctx); err != nil {
			return err
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
		if err := deployer.client.DeleteAllOf(ctx, &v1.Deployment{},
			client.InNamespace(ctx.Stack.Name),
			matchingLabels,
		); err != nil {
			return err
		}
	}

	ctx.Stack.Status.LightMode = ctx.Configuration.Spec.LightMode

	return nil
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

func (services Services) prepare(ctx PrepareContext, moduleName string) {
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

		err := service.Install(ctx, deployer, serviceName)
		if err != nil {
			me.setError(serviceName, err)
		}
	}
	if len(me.errors) > 0 {
		return me
	}
	return nil
}

type Module struct {
	Postgres func(ctx Context) v1beta3.PostgresConfig
	Services func(ctx Context) Services
}

func (module Module) prepare(ctx Context, portAllocator PortAllocator, moduleName string) (*PrepareContext, Services, error) {
	prepareContext := PrepareContext{
		Module:        moduleName,
		Context:       ctx,
		PortAllocator: portAllocator,
	}
	if module.Postgres != nil {
		postgresConfig := module.Postgres(ctx)
		conn, err := pgx.Connect(ctx, postgresConfig.DSN())
		if err != nil {
			return nil, nil, err
		}
		_, err = conn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, ctx.Stack.GetServiceName(moduleName)))
		if err != nil {
			pgErr := &pgconn.PgError{}
			if !errors.As(err, &pgErr) || pgErr.Code != "42P04" { // Database already exists error
				return nil, nil, err
			}
		}
		prepareContext.Postgres = &postgresConfig
	}

	services := module.Services(ctx)
	sort.Stable(services)
	services.prepare(prepareContext, moduleName)

	return &prepareContext, services, nil
}

var modules = map[string]Module{}

func Register(name string, module Module) {
	modules[name] = module
}
