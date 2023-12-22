# Webhooks

Webhooks is a service used to manage user configs and send webhooks to endpoints.
A user config is made of the following information:
- Endpoint: a single URL where messages are sent to.
- EventTypes: an array of string identifiers denoting the type of message being sent and are the primary way for webhook consumers to configure what events they are interested in receiving. Are stored in lower-case format.
- Secret: a string used to verify received webhooks. Every webhook and its metadata is signed with a unique key for each endpoint. This signature can then be used to verify the webhook indeed comes from this service.
  The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)

The service has 3 starting modes, split into 3 separate commands:

- `server`: REST web service API managing webhooks configs for users.
- `worker`: background service consuming kafka events on selected topics to send webhooks based on user configs and periodically finding failed webhooks requests to retry and sending new attempts.

## Run linters and tests 

Run the linters:
```
earthly +lint
```

Run the tests:
```
earthly -P +tests
```

## Usage
```
Usage:
  webhooks [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  serve       Run webhooks server
  version     Get webhooks version
  worker      Run webhooks worker

Flags:
      --abort-after duration                          consider a webhook as failed after retrying it for this duration. (default 720h0m0s)
      --debug                                         Debug mode
  -h, --help                                          help for webhooks
      --kafka-topics strings                          Kafka topics (default [default])
      --listen string                                 server HTTP bind address (default ":8080")
      --log-level string                              Log level (default "info")
      --max-backoff-delay duration                    maximum backoff delay (default 1h0m0s)
      --min-backoff-delay duration                    minimum backoff delay (default 1m0s)
      --otel-resource-attributes strings              Additional OTLP resource attributes
      --otel-service-name string                      OpenTelemetry service name
      --otel-traces                                   Enable OpenTelemetry traces support
      --otel-traces-batch                             Use OpenTelemetry batching
      --otel-traces-exporter string                   OpenTelemetry traces exporter (default "stdout")
      --otel-traces-exporter-jaeger-endpoint string   OpenTelemetry traces Jaeger exporter endpoint
      --otel-traces-exporter-jaeger-password string   OpenTelemetry traces Jaeger exporter password
      --otel-traces-exporter-jaeger-user string       OpenTelemetry traces Jaeger exporter user
      --otel-traces-exporter-otlp-endpoint string     OpenTelemetry traces grpc endpoint
      --otel-traces-exporter-otlp-insecure            OpenTelemetry traces grpc insecure
      --otel-traces-exporter-otlp-mode string         OpenTelemetry traces OTLP exporter mode (grpc|http) (default "grpc")
      --publisher-http-enabled                        Sent write event to http endpoint
      --publisher-kafka-broker strings                Kafka address is kafka enabled (default [localhost:9092])
      --publisher-kafka-enabled                       Publish write events to kafka
      --publisher-kafka-sasl-enabled                  Enable SASL authentication on kafka publisher
      --publisher-kafka-sasl-mechanism string         SASL authentication mechanism
      --publisher-kafka-sasl-password string          SASL password
      --publisher-kafka-sasl-scram-sha-size int       SASL SCRAM SHA size (default 512)
      --publisher-kafka-sasl-username string          SASL username
      --publisher-kafka-tls-enabled                   Enable TLS to connect on kafka
      --publisher-nats-client-id string               Nats client ID
      --publisher-nats-enabled                        Publish write events to nats
      --publisher-nats-max-reconnect int              Nats: set the maximum number of reconnect attempts. (default 30)
      --publisher-nats-reconnect-wait duration        Nats: the wait time between reconnect attempts. (default 2s)
      --publisher-nats-url string                     Nats url
      --publisher-topic-mapping strings               Define mapping between internal event types and topics
      --retry-period duration                         worker retry period (default 1m0s)
      --storage-postgres-conn-string string           Postgres connection string (default "postgresql://webhooks:webhooks@127.0.0.1/webhooks?sslmode=disable")
      --worker                                        Enable worker on server
```
