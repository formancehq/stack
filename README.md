# Webhooks

Webhooks is a service used to manage configs and to send webhooks to their endpoints.
A config is made of the following information:
- Endpoint: a single URL where messages are sent to.
- EventTypes: an array of string identifiers denoting the type of message being sent and are the primary way for webhook consumers to configure what events they are interested in receiving.
- Secret: a string used to verify received webhooks. Every webhook and its metadata is signed with a unique key for each endpoint. This signature can then be used to verify the webhook indeed comes from this service.
  The format is a random string of bytes of size 24, base64 encoded. (larger size after encoding)

Webhooks has two starting modes, split into two separate commands:

- `server`: REST web service API managing webhooks configs for users.
- `worker`: background service consuming kafka events on selected topics and sending webhooks for each configs.

## Run the stack locally
```
docker compose up
```

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
      --kafka-topics strings                 Kafka topics (default [default])
      --kafka-username string                Kafka username
      --log-level string                     Log level (default "info")
      --storage-mongo-conn-string string     Mongo connection string (default "mongodb://admin:admin@localhost:27017/")
      --storage-mongo-database-name string   Mongo database name (default "webhooks")

Use "webhooks [command] --help" for more information about a command. 
```
