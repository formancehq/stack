package env

import (
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/numary/webhooks-cloud/cmd/constants"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Init(flagSet *pflag.FlagSet) error {
	flagSet.String(constants.ServerHttpBindAddressFlag, constants.DefaultBindAddress, "API bind address")
	flagSet.String(constants.StorageMongoConnStringFlag, constants.DefaultMongoConnString, "Mongo connection string")

	flagSet.StringSlice(constants.KafkaBrokersFlag, []string{constants.DefaultKafkaBroker}, "Kafka brokers")
	flagSet.String(constants.KafkaGroupIDFlag, constants.DefaultKafkaGroupID, "Kafka consumer group")
	flagSet.String(constants.KafkaTopicFlag, constants.DefaultKafkaTopic, "Kafka topic")
	flagSet.Bool(constants.KafkaTLSEnabledFlag, false, "Kafka TLS enabled")
	flagSet.Bool(constants.KafkaTLSInsecureSkipVerifyFlag, false, "Kafka TLS insecure skip verify")
	flagSet.Bool(constants.KafkaSASLEnabledFlag, false, "Kafka SASL enabled")
	flagSet.String(constants.KafkaSASLMechanismFlag, "", "Kafka SASL mechanism")
	flagSet.String(constants.KafkaUsernameFlag, "", "Kafka username")
	flagSet.String(constants.KafkaPasswordFlag, "", "Kafka password")

	flagSet.String(constants.SvixTokenFlag, "", "Svix auth token")
	flagSet.String(constants.SvixAppIdFlag, "", "Svix app ID")

	if err := viper.BindPFlags(flagSet); err != nil {
		return err
	}

	LoadEnv(viper.GetViper())

	spew.Dump("env.Init", viper.AllSettings())
	return nil
}

func LoadEnv(v *viper.Viper) {
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
}
