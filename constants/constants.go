package constants

import "time"

const (
	LogLevelFlag                      = "log-level"
	HttpBindAddressServerFlag         = "http-bind-address-server"
	HttpBindAddressWorkerMessagesFlag = "http-bind-address-worker-messages"
	HttpBindAddressWorkerRetriesFlag  = "http-bind-address-worker-retries"

	RetriesScheduleFlag = "retries-schedule"
	RetriesCronFlag     = "retries-cron"

	StorageMongoConnStringFlag   = "storage-mongo-conn-string"
	StorageMongoDatabaseNameFlag = "storage-mongo-database-name"

	KafkaBrokersFlag       = "kafka-brokers"
	KafkaGroupIDFlag       = "kafka-consumer-group"
	KafkaTopicsFlag        = "kafka-topics"
	KafkaTLSEnabledFlag    = "kafka-tls-enabled"
	KafkaSASLEnabledFlag   = "kafka-sasl-enabled"
	KafkaSASLMechanismFlag = "kafka-sasl-mechanism"
	KafkaUsernameFlag      = "kafka-username"
	KafkaPasswordFlag      = "kafka-password"
)

const (
	DefaultBindAddressServer         = ":8080"
	DefaultBindAddressWorkerMessages = ":8081"
	DefaultBindAddressWorkerRetries  = ":8082"

	DefaultMongoConnString   = "mongodb://admin:admin@localhost:27017/"
	DefaultMongoDatabaseName = "webhooks"

	MongoCollectionConfigs  = "configs"
	MongoCollectionAttempts = "attempts"

	DefaultKafkaTopic   = "default"
	DefaultKafkaBroker  = "localhost:9092"
	DefaultKafkaGroupID = "webhooks"
)

var (
	DefaultRetriesSchedule = []time.Duration{time.Minute, 5 * time.Minute, 30 * time.Minute, 5 * time.Hour, 24 * time.Hour}
	DefaultRetriesCron     = time.Minute
)
