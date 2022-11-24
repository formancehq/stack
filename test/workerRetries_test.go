package test_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/formancehq/webhooks/cmd/flag"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/server"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/formancehq/webhooks/pkg/worker/retries"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestWorkerRetries(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx,
		options.Client().ApplyURI(
			viper.GetString(flag.StorageMongoConnString)))
	require.NoError(t, err)

	t.Run("1 attempt to retry with success", func(t *testing.T) {
		require.NoError(t, mongoClient.Database(
			viper.GetString(flag.StorageMongoDatabaseName)).
			Collection(storage.CollectionAttempts).Drop(context.Background()))

		// New test server with success handler
		httpServerSuccess := httptest.NewServer(http.HandlerFunc(webhooksSuccessHandler))
		defer func() {
			httpServerSuccess.CloseClientConnections()
			httpServerSuccess.Close()
		}()

		failedAttempt := webhooks.Attempt{
			Date:      time.Now().UTC(),
			WebhookID: uuid.NewString(),
			Config: webhooks.Config{
				ConfigUser: webhooks.ConfigUser{
					Endpoint:   httpServerSuccess.URL,
					Secret:     secret,
					EventTypes: []string{type1},
				},
				ID:        uuid.NewString(),
				Active:    true,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
			Payload:        fmt.Sprintf("{\"type\":\"%s\"}", type1),
			StatusCode:     http.StatusNotFound,
			Status:         webhooks.StatusAttemptToRetry,
			RetryAttempt:   0,
			NextRetryAfter: time.Now().UTC(),
		}

		_, err = mongoClient.Database(
			viper.GetString(flag.StorageMongoDatabaseName)).
			Collection(storage.CollectionAttempts).InsertOne(context.Background(), failedAttempt)
		require.NoError(t, err)

		retriesSchedule = []time.Duration{time.Second, time.Second, time.Second}
		viper.Set(flag.RetriesSchedule, retriesSchedule)

		workerRetriesApp := fxtest.New(t,
			fx.Supply(httpServerSuccess.Client()),
			retries.StartModule(
				viper.GetString(flag.HttpBindAddressWorkerRetries),
				viper.GetDuration(flag.RetriesCron),
				retriesSchedule))
		require.NoError(t, workerRetriesApp.Start(context.Background()))

		requestWorkerRetries(t, http.MethodGet, server.PathHealthCheck, http.StatusOK)

		expectedAttempts := 2

		attempts := 0
		for attempts != expectedAttempts {
			opts := options.Find().SetSort(bson.M{webhooks.KeyID: -1})
			cur, err := mongoClient.Database(
				viper.GetString(flag.StorageMongoDatabaseName)).
				Collection(storage.CollectionAttempts).
				Find(context.Background(), bson.M{}, opts)
			require.NoError(t, err)
			var results []webhooks.Attempt
			require.NoError(t, cur.All(context.Background(), &results))
			attempts = len(results)
			if attempts != expectedAttempts {
				time.Sleep(time.Second)
			} else {
				// First attempt should be successful
				require.Equal(t, webhooks.StatusAttemptSuccess, results[0].Status)
				require.Equal(t, expectedAttempts-1, results[0].RetryAttempt)
			}
		}
		time.Sleep(time.Second)
		require.Equal(t, expectedAttempts, attempts)

		require.NoError(t, workerRetriesApp.Stop(context.Background()))
	})

	t.Run("retrying an attempt until failed at the end of the schedule", func(t *testing.T) {
		require.NoError(t, mongoClient.Database(
			viper.GetString(flag.StorageMongoDatabaseName)).
			Collection(storage.CollectionAttempts).Drop(context.Background()))

		// New test server with fail handler
		httpServerFail := httptest.NewServer(http.HandlerFunc(webhooksFailHandler))
		defer func() {
			httpServerFail.CloseClientConnections()
			httpServerFail.Close()
		}()

		failedAttempt := webhooks.Attempt{
			Date:      time.Now().UTC(),
			WebhookID: uuid.NewString(),
			Config: webhooks.Config{
				ConfigUser: webhooks.ConfigUser{
					Endpoint:   httpServerFail.URL,
					Secret:     secret,
					EventTypes: []string{type1},
				},
				ID:        uuid.NewString(),
				Active:    true,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
			Payload:        fmt.Sprintf("{\"type\":\"%s\"}", type1),
			StatusCode:     http.StatusNotFound,
			Status:         webhooks.StatusAttemptToRetry,
			RetryAttempt:   0,
			NextRetryAfter: time.Now().UTC(),
		}

		_, err = mongoClient.Database(
			viper.GetString(flag.StorageMongoDatabaseName)).
			Collection(storage.CollectionAttempts).InsertOne(context.Background(), failedAttempt)
		require.NoError(t, err)

		retriesSchedule = []time.Duration{time.Second, time.Second, time.Second}
		viper.Set(flag.RetriesSchedule, retriesSchedule)

		workerRetriesApp := fxtest.New(t,
			fx.Supply(httpServerFail.Client()),
			retries.StartModule(
				viper.GetString(flag.HttpBindAddressWorkerRetries),
				viper.GetDuration(flag.RetriesCron),
				retriesSchedule))
		require.NoError(t, workerRetriesApp.Start(context.Background()))

		requestWorkerRetries(t, http.MethodGet, server.PathHealthCheck, http.StatusOK)

		expectedAttempts := 4

		attempts := 0
		for attempts != expectedAttempts {
			opts := options.Find().SetSort(bson.M{webhooks.KeyID: -1})
			cur, err := mongoClient.Database(
				viper.GetString(flag.StorageMongoDatabaseName)).
				Collection(storage.CollectionAttempts).
				Find(context.Background(), bson.M{}, opts)
			require.NoError(t, err)
			var results []webhooks.Attempt
			require.NoError(t, cur.All(context.Background(), &results))
			attempts = len(results)
			if attempts != expectedAttempts {
				time.Sleep(time.Second)
			} else {
				// First attempt should be failed
				require.Equal(t, webhooks.StatusAttemptFailed, results[0].Status)
				require.Equal(t, expectedAttempts-1, results[0].RetryAttempt)
			}
		}
		time.Sleep(time.Second)
		require.Equal(t, expectedAttempts, attempts)

		require.NoError(t, workerRetriesApp.Stop(context.Background()))
	})

	t.Run("retry long schedule", func(t *testing.T) {
		retriesSchedule = []time.Duration{time.Hour}
		viper.Set(flag.RetriesSchedule, retriesSchedule)

		require.NoError(t, mongoClient.Database(
			viper.GetString(flag.StorageMongoDatabaseName)).
			Collection(storage.CollectionAttempts).Drop(context.Background()))

		// New test server with fail handler
		httpServerFail := httptest.NewServer(http.HandlerFunc(webhooksFailHandler))
		defer func() {
			httpServerFail.CloseClientConnections()
			httpServerFail.Close()
		}()

		failedAttempt := webhooks.Attempt{
			Date:      time.Now().UTC(),
			WebhookID: uuid.NewString(),
			Config: webhooks.Config{
				ConfigUser: webhooks.ConfigUser{
					Endpoint:   httpServerFail.URL,
					Secret:     secret,
					EventTypes: []string{type1},
				},
				ID:        uuid.NewString(),
				Active:    true,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
			Payload:        fmt.Sprintf("{\"type\":\"%s\"}", type1),
			StatusCode:     http.StatusNotFound,
			Status:         webhooks.StatusAttemptToRetry,
			RetryAttempt:   0,
			NextRetryAfter: time.Now().UTC(),
		}

		_, err = mongoClient.Database(
			viper.GetString(flag.StorageMongoDatabaseName)).
			Collection(storage.CollectionAttempts).InsertOne(context.Background(), failedAttempt)
		require.NoError(t, err)

		workerRetriesApp := fxtest.New(t,
			fx.Supply(httpServerFail.Client()),
			retries.StartModule(
				viper.GetString(flag.HttpBindAddressWorkerRetries),
				viper.GetDuration(flag.RetriesCron),
				retriesSchedule))
		require.NoError(t, workerRetriesApp.Start(context.Background()))

		requestWorkerRetries(t, http.MethodGet, server.PathHealthCheck, http.StatusOK)

		time.Sleep(3 * time.Second)

		cur, err := mongoClient.Database(
			viper.GetString(flag.StorageMongoDatabaseName)).
			Collection(storage.CollectionAttempts).
			Find(context.Background(), bson.M{}, nil)
		require.NoError(t, err)
		var results []webhooks.Attempt
		require.NoError(t, cur.All(context.Background(), &results))
		attempts := len(results)
		require.Equal(t, 2, attempts)
		require.Equal(t, webhooks.StatusAttemptFailed, results[0].Status)
		require.Equal(t, webhooks.StatusAttemptFailed, results[1].Status)

		require.NoError(t, workerRetriesApp.Stop(context.Background()))
	})
}
