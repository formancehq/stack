package constants

const (
	LogLevelFlag              = "log-level"
	ServerHttpBindAddressFlag = "server-http-bind-address"

	StorageMongoConnStringFlag   = "storage-mongo-conn-string"
	StorageMongoDatabaseNameFlag = "storage-mongo-database-name"

	KafkaBrokersFlag               = "kafka-brokers"
	KafkaGroupIDFlag               = "kafka-consumer-group"
	KafkaTopicsFlag                = "kafka-topics"
	KafkaTLSEnabledFlag            = "kafka-tls-enabled"
	KafkaTLSInsecureSkipVerifyFlag = "kafka-tls-insecure-skip-verify"
	KafkaSASLEnabledFlag           = "kafka-sasl-enabled"
	KafkaSASLMechanismFlag         = "kafka-sasl-mechanism"
	KafkaUsernameFlag              = "kafka-username"
	KafkaPasswordFlag              = "kafka-password"

	SvixTokenFlag = "svix-token"
	SvixAppIdFlag = "svix-app-id"
)

const (
	DefaultBindAddress = ":8080"

	DefaultMongoConnString   = "mongodb://admin:admin@localhost:27017/"
	DefaultMongoDatabaseName = "webhooks"

	DefaultKafkaTopic   = "default"
	DefaultKafkaBroker  = "localhost:9092"
	DefaultKafkaGroupID = "webhooks"
)
