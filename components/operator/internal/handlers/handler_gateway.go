package handlers

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/formancehq/operator/internal/modules"
)

const (
	gatewayPort = 8080
)

func init() {
	modules.Register("gateway", modules.Module{
		Services: func(ctx modules.Context) modules.Services {
			return modules.Services{{
				Port: gatewayPort,
				Path: "/",
				Configs: func(resolveContext modules.InstallContext) modules.Configs {
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
						Image:    modules.GetImage("gateway", resolveContext.Versions.Spec.Gateway),
						Liveness: modules.LivenessDisable,
					}
				},
			}}
		},
	})
}

func createCaddyfile(context modules.InstallContext) string {
	tpl := template.Must(template.New("caddyfile").Parse(caddyfile))
	buf := bytes.NewBufferString("")

	services := make(map[string]modules.Service)
	for moduleName, module := range context.RegisteredModules {
		if moduleName == "gateway" {
			continue
		}
		for _, service := range module.Services {
			if service.Port == 0 {
				continue
			}
			serviceName := moduleName
			if service.Name != "" {
				serviceName += "-" + service.Name
			}
			services[serviceName] = service
		}
	}

	if err := tpl.Execute(buf, map[string]any{
		"JWK_URL":  fmt.Sprintf("http://auth:%d/keys", context.RegisteredModules["auth"].Services[0].Port),
		"Services": services,
		"Debug":    context.Stack.Spec.Debug,
		"Fallback": fmt.Sprintf("control:%d", context.RegisteredModules["control"].Services[0].Port),
		"Port":     gatewayPort,
	}); err != nil {
		panic(err)
	}
	return buf.String()
}

const caddyfile = `
(cors) {
	header {
		Access-Control-Allow-Methods "GET,OPTIONS,PUT,POST,DELETE,HEAD,PATCH"
		Access-Control-Allow-Headers content-type
		Access-Control-Max-Age 100
		Access-Control-Allow-Origin *
	}
}

(handle_route_without_jwt) {
	# handle does not strips the prefix from the request path
	handle {args.0}/* {
		reverse_proxy {args.1}

		import cors
	}
}

(handle_path_route_with_jwt) {
	# handle_path automatically strips the prefix from the request path
	handle_path {args.0}* {
		reverse_proxy {args.1}

		import cors

		import jwt
	}
}

(handle_path_route_without_jwt) {
	# handle_path automatically strips the prefix from the request path
	handle_path {args.0}* {
		reverse_proxy {args.1}

		import cors
	}
}

(jwt) {
	jwtauth {
		sign_alg RS256
		jwk_url {{ .JWK_URL }}
		from_header Authorization
		refresh_window 10s
		min_refresh_interval 10s
	}
}

{
	# Many directives manipulate the HTTP handler chain and the order in which
	# those directives are evaluated matters. So the jwtauth directive must be
	# ordered.
	# c.f. https://caddyserver.com/docs/caddyfile/directives#directive-order
	order jwtauth before basicauth
	order versions after metrics

	{{ if .Debug }}debug{{ end }}
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

	{{- range $service, $target := .Services }}
		{{- if not (eq $service "control") }}
			{{- if not $target.Secured }}
	import handle_path_route_with_jwt /api/{{ $service }} {{ $service }}:{{ $target.Port }}
			{{- else }}
	import handle_path_route_without_jwt /api/{{ $service }} {{ $service }}:{{ $target.Port }}
			{{- end }}
		{{- end }}
	{{- end }}

	handle /versions {
		versions {
			{{- range $service, $target := .Services }}
				{{- if $target.HasVersionEndpoint }}
			{{ $service }} http://{{ $service }}:{{ $target.Port}}/_info
				{{- end }}
			{{- end }}
		}
	}

	# handle all other requests
	handle {
		reverse_proxy {{ .Fallback }}
		import cors
	}
}
`
