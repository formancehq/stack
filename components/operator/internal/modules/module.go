package modules

import (
	"errors"
	"fmt"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func HandleStack(ctx Context, deployer *StackWideDeployer) error {

	type servicesWithContext struct {
		services       Services
		prepareContext PrepareContext
	}

	allServices := make(map[string]servicesWithContext)
	for moduleName, module := range modules {
		serviceContext, err := module.prepare(ctx, moduleName)
		if err != nil {
			return err
		}

		services := module.Services(ctx)
		services.prepare(*serviceContext, moduleName)
		allServices[moduleName] = servicesWithContext{
			services:       services,
			prepareContext: *serviceContext,
		}
	}

	registeredModules := RegisteredModules{}
	for moduleName, servicesWithContext := range allServices {
		registeredModules[moduleName] = RegisteredModule{
			Module:   modules[moduleName],
			Services: servicesWithContext.services,
		}
	}

	for moduleName, holder := range allServices {
		err := holder.services.install(InstallContext{
			PrepareContext:    holder.prepareContext,
			RegisteredModules: registeredModules,
		}, deployer, moduleName)
		if err != nil {
			return err
		}
	}

	return nil
}

type Services []Service

func (services Services) prepare(ctx PrepareContext, moduleName string) {
	for _, service := range services {
		serviceName := moduleName
		if service.Name != "" {
			serviceName = serviceName + "-" + service.Name
		}

		service.Prepare(ctx, serviceName)
	}
}

func (services Services) install(ctx InstallContext, deployer *StackWideDeployer, moduleName string) error {
	me := &serviceErrors{}
	for _, service := range services {
		serviceName := moduleName
		if service.Name != "" {
			serviceName = serviceName + "-" + service.Name
		}

		err := service.Install(ctx, deployer.ForService(serviceName), serviceName)
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

func (module Module) prepare(ctx Context, moduleName string) (*PrepareContext, error) {
	if module.Postgres != nil {
		postgresConfig := module.Postgres(ctx)
		conn, err := pgx.Connect(ctx, postgresConfig.DSN())
		if err != nil {
			return nil, err
		}
		_, err = conn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, ctx.Stack.GetServiceName(moduleName)))
		if err != nil {
			pgErr := &pgconn.PgError{}
			if !errors.As(err, &pgErr) || pgErr.Code != "42P04" { // Database already exists error
				return nil, err
			}
		}
		return &PrepareContext{
			Context:  ctx,
			Postgres: &postgresConfig,
			Module:   moduleName,
		}, nil
	}
	return &PrepareContext{
		Module:  moduleName,
		Context: ctx,
	}, nil
}

var modules = map[string]Module{}

func Register(name string, module Module) {
	modules[name] = module
}
