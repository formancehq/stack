package temporal

import (
	"context"
	"crypto/tls"

	"github.com/formancehq/orchestration/internal/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"
	"go.uber.org/fx"
)

func NewClientModule(address, namespace string, certStr string, key string) fx.Option {
	return fx.Options(
		fx.Provide(func() (client.Options, error) {

			var cert *tls.Certificate
			if key != "" && certStr != "" {
				clientCert, err := tls.X509KeyPair([]byte(certStr), []byte(key))
				if err != nil {
					return client.Options{}, err
				}
				cert = &clientCert
			}

			tracingInterceptor, err := opentelemetry.NewTracingInterceptor(opentelemetry.TracerOptions{
				Tracer: workflow.Tracer,
			})
			if err != nil {
				return client.Options{}, err
			}

			options := client.Options{
				Namespace:    namespace,
				HostPort:     address,
				Interceptors: []interceptor.ClientInterceptor{tracingInterceptor},
			}
			if cert != nil {
				options.ConnectionOptions = client.ConnectionOptions{
					TLS: &tls.Config{Certificates: []tls.Certificate{*cert}},
				}
			}
			return options, nil
		}),
		fx.Provide(client.Dial),
		fx.Invoke(func(lifecycle fx.Lifecycle, c client.Client) {
			lifecycle.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					c.Close()
					return nil
				},
			})
		}),
	)
}
