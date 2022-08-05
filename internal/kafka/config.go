package kafka

import (
	"crypto/tls"
	"fmt"

	"github.com/numary/webhooks-cloud/cmd/constants"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/spf13/viper"
)

func NewKafkaReaderConfig() (kafka.ReaderConfig, error) {
	dialer := kafka.DefaultDialer
	if viper.GetBool(constants.KafkaTLSEnabledFlag) {
		dialer.TLS = &tls.Config{
			InsecureSkipVerify: viper.GetBool(constants.KafkaTLSInsecureSkipVerifyFlag),
		}
	}

	if viper.GetBool(constants.KafkaSASLEnabledFlag) {
		var alg scram.Algorithm
		switch mechanism := viper.GetString(constants.KafkaSASLMechanismFlag); mechanism {
		case "SCRAM-SHA-512":
			alg = scram.SHA512
		case "SCRAM-SHA-256":
			alg = scram.SHA256
		default:
			return kafka.ReaderConfig{}, fmt.Errorf("unrecognized SASL mechanism: %s", mechanism)
		}
		mechanism, err := scram.Mechanism(alg,
			viper.GetString(constants.KafkaUsernameFlag), viper.GetString(constants.KafkaPasswordFlag))
		if err != nil {
			return kafka.ReaderConfig{}, err
		}
		dialer.SASLMechanism = mechanism
	}

	brokers := viper.GetStringSlice(constants.KafkaBrokersFlag)
	if len(brokers) == 0 {
		brokers = []string{constants.DefaultKafkaBroker}
	}

	topic := viper.GetString(constants.KafkaTopicFlag)
	if topic == "" {
		topic = constants.DefaultKafkaTopic
	}

	groupID := viper.GetString(constants.KafkaGroupIDFlag)
	if groupID == "" {
		groupID = constants.DefaultKafkaGroupID
	}

	return kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1,
		MaxBytes: 10e5,
		Dialer:   dialer,
	}, nil
}
