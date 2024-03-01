package temporalclient

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/formancehq/orchestration/internal/triggers"
	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/operatorservice/v1"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"
	"go.uber.org/fx"
)

func NewModule(address, namespace string, certStr string, key string) fx.Option {
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
				Tracer: workflow.Tracer,
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
					return createSearchAttributes(ctx, c, namespace)
				},
				OnStop: func(ctx context.Context) error {
					c.Close()
					return nil
				},
			})
		}),
	)
}

func createSearchAttributes(ctx context.Context, c client.Client, namespace string) error {
	_, err := c.OperatorService().AddSearchAttributes(logging.TestingContext(), &operatorservice.AddSearchAttributesRequest{
		SearchAttributes: map[string]enums.IndexedValueType{
			workflow.SearchAttributeWorkflowID: enums.INDEXED_VALUE_TYPE_TEXT,
			triggers.SearchAttributeTriggerID:  enums.INDEXED_VALUE_TYPE_TEXT,
		},
		Namespace: namespace,
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

		if ret.CustomAttributes[workflow.SearchAttributeWorkflowID] != enums.INDEXED_VALUE_TYPE_UNSPECIFIED &&
			ret.CustomAttributes[triggers.SearchAttributeTriggerID] != enums.INDEXED_VALUE_TYPE_UNSPECIFIED {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
	}
}
