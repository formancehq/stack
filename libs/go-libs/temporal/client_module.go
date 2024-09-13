package temporal

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/operatorservice/v1"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"
	"go.uber.org/fx"
)

type SearchAttributes struct {
	SearchAttributes map[string]enums.IndexedValueType
}

func NewModule(tracer trace.Tracer, searchAttributes SearchAttributes) fx.Option {
	address := viper.GetString(TemporalAddressFlag)
	namespace := viper.GetString(TemporalNamespaceFlag)
	certStr := viper.GetString(TemporalSSLClientCertFlag)
	key := viper.GetString(TemporalSSLClientKeyFlag)
	initSearchAttributes := viper.GetBool(TemporalInitSearchAttributes)

	return fx.Options(
		fx.Provide(func(logger logging.Logger) (client.Options, error) {

			var cert *tls.Certificate
			if key != "" && certStr != "" {
				clientCert, err := tls.X509KeyPair([]byte(certStr), []byte(key))
				if err != nil {
					return client.Options{}, err
				}
				cert = &clientCert
			}

			tracingInterceptor, err := opentelemetry.NewTracingInterceptor(opentelemetry.TracerOptions{
				Tracer: tracer,
			})
			if err != nil {
				return client.Options{}, err
			}

			options := client.Options{
				Namespace:    namespace,
				HostPort:     address,
				Interceptors: []interceptor.ClientInterceptor{tracingInterceptor},
				Logger:       newLogger(logger),
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
				OnStart: func(ctx context.Context) error {
					if initSearchAttributes {
						return createSearchAttributes(ctx, c, namespace, searchAttributes.SearchAttributes)
					}
					return nil
				},
				OnStop: func(ctx context.Context) error {
					c.Close()
					return nil
				},
			})
		}),
	)
}

func createSearchAttributes(ctx context.Context, c client.Client, namespace string, searchAttributes map[string]enums.IndexedValueType) error {
	_, err := c.OperatorService().AddSearchAttributes(logging.TestingContext(), &operatorservice.AddSearchAttributesRequest{
		SearchAttributes: searchAttributes,
		Namespace:        namespace,
	})
	if err != nil {
		if _, ok := err.(*serviceerror.AlreadyExists); !ok {
			return err
		}
	}
	// Search attributes are created asynchronously, so poll the list, until it is ready
	for {
		ret, err := c.OperatorService().ListSearchAttributes(ctx, &operatorservice.ListSearchAttributesRequest{
			Namespace: namespace,
		})
		if err != nil {
			panic(err)
		}

		done := true
		for key := range searchAttributes {
			if ret.CustomAttributes[key] == enums.INDEXED_VALUE_TYPE_UNSPECIFIED {
				done = false
				break
			}
		}

		if done {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
}
