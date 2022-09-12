package env

import (
	"fmt"
	"strings"
	"time"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/go-libs/sharedlogging/sharedlogginglogrus"
	"github.com/numary/webhooks/constants"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var ErrScheduleInvalid = errors.New("the retries schedule should only contain durations of at least 1 second")

func Init(flagSet *pflag.FlagSet) (retriesSchedule []time.Duration, err error) {
	flagSet.String(constants.LogLevelFlag, logrus.InfoLevel.String(), "Log level")

	flagSet.String(constants.HttpBindAddressServerFlag, constants.DefaultBindAddressServer, "server HTTP bind address")
	flagSet.String(constants.HttpBindAddressWorkerMessagesFlag, constants.DefaultBindAddressWorkerMessages, "worker messages HTTP bind address")
	flagSet.String(constants.HttpBindAddressWorkerRetriesFlag, constants.DefaultBindAddressWorkerRetries, "worker retries HTTP bind address")
	flagSet.DurationSlice(constants.RetriesScheduleFlag, constants.DefaultRetriesSchedule, "worker retries schedule")
	flagSet.Duration(constants.RetriesCronFlag, constants.DefaultRetriesCron, "worker retries cron")
	flagSet.String(constants.StorageMongoConnStringFlag, constants.DefaultMongoConnString, "Mongo connection string")
	flagSet.String(constants.StorageMongoDatabaseNameFlag, constants.DefaultMongoDatabaseName, "Mongo database name")

	flagSet.StringSlice(constants.KafkaBrokersFlag, []string{constants.DefaultKafkaBroker}, "Kafka brokers")
	flagSet.String(constants.KafkaGroupIDFlag, constants.DefaultKafkaGroupID, "Kafka consumer group")
	flagSet.StringSlice(constants.KafkaTopicsFlag, []string{constants.DefaultKafkaTopic}, "Kafka topics")
	flagSet.Bool(constants.KafkaTLSEnabledFlag, false, "Kafka TLS enabled")
	flagSet.Bool(constants.KafkaSASLEnabledFlag, false, "Kafka SASL enabled")
	flagSet.String(constants.KafkaSASLMechanismFlag, "", "Kafka SASL mechanism")
	flagSet.String(constants.KafkaUsernameFlag, "", "Kafka username")
	flagSet.String(constants.KafkaPasswordFlag, "", "Kafka password")

	if err := viper.BindPFlags(flagSet); err != nil {
		return nil, fmt.Errorf("viper.BinPFlags: %w", err)
	}

	LoadEnv(viper.GetViper())

	logger := logrus.New()
	lvl, err := logrus.ParseLevel(viper.GetString(constants.LogLevelFlag))
	if err != nil {
		return nil, fmt.Errorf("logrus.ParseLevel: %w", err)
	}
	logger.SetLevel(lvl)

	retriesSchedule, err = flagSet.GetDurationSlice(constants.RetriesScheduleFlag)
	if err != nil {
		return nil, errors.Wrap(err, "flagSet.GetDurationSlice")
	}

	// Check that the schedule is valid
	for _, s := range retriesSchedule {
		if s < time.Second {
			return nil, ErrScheduleInvalid
		}
	}

	sharedlogging.SetFactory(
		sharedlogging.StaticLoggerFactory(
			sharedlogginglogrus.New(logger)))

	return retriesSchedule, nil
}

func LoadEnv(v *viper.Viper) {
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
}
