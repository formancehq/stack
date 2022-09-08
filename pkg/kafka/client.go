package kafka

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/constants"
	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

type Client interface {
	PollFetches(ctx context.Context) kgo.Fetches
	PauseFetchTopics(topics ...string) []string
	ResumeFetchTopics(topics ...string)
	Close()
}

var ErrMechanism = errors.New("unrecognized SASL mechanism")

func NewClient() (*kgo.Client, []string, error) {
	sharedlogging.Infof("connecting to new kafka client...")
	var opts []kgo.Opt
	if viper.GetBool(constants.KafkaTLSEnabledFlag) {
		opts = append(opts, kgo.DialTLSConfig(&tls.Config{
			MinVersion: tls.VersionTLS12,
		}))
	}

	if viper.GetBool(constants.KafkaSASLEnabledFlag) {
		a := scram.Auth{
			User: viper.GetString(constants.KafkaUsernameFlag),
			Pass: viper.GetString(constants.KafkaPasswordFlag),
		}
		switch mechanism := viper.GetString(constants.KafkaSASLMechanismFlag); mechanism {
		case "SCRAM-SHA-512":
			opts = append(opts, kgo.SASL(a.AsSha512Mechanism()))
		case "SCRAM-SHA-256":
			opts = append(opts, kgo.SASL(a.AsSha256Mechanism()))
		default:
			return nil, []string{}, ErrMechanism
		}
	}

	brokers := viper.GetStringSlice(constants.KafkaBrokersFlag)
	opts = append(opts, kgo.SeedBrokers(brokers...))

	groupID := viper.GetString(constants.KafkaGroupIDFlag)
	opts = append(opts, kgo.ConsumerGroup(groupID))

	topics := viper.GetStringSlice(constants.KafkaTopicsFlag)
	opts = append(opts, kgo.ConsumeTopics(topics...))

	opts = append(opts, kgo.AllowAutoTopicCreation())

	kafkaClient, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, []string{}, fmt.Errorf("kgo.NewClient: %w", err)
	}

	healthy := false
	for !healthy {
		if err := kafkaClient.Ping(context.Background()); err != nil {
			sharedlogging.Infof("trying to reach broker: %s", err)
			time.Sleep(3 * time.Second)
		} else {
			healthy = true
		}
	}

	return kafkaClient, topics, nil
}
