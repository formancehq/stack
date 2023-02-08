package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/webhooks/cmd"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
)

var ErrMechanism = errors.New("unrecognized SASL mechanism")

func NewClient() (*kgo.Client, []string, error) {
	logging.Infof("connecting to new kafka client...")
	var opts []kgo.Opt

	brokers := viper.GetStringSlice(publish.PublisherKafkaBrokerFlag)
	opts = append(opts, kgo.SeedBrokers(brokers...))

	groupID := viper.GetString(cmd.ServiceName)
	opts = append(opts, kgo.ConsumerGroup(groupID))

	topics := viper.GetStringSlice(flag.KafkaTopics)
	opts = append(opts, kgo.ConsumeTopics(topics...))

	opts = append(opts, kgo.AllowAutoTopicCreation())

	kafkaClient, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, []string{}, fmt.Errorf("kgo.NewClient: %w", err)
	}

	healthy := false
	for !healthy {
		if err := kafkaClient.Ping(context.Background()); err != nil {
			logging.Infof("trying to reach broker: %s", err)
			time.Sleep(3 * time.Second)
		} else {
			healthy = true
		}
	}

	return kafkaClient, topics, nil
}
