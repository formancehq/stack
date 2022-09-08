package constants

const (
	LogLevelFlag              = "log-level"
	HttpBindAddressServerFlag = "http-bind-address-server"
	HttpBindAddressWorkerFlag = "http-bind-address-worker"

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
	DefaultBindAddressServer = ":8080"
	DefaultBindAddressWorker = ":8081"

	DefaultMongoConnString   = "mongodb://admin:admin@localhost:27017/"
	DefaultMongoDatabaseName = "webhooks"

	MongoCollectionConfigs  = "configs"
	MongoCollectionRequests = "requests"

	DefaultKafkaTopic   = "default"
	DefaultKafkaBroker  = "localhost:9092"
	DefaultKafkaGroupID = "webhooks"
)
