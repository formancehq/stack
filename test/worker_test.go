package test_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/constants"
	"github.com/numary/webhooks/pkg/kafka"
	"github.com/numary/webhooks/pkg/model"
	"github.com/numary/webhooks/pkg/server"
	"github.com/numary/webhooks/pkg/worker"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/fx/fxtest"
)

func TestWorker(t *testing.T) {
	serverApp := fxtest.New(t,
		server.StartModule(
			httpClient, viper.GetString(constants.HttpBindAddressServerFlag)))
	workerApp := fxtest.New(t,
		worker.StartModule(
			httpClient, viper.GetString(constants.HttpBindAddressWorkerFlag)))

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

	resBody := requestServer(t, http.MethodGet, server.PathConfigs, http.StatusOK)
	cur := decodeCursorResponse[model.ConfigInserted](t, resBody)
	for _, cfg := range cur.Data {
		requestServer(t, http.MethodDelete, server.PathConfigs+"/"+cfg.ID, http.StatusOK)
	}
	require.NoError(t, resBody.Close())

	eventType := "TYPE"
	endpoint := "https://example.com"
	cfg := model.Config{
		Endpoint:   endpoint,
		Secret:     model.NewSecret(),
		EventTypes: []string{"OTHER_TYPE", eventType},
	}
	require.NoError(t, cfg.Validate())

	var insertedId string
	resBody = requestServer(t, http.MethodPost, server.PathConfigs, http.StatusOK, cfg)
	require.NoError(t, json.NewDecoder(resBody).Decode(&insertedId))
	require.NoError(t, resBody.Close())

	n := 3
	var messages []kafkago.Message
	for i := 0; i < n; i++ {
		messages = append(messages, newEventMessage(t, eventType, i))
	}
	nbBytes, err := conn.WriteMessages(messages...)
	require.NoError(t, err)
	require.NotEqual(t, 0, nbBytes)

	t.Run("health check", func(t *testing.T) {
		requestWorker(t, http.MethodGet, server.PathHealthCheck, http.StatusOK)
	})

	t.Run("messages", func(t *testing.T) {
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
		time.Sleep(time.Second)
		assert.Equal(t, n, msgs)
	})

	requestServer(t, http.MethodDelete, server.PathConfigs+"/"+insertedId, http.StatusOK)

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
			"id":   id,
			"type": eventType,
		},
	}

	by, err := json.Marshal(ev)
	require.NoError(t, err)

	return kafkago.Message{
		Key:   []byte("key"),
		Value: by,
	}
}
