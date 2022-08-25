package worker

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/numary/webhooks/constants"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/spf13/viper"
)

var ErrMechanism = errors.New("unrecognized SASL mechanism")

func NewKafkaReaderConfig() (kafkago.ReaderConfig, error) {
	dialer := kafkago.DefaultDialer
	if viper.GetBool(constants.KafkaTLSEnabledFlag) {
		dialer.TLS = &tls.Config{
			MinVersion: tls.VersionTLS12,
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
			return kafkago.ReaderConfig{}, ErrMechanism
		}
		mechanism, err := scram.Mechanism(alg,
			viper.GetString(constants.KafkaUsernameFlag), viper.GetString(constants.KafkaPasswordFlag))
		if err != nil {
			return kafkago.ReaderConfig{}, fmt.Errorf("scram.Mechanism: %w", err)
		}
		dialer.SASLMechanism = mechanism
	}

	brokers := viper.GetStringSlice(constants.KafkaBrokersFlag)
	if len(brokers) == 0 {
		brokers = []string{constants.DefaultKafkaBroker}
	}

	topics := viper.GetStringSlice(constants.KafkaTopicsFlag)
	if len(topics) == 0 {
		topics = []string{constants.DefaultKafkaTopic}
	}

	groupID := viper.GetString(constants.KafkaGroupIDFlag)
	if groupID == "" {
		groupID = constants.DefaultKafkaGroupID
	}

	return kafkago.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		GroupTopics: topics,
		MinBytes:    1,
		MaxBytes:    10e5,
		Dialer:      dialer,
	}, nil
}
