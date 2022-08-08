package kafka

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/numary/webhooks-cloud/cmd/constants"
	"github.com/numary/webhooks-cloud/internal/env"
	"github.com/numary/webhooks-cloud/internal/storage/mongo"
	"github.com/numary/webhooks-cloud/internal/svix"
	"github.com/numary/webhooks-cloud/pkg/model"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newEventMessage(t *testing.T, eventType string) kafkago.Message {
	ev := Event{
		Date: time.Now().UTC(),
		Type: eventType,
		Payload: map[string]any{
			"data": "test",
		},
	}

	by, err := json.Marshal(ev)
	require.NoError(t, err)

	return kafkago.Message{
		Key:   []byte("key"),
		Value: by,
	}
}

func TestWorker(t *testing.T) {
	flagSet := pflag.NewFlagSet("TestWorker", pflag.ContinueOnError)
	require.NoError(t, env.Init(flagSet))

	store, err := mongo.NewConfigStore()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, store.Close(context.Background()))
	}()

	conn, err := kafkago.DialLeader(context.Background(),
		"tcp",
		viper.GetStringSlice(constants.KafkaBrokersFlag)[0],
		viper.GetStringSlice(constants.KafkaTopicsFlag)[0], 0)
	require.NoError(t, err)

	eventType := "COMMITTED_TRANSACTIONS"
	nbBytes, err := conn.WriteMessages(newEventMessage(t, eventType))
	require.NoError(t, err)
	require.NotEqual(t, 0, nbBytes)

	endpoint := "https://example.com"
	cfg := model.Config{
		Endpoint:   endpoint,
		EventTypes: []string{eventType},
	}
	require.NoError(t, cfg.Validate())

	id, err := store.InsertOneConfig(context.Background(), cfg)
	require.NoError(t, err)

	kcfg, err := NewKafkaReaderConfig()
	require.NoError(t, err)

	reader := kafkago.NewReader(kcfg)
	defer func(reader *kafkago.Reader) {
		require.NoError(t, reader.Close())
	}(reader)

	svixClient, svixAppId, err := svix.New()
	require.NoError(t, err)
	require.NoError(t, svix.CreateEndpoint(id, cfg, svixClient, svixAppId))

	ctx, cancel := context.WithTimeout(
		context.Background(), 10*time.Second)
	defer cancel()
	fetchedMsgs, sentWebhooks, err := NewWorker(reader, store, svixClient, svixAppId).Run(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, fetchedMsgs)
	assert.Equal(t, 1, sentWebhooks)

	cur, err := store.FindAllConfigs(context.Background())
	require.NoError(t, err)
	spew.Dump(cur)

	deletedCount, err := store.DeleteOneConfig(context.Background(), id)
	assert.Equal(t, int64(1), deletedCount)
	require.NoError(t, err)

	require.NoError(t, svix.DeleteEndpoint(id, svixClient, svixAppId))
}
