package kafka

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/constants"
	"github.com/numary/webhooks/internal/env"
	"github.com/numary/webhooks/internal/storage/mongo"
	"github.com/numary/webhooks/internal/svix"
	"github.com/numary/webhooks/pkg/model"
	"github.com/numary/webhooks/pkg/service"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorker_Start(t *testing.T) {
	flagSet := pflag.NewFlagSet("TestWorker", pflag.ContinueOnError)
	require.NoError(t, env.Init(flagSet))

	configStore, err := mongo.NewConfigStore()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, configStore.Close(context.Background()))
	}()

	eventType := "TYPE"
	endpoint := "https://example.com"
	cfg := model.Config{
		Endpoint:   endpoint,
		EventTypes: []string{eventType},
	}
	require.NoError(t, cfg.Validate())

	svixClient, svixAppId, err := svix.New()
	require.NoError(t, err)

	id, err := service.InsertOneConfig(cfg, context.Background(), configStore, svixClient, svixAppId)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, service.DeleteOneConfig(context.Background(), id, configStore, svixClient, svixAppId))
	}()

	kcfg, err := NewKafkaReaderConfig()
	require.NoError(t, err)
	kafkaReader := kafkago.NewReader(kcfg)
	defer func() {
		require.NoError(t, kafkaReader.Close())
	}()

	var conn *kafkago.Conn
	for conn == nil {
		conn, err = kafkago.DialLeader(context.Background(), "tcp",
			viper.GetStringSlice(constants.KafkaBrokersFlag)[0],
			viper.GetStringSlice(constants.KafkaTopicsFlag)[0], 0)
		if err != nil {
			sharedlogging.GetLogger(context.Background()).Debug("connecting to kafka: err: ", err)
			time.Sleep(3 * time.Second)
		}
	}
	defer func() {
		require.NoError(t, conn.Close())
	}()

	// TODO: delete when tests are working well
	t.Run("cleanup kafka", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for {
			m, err := kafkaReader.ReadMessage(ctx)
			if err != nil {
				sharedlogging.GetLogger(context.Background()).Debugf(
					"cleanup: kafka.Reader.FetchMessage: %s", err)
				break
			}
			sharedlogging.GetLogger(context.Background()).Debugf(
				"cleanup: new kafka message fetched: %s", string(m.Value))
		}
	})

	t.Run("no messages", func(t *testing.T) {
		w := NewWorker(kafkaReader, configStore, svixClient, svixAppId)
		go w.Start(context.Background())
		w.Stop(context.Background())

		event, ok := <-w.eventChan
		assert.Nil(t, event)
		assert.False(t, ok)

		message, ok := <-w.messageChan
		assert.Nil(t, message)
		assert.False(t, ok)

		err, ok = <-w.errChan
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("one message unknown event type", func(t *testing.T) {
		nbBytes, err := conn.WriteMessages(
			newEventMessage(t, "unknown", 0))
		require.NoError(t, err)
		require.NotEqual(t, 0, nbBytes)

		w := NewWorker(kafkaReader, configStore, svixClient, svixAppId)
		go w.Start(context.Background())

		event := <-w.eventChan
		assert.Equal(t, "unknown", event.Type)

		w.Stop(context.Background())

		var ok bool
		event, ok = <-w.eventChan
		assert.Nil(t, event)
		assert.False(t, ok)

		message, ok := <-w.messageChan
		assert.Nil(t, message)
		assert.False(t, ok)

		err, ok = <-w.errChan
		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("n messages", func(t *testing.T) {
		n := 3
		var messages []kafkago.Message
		for i := 0; i < n; i++ {
			messages = append(messages, newEventMessage(t, eventType, i))
		}
		nbBytes, err := conn.WriteMessages(messages...)
		require.NoError(t, err)
		require.NotEqual(t, 0, nbBytes)

		w := NewWorker(kafkaReader, configStore, svixClient, svixAppId)
		go w.Start(context.Background())

		receivedEvents, receivedMessages, receivedErrors := 0, 0, 0
		for i := 0; i < n*2; i++ {
			select {
			case <-w.eventChan:
				receivedEvents++
			case <-w.messageChan:
				receivedMessages++
			case <-w.errChan:
				receivedErrors++
			}
		}

		w.Stop(context.Background())

		assert.Equal(t, n, receivedEvents)
		assert.Equal(t, n, receivedMessages)
		assert.Equal(t, 0, receivedErrors)

		event, ok := <-w.eventChan
		assert.Nil(t, event)
		assert.False(t, ok)

		message, ok := <-w.messageChan
		assert.Nil(t, message)
		assert.False(t, ok)

		err, ok = <-w.errChan
		assert.NoError(t, err)
		assert.False(t, ok)
	})
}

func newEventMessage(t *testing.T, eventType string, id int) kafkago.Message {
	ev := Event{
		Date: time.Now().UTC(),
		Type: eventType,
		Payload: map[string]any{
			"id": id,
		},
	}

	by, err := json.Marshal(ev)
	require.NoError(t, err)

	return kafkago.Message{
		Key:   []byte("key"),
		Value: by,
	}
}
