module github.com/formancehq/wallets

go 1.19

require (
	github.com/formancehq/formance-sdk-go v1.0.0-rc.1.0.20230213140938-3b8735cf0168
	github.com/formancehq/stack/libs/go-libs v0.0.0-20230222164357-55840b21a337
	github.com/go-chi/chi/v5 v5.0.8
	github.com/go-chi/render v1.0.2
	github.com/google/uuid v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/riandyrn/otelchi v0.5.1
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.15.0
	github.com/stretchr/testify v1.8.1
	go.uber.org/fx v1.19.2
	golang.org/x/oauth2 v0.5.0
)

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/uptrace/opentelemetry-go-extra/otellogrus v0.1.21 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelutil v0.1.21 // indirect
	go.opentelemetry.io/contrib v1.14.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.39.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.14.0 // indirect
	go.opentelemetry.io/otel v1.13.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.13.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.13.0 // indirect
	go.opentelemetry.io/otel/metric v0.36.0 // indirect
	go.opentelemetry.io/otel/sdk v1.13.0 // indirect
	go.opentelemetry.io/otel/trace v1.13.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/dig v1.16.1 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230221151758-ace64dc21148 // indirect
	google.golang.org/grpc v1.53.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/formancehq/stack/libs/go-libs => ../../libs/go-libs

replace github.com/formancehq/formance-sdk-go => ../../sdks/go
