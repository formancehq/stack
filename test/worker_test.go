package test_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/constants"
	"github.com/numary/webhooks/internal/kafka"
	"github.com/numary/webhooks/internal/model"
	"github.com/numary/webhooks/internal/server"
	"github.com/numary/webhooks/internal/worker"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/fx/fxtest"
)

func TestWorker(t *testing.T) {
	serverApp := fxtest.New(t, server.StartModule(context.Background(), httpClient))
	workerApp := fxtest.New(t, worker.StartModule(context.Background(), httpClient))

	require.NoError(t, serverApp.Start(context.Background()))
	require.NoError(t, workerApp.Start(context.Background()))

	require.NoError(t, mongoClient.Database(
		viper.GetString(constants.StorageMongoDatabaseNameFlag)).
		Collection("messages").Drop(context.Background()))

	var err error
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

	t.Run("health check", func(t *testing.T) {
		resp, err := http.Get(serverBaseURL + server.PathHealthCheck)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("clean existing configs", func(t *testing.T) {
		resp, err := http.Get(serverBaseURL + server.PathConfigs)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur := decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		for _, cfg := range cur.Data {
			req, err := http.NewRequest(http.MethodDelete, serverBaseURL+server.PathConfigs+"/"+cfg.ID, nil)
			require.NoError(t, err)
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})

	eventType := "TYPE_TO_SEND"
	endpoint := "https://example.com"
	cfg := model.Config{
		Endpoint:   endpoint,
		Secret:     model.NewSecret(),
		EventTypes: []string{"OTHER_TYPE", eventType},
	}
	require.NoError(t, cfg.Validate())

	var insertedId string

	t.Run("POST "+server.PathConfigs, func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, serverBaseURL+server.PathConfigs, buffer(t, cfg))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&insertedId))
	})

	n := 3
	var messages []kafkago.Message
	for i := 0; i < n; i++ {
		messages = append(messages, newEventMessage(t, eventType, i))
	}
	messages = append(messages, newEventMessage(t, "TYPE_NOT_TO_SEND", n))
	nbBytes, err := conn.WriteMessages(messages...)
	require.NoError(t, err)
	require.NotEqual(t, 0, nbBytes)

	msgs := 0
	for msgs != n {
		cur, err := mongoClient.Database(
			viper.GetString(constants.StorageMongoDatabaseNameFlag)).
			Collection("messages").Find(context.Background(), bson.M{}, nil)
		require.NoError(t, err)
		var results []message
		require.NoError(t, cur.All(context.Background(), &results))
		msgs = len(results)
		time.Sleep(time.Second)
	}

	assert.Equal(t, n, msgs)

	t.Run("DELETE "+server.PathConfigs, func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, serverBaseURL+server.PathConfigs+"/"+insertedId, nil)
		require.NoError(t, err)
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	require.NoError(t, mongoClient.Database(
		viper.GetString(constants.StorageMongoDatabaseNameFlag)).
		Collection("messages").Drop(context.Background()))

	require.NoError(t, serverApp.Stop(context.Background()))
	require.NoError(t, workerApp.Stop(context.Background()))
}

func newEventMessage(t *testing.T, eventType string, id int) kafkago.Message {
	ev := kafka.Event{
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
