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
	"golang.org/x/mod/semver"

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
							Env: modules.NewEnv().
								Append(modules.BrokerEnvVarsWithPrefix(resolveContext.Configuration.Spec.Broker, "gateway")...),
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

	fallback := ""
	redirect := false
	if context.Configuration.Spec.Services.Gateway.Fallback != nil && *context.Configuration.Spec.Services.Gateway.Fallback != "" {
		fallback = *context.Configuration.Spec.Services.Gateway.Fallback
		redirect = true
	} else {
		if !context.IsDisabled("control") {
			fallback = fmt.Sprintf("control:%d", servicesMap["control"].Port)
		}
	}

	if err := caddyfileTemplate.Execute(buf, map[string]any{
		"Region":   context.Platform.Region,
		"Env":      context.Platform.Environment,
		"Issuer":   fmt.Sprintf("%s/api/auth", context.Stack.URL()),
		"Services": services,
		"Debug":    context.Stack.Spec.Debug,
		"Broker": func() string {
			if context.Configuration.Spec.Broker.Kafka != nil {
				return "kafka"
			}
			if context.Configuration.Spec.Broker.Nats != nil {
				return "nats"
			}
			return ""
		}(),
		"Port":     gatewayPort,
		"Fallback": fallback,
		"Redirect": redirect,
		"EnableAudit": func() bool {
			gatewayVersion := context.Versions.Spec.Gateway
			enable := context.Configuration.Spec.Services.Gateway.EnableAuditPlugin
			if enable == nil {
				return false
			}

			if !semver.IsValid(gatewayVersion) {
				return *enable
			}

			return *enable && semver.Compare("v0.2.0", gatewayVersion) <= 0
		}(),
	}); err != nil {
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

{{- if .EnableAudit }}
(audit) {
	audit {
		# Kafka publisher
		{{- if (eq .Broker "kafka") }}
		publisher_kafka_broker {$PUBLISHER_KAFKA_BROKER:redpanda:29092}
		publisher_kafka_enabled {$PUBLISHER_KAFKA_ENABLED:false}
		publisher_kafka_tls_enabled {$PUBLISHER_KAFKA_TLS_ENABLED:false}
		publisher_kafka_sasl_enabled {$PUBLISHER_KAFKA_SASL_ENABLED:false}
		publisher_kafka_sasl_username {$PUBLISHER_KAFKA_SASL_USERNAME}
		publisher_kafka_sasl_password {$PUBLISHER_KAFKA_SASL_PASSWORD}
		publisher_kafka_sasl_mechanism {$PUBLISHER_KAFKA_SASL_MECHANISM}
		publisher_kafka_sasl_scram_sha_size {$PUBLISHER_KAFKA_SASL_SCRAM_SHA_SIZE}
		{{- end }}
		{{- if (eq .Broker "nats") }}
		# Nats publisher
		publisher_nats_enabled {$PUBLISHER_NATS_ENABLED:true}
		publisher_nats_url {$PUBLISHER_NATS_URL:nats://nats:4222}
		publisher_nats_client_id {$PUBLISHER_NATS_CLIENT_ID:gateway}
		{{- end }}
	}
}
{{- end }}

{
	{{ if .Debug }}debug{{ end }}
	# Many directives manipulate the HTTP handler chain and the order in which
	# those directives are evaluated matters. So the jwtauth directive must be
	# ordered.
	# c.f. https://caddyserver.com/docs/caddyfile/directives#directive-order
	order auth before basicauth
	order versions after metrics
	order audit after encode
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

	{{- if .EnableAudit }}
	import audit
	{{- end }}

	{{- range $i, $service := .Services }}
		{{- if not (eq $service.Name "control") }}
			{{- range $i, $path := $service.Paths }}
			@{{ $path.Name }}matcher {
				path /api/{{ $service.RoutingPath }}{{ $path.Path }}*
				{{- if gt ($path.Methods | len) 0 }}
				method {{ join $path.Methods " " }}
				{{- end }}
			}
			handle @{{ $path.Name }}matcher {
				uri strip_prefix /api/{{ $service.RoutingPath }}
				reverse_proxy {{ $service.Hostname }}:{{ $service.Port }}
				import cors
				{{- if not $service.Secured }}
				import auth
				{{- end }}
			}
			{{- end }}

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

	# Respond 502 if service does not exists
	handle /api/* {
		respond "Bad Gateway" 502
	}

	# handle all other requests
	{{- if not (eq .Fallback "") }}
	handle {
		{{- if .Redirect }}
		redir {{ .Fallback }}
		{{- else }}
		reverse_proxy {{ .Fallback }}
		{{- end }}
		import cors
	}
	{{ end }}
}`
