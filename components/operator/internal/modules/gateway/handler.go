package gateway

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/formancehq/operator/internal/modules/auth"
	"github.com/formancehq/operator/internal/modules/control"
	"github.com/formancehq/operator/internal/modules/ledger"
	"github.com/formancehq/operator/internal/modules/orchestration"
	"github.com/formancehq/operator/internal/modules/payments"
	"github.com/formancehq/operator/internal/modules/search"
	"github.com/formancehq/operator/internal/modules/wallets"
	"github.com/formancehq/operator/internal/modules/webhooks"

	"github.com/formancehq/operator/internal/modules"
)

const (
	gatewayPort = 8000
)

type module struct{}

func (g module) Name() string {
	return "gateway"
}

func (g module) DependsOn() []modules.Module {
	return []modules.Module{
		auth.Module,
		control.Module,
		ledger.Module,
		orchestration.Module,
		payments.Module,
		search.Module,
		wallets.Module,
		webhooks.Module,
	}
}

func (g module) Versions() map[string]modules.Version {
	return map[string]modules.Version{
		"v0.0.0": {
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				return modules.Services{{
					Port: gatewayPort,
					ExposeHTTP: &modules.ExposeHTTP{
						Path: "/",
					},
					Liveness:    modules.LivenessDisable,
					Annotations: ctx.Configuration.Spec.Services.Gateway.Annotations.Service,
					Configs: func(resolveContext modules.ServiceInstallConfiguration) modules.Configs {
						return modules.Configs{
							"config": modules.Config{
								Data: map[string]string{
									"Caddyfile": createCaddyfile(resolveContext),
								},
								Mount: true,
							},
						}
					},
					Container: func(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
						return modules.Container{
							Command: []string{"/usr/bin/caddy"},
							Args: []string{
								"run",
								"--config", resolveContext.GetConfig("config").GetMountPath() + "/Caddyfile",
								"--adapter", "caddyfile",
							},
							Image: modules.GetImage("gateway", resolveContext.Versions.Spec.Gateway),
							Env:   modules.NewEnv(),
							Resources: modules.GetResourcesWithDefault(
								resolveContext.Configuration.Spec.Services.Gateway.ResourceProperties,
								modules.ResourceSizeSmall(),
							),
						}
					},
				}}
			},
		},
	}
}

var Module = &module{}

var _ modules.Module = Module
var _ modules.DependsOnAwareModule = Module

var caddyfileTemplate = template.Must(template.New("caddyfile").Funcs(map[string]any{
	"join": strings.Join,
}).Parse(caddyfile))

func init() {
	modules.Register(Module)
}

func createCaddyfile(context modules.ServiceInstallConfiguration) string {
	buf := bytes.NewBufferString("")

	type service struct {
		modules.RegisteredService
		Name        string
		Port        int32
		Hostname    string
		HealthPath  string
		Methods     []string
		RoutingPath string
	}

	servicesMap := make(map[string]service, 0)
	keys := make([]string, 0)
	for _, registeredModule := range context.RegisteredModules {
		if registeredModule.Module.Name() == "gateway" {
			continue
		}
		for _, s := range registeredModule.Services {
			if s.ExposeHTTP == nil {
				continue
			}
			usedPort := s.Port
			if usedPort == 0 {
				continue
			}
			serviceName := registeredModule.Module.Name()
			if s.Name != "" {
				serviceName += "-" + s.Name
			}
			serviceRoutingPath := registeredModule.Module.Name()
			if s.ExposeHTTP.Name != "" {
				serviceRoutingPath = serviceRoutingPath + "-" + s.ExposeHTTP.Name
			}

			healthPath := "_healthcheck"
			if s.Liveness == modules.LivenessLegacy {
				healthPath = "_health"
			}
			hostname := serviceName
			if context.Configuration.Spec.LightMode {
				hostname = "127.0.0.1"
			}

			servicesMap[serviceName] = service{
				Name:              serviceName,
				RegisteredService: s,
				Port:              usedPort,
				Hostname:          hostname,
				HealthPath:        healthPath,
				Methods:           s.ExposeHTTP.Methods,
				RoutingPath:       serviceRoutingPath,
			}
			keys = append(keys, serviceName)
		}
	}

	sort.Strings(keys)
	services := make([]service, 0)
	for _, key := range keys {
		services = append(services, servicesMap[key])
	}

	data := map[string]any{
		"Region":   context.Platform.Region,
		"Env":      context.Platform.Environment,
		"Issuer":   fmt.Sprintf("%s/api/auth", context.Stack.URL()),
		"Services": services,
		"Debug":    context.Stack.Spec.Debug,
		"Fallback": fmt.Sprintf("control:%d", servicesMap["control"].Port),
		"Port":     gatewayPort,
	}
	control, ok := context.RegisteredModules["control"]
	if ok {
		data["Fallback"] = fmt.Sprintf("control:%d", control.Services["control"].Port)
	}

	if err := caddyfileTemplate.Execute(buf, data); err != nil {
		panic(err)
	}
	return buf.String()
}

const caddyfile = `(cors) {
	header {
		Access-Control-Allow-Methods "GET,OPTIONS,PUT,POST,DELETE,HEAD,PATCH"
		Access-Control-Allow-Headers content-type
		Access-Control-Max-Age 100
		Access-Control-Allow-Origin *
	}
}

(auth) {
	auth {
		issuer {{ .Issuer }}

		read_key_set_max_retries 10
	}
}

{
	{{ if .Debug }}debug{{ end }}

	# Many directives manipulate the HTTP handler chain and the order in which
	# those directives are evaluated matters. So the jwtauth directive must be
	# ordered.
	# c.f. https://caddyserver.com/docs/caddyfile/directives#directive-order
	order auth before basicauth
	order versions after metrics
}

:{{ .Port }} {
	tracing {
		span gateway
	}
	log {
		output stdout
		{{- if .Debug }}
		level  DEBUG
		{{- end }}
	}

	{{- range $i, $service := .Services }}
		{{- if not (eq $service.Name "control") }}
			@{{ $service.Name }}matcher {
				path /api/{{ $service.RoutingPath }}*
				{{- if gt ($service.Methods | len) 0 }}
				method {{ join $service.Methods " " }}
				{{- end }}
			}
			handle @{{ $service.Name }}matcher {
				uri strip_prefix /api/{{ $service.RoutingPath }}
				reverse_proxy {{ $service.Hostname }}:{{ $service.Port }}
		
				import cors
				{{- if not $service.Secured }}
				import auth
				{{- end }}
			}
		{{- end }}
	{{- end }}

	handle /versions {
		versions {
			region "{{ .Region }}"
			env "{{ .Env }}"
			endpoints {
				{{- range $i, $service := .Services }}
					{{- if $service.HasVersionEndpoint }}
				{{ $service.Name }} http://{{ $service.Hostname }}:{{ $service.Port }}/_info http://{{ $service.Hostname }}:{{ $service.Port }}/{{ $service.HealthPath }}
					{{- end }}
				{{- end }}
			}
		}
	}

	# handle all other requests
	handle {
		reverse_proxy {{ .Fallback }}
		import cors
	}
}`
