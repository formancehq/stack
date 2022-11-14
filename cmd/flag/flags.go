package flag

import (
	"fmt"
	"strings"
	"time"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/go-libs/sharedlogging/sharedlogginglogrus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	LogLevel              = "log-level"
	HttpBindAddressServer = "http-bind-address-server"
	HttpBindAddressWorker = "http-bind-address-worker"

	RetrySchedule = "retry-schedule"
	RetryCron     = "retry-cron"

	StorageMongoConnString   = "storage-mongo-conn-string"
	StorageMongoDatabaseName = "storage-mongo-database-name"

	KafkaBrokers       = "kafka-brokers"
	KafkaGroupID       = "kafka-consumer-group"
	KafkaTopics        = "kafka-topics"
	KafkaTLSEnabled    = "kafka-tls-enabled"
	KafkaSASLEnabled   = "kafka-sasl-enabled"
	KafkaSASLMechanism = "kafka-sasl-mechanism"
	KafkaUsername      = "kafka-username"
	KafkaPassword      = "kafka-password"
)

const (
	DefaultBindAddressServer = ":8080"
	DefaultBindAddressWorker = ":8081"

	DefaultMongoConnString   = "mongodb://admin:admin@localhost:27017/"
	DefaultMongoDatabaseName = "webhooks"

	DefaultKafkaTopic   = "default"
	DefaultKafkaBroker  = "localhost:9092"
	DefaultKafkaGroupID = "webhooks"
)

var (
	DefaultRetrySchedule = []time.Duration{time.Minute, 5 * time.Minute, 30 * time.Minute, 5 * time.Hour, 24 * time.Hour}
	DefaultRetryCron     = time.Minute
)

var ErrScheduleInvalid = errors.New("the retry schedule should only contain durations of at least 1 second")

func Init(flagSet *pflag.FlagSet) (retrySchedule []time.Duration, err error) {
	flagSet.String(LogLevel, logrus.InfoLevel.String(), "Log level")

	flagSet.String(HttpBindAddressServer, DefaultBindAddressServer, "server HTTP bind address")
	flagSet.String(HttpBindAddressWorker, DefaultBindAddressWorker, "worker HTTP bind address")
	flagSet.DurationSlice(RetrySchedule, DefaultRetrySchedule, "worker retry schedule")
	flagSet.Duration(RetryCron, DefaultRetryCron, "worker retry cron")
	flagSet.String(StorageMongoConnString, DefaultMongoConnString, "Mongo connection string")
	flagSet.String(StorageMongoDatabaseName, DefaultMongoDatabaseName, "Mongo database name")

	flagSet.StringSlice(KafkaBrokers, []string{DefaultKafkaBroker}, "Kafka brokers")
	flagSet.String(KafkaGroupID, DefaultKafkaGroupID, "Kafka consumer group")
	flagSet.StringSlice(KafkaTopics, []string{DefaultKafkaTopic}, "Kafka topics")
	flagSet.Bool(KafkaTLSEnabled, false, "Kafka TLS enabled")
	flagSet.Bool(KafkaSASLEnabled, false, "Kafka SASL enabled")
	flagSet.String(KafkaSASLMechanism, "", "Kafka SASL mechanism")
	flagSet.String(KafkaUsername, "", "Kafka username")
	flagSet.String(KafkaPassword, "", "Kafka password")

	if err := viper.BindPFlags(flagSet); err != nil {
		return nil, fmt.Errorf("viper.BinPFlags: %w", err)
	}

	LoadEnv(viper.GetViper())

	logger := logrus.New()
	lvl, err := logrus.ParseLevel(viper.GetString(LogLevel))
	if err != nil {
		return nil, fmt.Errorf("logrus.ParseLevel: %w", err)
	}
	logger.SetLevel(lvl)
	if logger.GetLevel() < logrus.DebugLevel {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	retrySchedule, err = flagSet.GetDurationSlice(RetrySchedule)
	if err != nil {
		return nil, errors.Wrap(err, "flagSet.GetDurationSlice")
	}

	// Check that the schedule is valid
	for _, s := range retrySchedule {
		if s < time.Second {
			return nil, ErrScheduleInvalid
		}
	}

	sharedlogging.SetFactory(
		sharedlogging.StaticLoggerFactory(
			sharedlogginglogrus.New(logger)))

	return retrySchedule, nil
}

func LoadEnv(v *viper.Viper) {
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
}

func init() {
	LoadEnv(viper.GetViper())
}
