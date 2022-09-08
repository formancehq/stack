package env

import (
	"fmt"
	"strings"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/go-libs/sharedlogging/sharedlogginglogrus"
	"github.com/numary/webhooks/constants"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Init(flagSet *pflag.FlagSet) error {
	flagSet.String(constants.LogLevelFlag, logrus.InfoLevel.String(), "Log level")

	flagSet.String(constants.HttpBindAddressServerFlag, constants.DefaultBindAddressServer, "server HTTP bind address")
	flagSet.String(constants.HttpBindAddressWorkerFlag, constants.DefaultBindAddressWorker, "worker HTTP bind address")
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
		return fmt.Errorf("viper.BinPFlags: %w", err)
	}

	LoadEnv(viper.GetViper())

	logger := logrus.New()
	lvl, err := logrus.ParseLevel(viper.GetString(constants.LogLevelFlag))
	if err != nil {
		return fmt.Errorf("logrus.ParseLevel: %w", err)
	}
	logger.SetLevel(lvl)

	sharedlogging.SetFactory(
		sharedlogging.StaticLoggerFactory(
			sharedlogginglogrus.New(logger)))

	return nil
}

func LoadEnv(v *viper.Viper) {
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
}
