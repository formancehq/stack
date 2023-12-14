---
title: Configuration variables
hide_table_of_contents: true
---

:::info
When using environment vars, the var name must be prefixed with `NUMARY_`.
As an example, DEBUG could either be passed as `numary server start --debug` or `NUMARY_DEBUG=true numary server start`.
:::


 |Flag                                   |Env var                              |Default value                    |Description                                                                |
 |-                                      |-                                    |-                                |-                                                                          |
 |--auth-basic-credentials               |AUTH_BASIC_CREDENTIALS               |[]                               |HTTP basic auth credentials (`<username>:<password>`)                        |
 |--auth-basic-enabled                   |AUTH_BASIC_ENABLED                   |false                            |Enable basic auth                                                          |
 |--auth-bearer-audience                 |AUTH_BEARER_AUDIENCE                 |[]                               |Allowed audiences                                                          |
 |--auth-bearer-audiences-wildcard       |AUTH_BEARER_AUDIENCES_WILDCARD       |false                            |Don't check audience                                                       |
 |--auth-bearer-enabled                  |AUTH_BEARER_ENABLED                  |false                            |Enable bearer auth                                                         |
 |--auth-bearer-introspect-url           |AUTH_BEARER_INTROSPECT_URL           |                                 |OAuth2 introspect URL                                                      |
 |--auth-bearer-use-scopes               |AUTH_BEARER_USE_SCOPES               |false                            |Use scopes as defined by rfc https://datatracker.ietf.org/doc/html/rfc8693 |
 |--commit-policy                        |COMMIT_POLICY                        |                                 |Transaction commit policy (default or allow-past-timestamps)               |
 |--debug                                |DEBUG                                |false                            |Debug mode                                                                 |
 |--lock-strategy                        |LOCK_STRATEGY                        |memory                           |Lock strategy (memory, none, redis)                                        |
 |--lock-strategy-redis-duration         |LOCK_STRATEGY_REDIS_DURATION         |1m0s                             |Lock duration                                                              |
 |--lock-strategy-redis-retry            |LOCK_STRATEGY_REDIS_RETRY            |1s                               |Retry lock period                                                          |
 |--lock-strategy-redis-tls-enabled      |LOCK_STRATEGY_REDIS_TLS_ENABLED      |false                            |Use tls on redis                                                           |
 |--lock-strategy-redis-tls-insecure     |LOCK_STRATEGY_REDIS_TLS_INSECURE     |false                            |Whether redis is trusted or not                                            |
 |--lock-strategy-redis-url              |LOCK_STRATEGY_REDIS_URL              |                                 |Redis url when using redis locking strategy                                |
 |--otel-metrics                         |OTEL_METRICS                         |false                            |Enable OpenTelemetry metrics support                                       |
 |--otel-metrics-exporter                |OTEL_METRICS_EXPORTER                |stdout                           |OpenTelemetry metrics exporter                                             |
 |--otel-metrics-exporter-otlp-endpoint  |OTEL_METRICS_EXPORTER_OTLP_ENDPOINT  |                                 |OpenTelemetry metrics grpc endpoint                                        |
 |--otel-metrics-exporter-otlp-insecure  |OTEL_METRICS_EXPORTER_OTLP_INSECURE  |false                            |OpenTelemetry metrics grpc insecure                                        |
 |--otel-metrics-exporter-otlp-mode      |OTEL_METRICS_EXPORTER_OTLP_MODE      |grpc                             |OpenTelemetry metrics OTLP exporter mode (grpc|http)                       |
 |--otel-traces                          |OTEL_TRACES                          |false                            |Enable OpenTelemetry traces support                                        |
 |--otel-traces-batch                    |OTEL_TRACES_BATCH                    |false                            |Use OpenTelemetry batching                                                 |
 |--otel-traces-exporter                 |OTEL_TRACES_EXPORTER                 |stdout                           |OpenTelemetry traces exporter                                              |
 |--otel-traces-exporter-jaeger-endpoint |OTEL_TRACES_EXPORTER_JAEGER_ENDPOINT |                                 |OpenTelemetry traces Jaeger exporter endpoint                              |
 |--otel-traces-exporter-jaeger-password |OTEL_TRACES_EXPORTER_JAEGER_PASSWORD |                                 |OpenTelemetry traces Jaeger exporter password                              |
 |--otel-traces-exporter-jaeger-user     |OTEL_TRACES_EXPORTER_JAEGER_USER     |                                 |OpenTelemetry traces Jaeger exporter user                                  |
 |--otel-traces-exporter-otlp-endpoint   |OTEL_TRACES_EXPORTER_OTLP_ENDPOINT   |                                 |OpenTelemetry traces grpc endpoint                                         |
 |--otel-traces-exporter-otlp-insecure   |OTEL_TRACES_EXPORTER_OTLP_INSECURE   |false                            |OpenTelemetry traces grpc insecure                                         |
 |--otel-traces-exporter-otlp-mode       |OTEL_TRACES_EXPORTER_OTLP_MODE       |grpc                             |OpenTelemetry traces OTLP exporter mode (grpc|http)                        |
 |--publisher-http-enabled               |PUBLISHER_HTTP_ENABLED               |false                            |Sent write event to http endpoint                                          |
 |--publisher-kafka-broker               |PUBLISHER_KAFKA_BROKER               |[]                               |Kafka address is kafka enabled                                             |
 |--publisher-kafka-enabled              |PUBLISHER_KAFKA_ENABLED              |false                            |Publish write events to kafka                                              |
 |--publisher-kafka-sasl-enabled         |PUBLISHER_KAFKA_SASL_ENABLED         |false                            |Enable SASL authentication on kafka publisher                              |
 |--publisher-kafka-sasl-mechanism       |PUBLISHER_KAFKA_SASL_MECHANISM       |                                 |SASL authentication mechanism                                              |
 |--publisher-kafka-sasl-password        |PUBLISHER_KAFKA_SASL_PASSWORD        |                                 |SASL password                                                              |
 |--publisher-kafka-sasl-scram-sha-size  |PUBLISHER_KAFKA_SASL_SCRAM_SHA_SIZE  |512                              |SASL SCRAM SHA size                                                        |
 |--publisher-kafka-sasl-username        |PUBLISHER_KAFKA_SASL_USERNAME        |                                 |SASL username                                                              |
 |--publisher-kafka-tls-enabled          |PUBLISHER_KAFKA_TLS_ENABLED          |false                            |Enable TLS to connect on kafka                                             |
 |--publisher-topic-mapping              |PUBLISHER_TOPIC_MAPPING              |[]                               |Define mapping between internal event types and topics                     |
 |--segment-application-id               |SEGMENT_APPLICATION_ID               |                                 |Segment application id                                                     |
 |--segment-enabled                      |SEGMENT_ENABLED                      |true                             |Is segment enabled                                                         |
 |--segment-heartbeat-interval           |SEGMENT_HEARTBEAT_INTERVAL           |24h0m0s                          |Segment heartbeat interval                                                 |
 |--segment-write-key                    |SEGMENT_WRITE_KEY                    |lAVEcNA5tKkhkQGp2CvTBSsbGqFsbCIF |Segment write key                                                          |
 |--server.http.basic_auth               |SERVER_HTTP_BASIC_AUTH               |                                 |Http basic auth                                                            |
 |--server.http.bind_address             |SERVER_HTTP_BIND_ADDRESS             |localhost:3068                   |API bind address                                                           |
 |--storage.cache                        |STORAGE_CACHE                        |true                             |Storage cache                                                              |
 |--storage.dir                          |STORAGE_DIR                          |/Users/clement/.numary/data      |Storage directory (for sqlite)                                             |
 |--storage.driver                       |STORAGE_DRIVER                       |sqlite                           |Storage driver                                                             |
 |--storage-postgres-conn-string         |STORAGE_POSTGRES_CONN_STRING         |postgresql://localhost/postgres  |Postgre connection string                                                  |
 |--storage.sqlite.db_name               |STORAGE_SQLITE_DB_NAME               |numary                           |SQLite database name                                                       |
 |--ui.http.bind_address                 |UI_HTTP_BIND_ADDRESS                 |localhost:3068                   |UI bind address                                                            |
