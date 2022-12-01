package test_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/formancehq/webhooks/cmd/flag"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/kafka"
	"github.com/formancehq/webhooks/pkg/server"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/formancehq/webhooks/pkg/worker/messages"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestWorkerMessages(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx,
		options.Client().ApplyURI(
			viper.GetString(flag.StorageMongoConnString)))
	require.NoError(t, err)

	// Cleanup collections
	require.NoError(t, mongoClient.Database(
		viper.GetString(flag.StorageMongoDatabaseName)).
		Collection(storage.CollectionConfigs).Drop(context.Background()))

	// New test server with success handler
	httpServerSuccess := httptest.NewServer(http.HandlerFunc(webhooksSuccessHandler))
	defer func() {
		httpServerSuccess.CloseClientConnections()
		httpServerSuccess.Close()
	}()

	// New test server with fail handler
	httpServerFail := httptest.NewServer(http.HandlerFunc(webhooksFailHandler))
	defer func() {
		httpServerFail.CloseClientConnections()
		httpServerFail.Close()
	}()

	serverApp := fxtest.New(t,
		fx.Supply(httpServerSuccess.Client()),
		server.StartModule(
			viper.GetString(flag.HttpBindAddressServer)))
	require.NoError(t, serverApp.Start(context.Background()))

	cfgSuccess := webhooks.ConfigUser{
		Endpoint:   httpServerSuccess.URL,
		Secret:     secret,
		EventTypes: []string{"unknown", fmt.Sprintf("%s.%s", app1, type1)},
	}
	require.NoError(t, cfgSuccess.Validate())

	cfgFail := webhooks.ConfigUser{
		Endpoint:   httpServerFail.URL,
		Secret:     secret,
		EventTypes: []string{"unknown", fmt.Sprintf("%s.%s", app2, type2)},
	}
	require.NoError(t, cfgFail.Validate())

	requestServer(t, http.MethodPost, server.PathConfigs, http.StatusOK, cfgSuccess)
	requestServer(t, http.MethodPost, server.PathConfigs, http.StatusOK, cfgFail)

	t.Run("success", func(t *testing.T) {
		require.NoError(t, mongoClient.Database(
			viper.GetString(flag.StorageMongoDatabaseName)).
			Collection(storage.CollectionAttempts).Drop(context.Background()))

		retriesSchedule = []time.Duration{time.Second}
		viper.Set(flag.RetriesSchedule, retriesSchedule)

		workerMessagesApp := fxtest.New(t,
			fx.Supply(httpServerSuccess.Client()),
			messages.StartModule(
				viper.GetString(flag.HttpBindAddressWorkerMessages),
				retriesSchedule,
			))
		require.NoError(t, workerMessagesApp.Start(context.Background()))

		t.Run("health check", func(t *testing.T) {
			requestWorkerMessages(t, http.MethodGet, server.PathHealthCheck, http.StatusOK)
		})

		expectedSentWebhooks := 1
		kafkaClient, kafkaTopics, err := kafka.NewClient()
		require.NoError(t, err)

		by1, err := json.Marshal(event1)
		require.NoError(t, err)
		by3, err := json.Marshal(event3)
		require.NoError(t, err)

		records := []*kgo.Record{
			{Topic: kafkaTopics[0], Value: by1},
			{Topic: kafkaTopics[0], Value: by3},
		}
		if err := kafkaClient.ProduceSync(context.Background(), records...).FirstErr(); err != nil {
			fmt.Printf("record had a produce error while synchronously producing: %v\n", err)
		}
		kafkaClient.Close()

		t.Run("webhooks", func(t *testing.T) {
			msgs := 0
			for msgs != expectedSentWebhooks {
				cur, err := mongoClient.Database(
					viper.GetString(flag.StorageMongoDatabaseName)).
					Collection(storage.CollectionAttempts).
					Find(context.Background(), bson.M{}, nil)
				require.NoError(t, err)
				var results []webhooks.Attempt
				require.NoError(t, cur.All(context.Background(), &results))
				msgs = len(results)
				if msgs != expectedSentWebhooks {
					time.Sleep(time.Second)
				} else {
					for _, res := range results {
						require.Equal(t, webhooks.StatusAttemptSuccess, res.Status)
						require.Equal(t, 0, res.RetryAttempt)
					}
				}
			}
			time.Sleep(time.Second)
			require.Equal(t, expectedSentWebhooks, msgs)
		})

		require.NoError(t, workerMessagesApp.Stop(context.Background()))
	})

	t.Run("failure", func(t *testing.T) {
		require.NoError(t, mongoClient.Database(
			viper.GetString(flag.StorageMongoDatabaseName)).
			Collection(storage.CollectionAttempts).Drop(context.Background()))

		retriesSchedule = []time.Duration{time.Second}
		viper.Set(flag.RetriesSchedule, retriesSchedule)

		workerMessagesApp := fxtest.New(t,
			fx.Supply(httpServerFail.Client()),
			messages.StartModule(
				viper.GetString(flag.HttpBindAddressWorkerMessages),
				retriesSchedule,
			))
		require.NoError(t, workerMessagesApp.Start(context.Background()))

		t.Run("health check", func(t *testing.T) {
			requestWorkerMessages(t, http.MethodGet, server.PathHealthCheck, http.StatusOK)
		})

		expectedSentWebhooks := 1
		kafkaClient, kafkaTopics, err := kafka.NewClient()
		require.NoError(t, err)

		by2, err := json.Marshal(event2)
		require.NoError(t, err)
		by3, err := json.Marshal(event3)
		require.NoError(t, err)

		records := []*kgo.Record{
			{Topic: kafkaTopics[0], Value: by2},
			{Topic: kafkaTopics[0], Value: by3},
		}
		if err := kafkaClient.ProduceSync(context.Background(), records...).FirstErr(); err != nil {
			fmt.Printf("record had a produce error while synchronously producing: %v\n", err)
		}
		kafkaClient.Close()

		t.Run("webhooks", func(t *testing.T) {
			msgs := 0
			for msgs != expectedSentWebhooks {
				cur, err := mongoClient.Database(
					viper.GetString(flag.StorageMongoDatabaseName)).
					Collection(storage.CollectionAttempts).
					Find(context.Background(), bson.M{}, nil)
				require.NoError(t, err)
				var results []webhooks.Attempt
				require.NoError(t, cur.All(context.Background(), &results))
				msgs = len(results)
				if msgs != expectedSentWebhooks {
					time.Sleep(time.Second)
				} else {
					for _, res := range results {
						require.Equal(t, webhooks.StatusAttemptToRetry, res.Status)
						require.Equal(t, 0, res.RetryAttempt)
					}
				}
			}
			time.Sleep(time.Second)
			require.Equal(t, expectedSentWebhooks, msgs)
		})

		require.NoError(t, workerMessagesApp.Stop(context.Background()))
	})

	require.NoError(t, serverApp.Stop(context.Background()))
}
