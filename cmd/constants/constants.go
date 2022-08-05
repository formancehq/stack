package constants

const (
	StorageMongoConnStringFlag = "storage-mongo-conn-string"
	ServerHttpBindAddressFlag  = "server-http-bind-address"

	KafkaBrokersFlag               = "kafka-brokers"
	KafkaGroupIDFlag               = "kafka-consumer-group"
	KafkaTopicFlag                 = "kafka-topic"
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
	DefaultMongoConnString = "mongodb://admin:admin@localhost:27017/"
	DefaultBindAddress     = ":8080"

	DefaultKafkaTopic   = "default"
	DefaultKafkaBroker  = "localhost:9092"
	DefaultKafkaGroupID = "webhooks"
)
