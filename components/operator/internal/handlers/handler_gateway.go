package handlers

import (
	"bytes"
	"fmt"
	"sort"
	"text/template"

	"github.com/formancehq/operator/internal/modules"
)

const (
	gatewayPort = 8000
)

func init() {
	modules.Register("gateway", modules.Module{
		Versions: map[string]modules.Version{
			"v0.0.0": {
				Services: func(ctx modules.ModuleContext) modules.Services {
					return modules.Services{{
						Port:        gatewayPort,
						Path:        "/",
						ExposeHTTP:  true,
						Liveness:    modules.LivenessDisable,
						Annotations: ctx.Configuration.Spec.Services.Gateway.Annotations.Service,
						Configs: func(resolveContext modules.ServiceInstallContext) modules.Configs {
							return modules.Configs{
								"config": modules.Config{
									Data: map[string]string{
										"Caddyfile": createCaddyfile(resolveContext),
									},
									Mount: true,
								},
							}
						},
						Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
							return modules.Container{
								Command: []string{"/usr/bin/caddy"},
								Args: []string{
									"run",
									"--config", resolveContext.GetConfig("config").GetMountPath() + "/Caddyfile",
									"--adapter", "caddyfile",
								},
								Image: modules.GetImage("gateway", resolveContext.Versions.Spec.Gateway),
								Env:   modules.NewEnv(),
								Resources: getResourcesWithDefault(
									resolveContext.Configuration.Spec.Services.Gateway.ResourceProperties,
									modules.ResourceSizeSmall(),
								),
							}
						},
					}}
				},
			},
		},
	})
}

func createCaddyfile(context modules.ServiceInstallContext) string {
	tpl := template.Must(template.New("caddyfile").Parse(caddyfile))
	buf := bytes.NewBufferString("")

	type service struct {
		Name string
		*modules.Service
		Port       int32
		Hostname   string
		HealthPath string
	}

	servicesMap := make(map[string]service, 0)
	keys := make([]string, 0)
	for moduleName, module := range context.RegisteredModules {
		if moduleName == "gateway" {
			continue
		}
		for _, s := range module.Services {
			usedPort := s.GetUsedPort()
			if usedPort == 0 {
				continue
			}
			serviceName := moduleName
			if s.Name != "" {
				serviceName += "-" + s.Name
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
				Name:       serviceName,
				Service:    s,
				Port:       usedPort,
				Hostname:   hostname,
				HealthPath: healthPath,
			}
			keys = append(keys, serviceName)
		}
	}

	sort.Strings(keys)
	services := make([]service, 0)
	for _, key := range keys {
		services = append(services, servicesMap[key])
	}

	if err := tpl.Execute(buf, map[string]any{
		"Region":   context.Region,
		"Env":      context.Environment,
		"Issuer":   fmt.Sprintf("%s/api/auth", context.Stack.URL()),
		"Services": services,
		"Debug":    context.Stack.Spec.Debug,
		"Fallback": fmt.Sprintf("control:%d", context.RegisteredModules["control"].Services[0].Port),
		"Port":     gatewayPort,
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

(handle_route_without_auth) {
	# handle does not strips the prefix from the request path
	handle {args.0}/* {
		reverse_proxy {args.1}

		import cors
	}
}

(handle_path_route_with_auth) {
	# handle_path automatically strips the prefix from the request path
	handle_path {args.0}* {
		reverse_proxy {args.1}

		import cors

		import auth
	}
}

(handle_path_route_without_auth) {
	# handle_path automatically strips the prefix from the request path
	handle_path {args.0}* {
		reverse_proxy {args.1}

		import cors
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
			{{- if not $service.Secured }}
	import handle_path_route_with_auth /api/{{ $service.Name }} {{ $service.Hostname }}:{{ $service.Port }}
			{{- else }}
	import handle_path_route_without_auth /api/{{ $service.Name }} {{ $service.Hostname }}:{{ $service.Port }}
			{{- end }}
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
