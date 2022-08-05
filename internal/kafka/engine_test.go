package kafka

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks-cloud/cmd/constants"
	"github.com/numary/webhooks-cloud/internal/env"
	"github.com/numary/webhooks-cloud/internal/storage/mongo"
	"github.com/numary/webhooks-cloud/internal/svix"
	"github.com/numary/webhooks-cloud/pkg/model"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/spf13/pflag"
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

func TestEngine(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(), 15*time.Second)
	defer cancel()
	sharedlogging.GetLogger(ctx).Infof("started TestEngine")

	flagSet := pflag.NewFlagSet("TestEngine", pflag.ContinueOnError)
	require.NoError(t, env.Init(flagSet))

	store, err := mongo.NewConfigStore()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, store.Close(ctx))
	}()

	require.NoError(t, store.DropConfigsCollection(ctx))

	topic := os.Getenv("KAFKA_TOPIC")
	conn, err := kafkago.DialLeader(context.Background(),
		"tcp", constants.DefaultKafkaBroker, topic, 0)
	require.NoError(t, err)
	require.NoError(t, conn.SetWriteDeadline(time.Now().Add(10*time.Second)))

	eventType := "COMMITTED_TRANSACTIONS"
	i, err := conn.WriteMessages(newEventMessage(t, eventType))
	require.NoError(t, err)
	require.NotEqual(t, 0, i)

	endpoint := "https://example.com"
	cfg := model.Config{
		Active:     true,
		EventTypes: []string{eventType},
		Endpoints:  []string{endpoint},
	}
	require.NoError(t, cfg.Validate())

	_, err = store.InsertOneConfig(ctx, cfg)
	require.NoError(t, err)

	kcfg, err := NewKafkaReaderConfig()
	require.NoError(t, err)

	reader := kafkago.NewReader(kcfg)
	defer func(reader *kafkago.Reader) {
		require.NoError(t, reader.Close())
	}(reader)

	svixClient, svixAppId, err := svix.New()
	require.NoError(t, err)
	require.NoError(t, svix.CreateEndpoint(svixClient, svixAppId, endpoint))

	e := NewEngine(reader, store, svixClient, svixAppId)
	fetchedMsgs, sentWebhooks, err := e.Run(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, fetchedMsgs)
	assert.Equal(t, 1, sentWebhooks)

	require.NoError(t, svix.DeleteAllEndpoints(svixClient, svixAppId))
}
