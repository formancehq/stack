# Webhooks

This service has two starting modes, for two responsibilities:

- server: RESTful web service API managing webhooks configs for users.
- worker: background service consuming kafka events on selected topics and sending webhooks for each matching event type in configs.

## Run linters and tests locally

```
task install:lint
task
```

Run the tests for a specific package:
```
task tests:local PKG=./pkg/model
```

Run a specific test (regexp):
```
task tests:local RUN=TestServer
```

## Build locally
```
task build:local
```

## Usage
```
$> ./webhooks
Usage:
  webhooks [command]

Available Commands:
  help        Help about any command
  server      Start webhooks server
  version     Get webhooks version
  worker      Start webhooks worker

Flags:
  -h, --help                                 help for webhooks
      --http-bind-address-server string      server HTTP bind address (default ":8080")
      --http-bind-address-worker string      worker HTTP bind address (default ":8081")
      --kafka-brokers strings                Kafka brokers (default [localhost:9092])
      --kafka-consumer-group string          Kafka consumer group (default "webhooks")
      --kafka-password string                Kafka password
      --kafka-sasl-enabled                   Kafka SASL enabled
      --kafka-sasl-mechanism string          Kafka SASL mechanism
      --kafka-tls-enabled                    Kafka TLS enabled
      --kafka-tls-insecure-skip-verify       Kafka TLS insecure skip verify
      --kafka-topics strings                 Kafka topics (default [default])
      --kafka-username string                Kafka username
      --log-level string                     Log level (default "info")
      --storage-mongo-conn-string string     Mongo connection string (default "mongodb://admin:admin@localhost:27017/")
      --storage-mongo-database-name string   Mongo database name (default "webhooks")
      --svix-app-id string                   Svix App ID
      --svix-app-name string                 Svix App Name
      --svix-server-url string               Svix Server URL
      --svix-token string                    Svix auth token

Use "webhooks [command] --help" for more information about a command. 
```
