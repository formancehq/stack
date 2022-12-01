package test_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/formancehq/webhooks/cmd/flag"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/kafka"
	"github.com/formancehq/webhooks/pkg/security"
	"github.com/formancehq/webhooks/pkg/server"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/formancehq/webhooks/pkg/worker"
	"github.com/google/uuid"
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

		retrySchedule = []time.Duration{time.Second}
		viper.Set(flag.RetrySchedule, retrySchedule)

		workerApp := fxtest.New(t,
			fx.Supply(httpServerSuccess.Client()),
			worker.StartModule(
				viper.GetString(flag.HttpBindAddressWorker),
				viper.GetDuration(flag.RetryCron),
				retrySchedule,
			))
		require.NoError(t, workerApp.Start(context.Background()))

		healthCheckWorker(t)

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

		require.NoError(t, workerApp.Stop(context.Background()))
	})

	t.Run("failure", func(t *testing.T) {
		require.NoError(t, mongoClient.Database(
			viper.GetString(flag.StorageMongoDatabaseName)).
			Collection(storage.CollectionAttempts).Drop(context.Background()))

		retrySchedule = []time.Duration{time.Second}
		viper.Set(flag.RetrySchedule, retrySchedule)

		workerApp := fxtest.New(t,
			fx.Supply(httpServerFail.Client()),
			worker.StartModule(
				viper.GetString(flag.HttpBindAddressWorker),
				viper.GetDuration(flag.RetryCron),
				retrySchedule,
			))
		require.NoError(t, workerApp.Start(context.Background()))

		healthCheckWorker(t)

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

		require.NoError(t, workerApp.Stop(context.Background()))
	})

	require.NoError(t, serverApp.Stop(context.Background()))
}

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

		retrySchedule = []time.Duration{time.Second, time.Second, time.Second}
		viper.Set(flag.RetrySchedule, retrySchedule)

		workerApp := fxtest.New(t,
			fx.Supply(httpServerSuccess.Client()),
			worker.StartModule(
				viper.GetString(flag.HttpBindAddressWorker),
				viper.GetDuration(flag.RetryCron),
				retrySchedule))
		require.NoError(t, workerApp.Start(context.Background()))

		healthCheckWorker(t)

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

		require.NoError(t, workerApp.Stop(context.Background()))
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

		retrySchedule = []time.Duration{time.Second, time.Second, time.Second}
		viper.Set(flag.RetrySchedule, retrySchedule)

		workerApp := fxtest.New(t,
			fx.Supply(httpServerFail.Client()),
			worker.StartModule(
				viper.GetString(flag.HttpBindAddressWorker),
				viper.GetDuration(flag.RetryCron),
				retrySchedule))
		require.NoError(t, workerApp.Start(context.Background()))

		healthCheckWorker(t)

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

		require.NoError(t, workerApp.Stop(context.Background()))
	})

	t.Run("retry long schedule", func(t *testing.T) {
		retrySchedule = []time.Duration{time.Hour}
		viper.Set(flag.RetrySchedule, retrySchedule)

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

		workerApp := fxtest.New(t,
			fx.Supply(httpServerFail.Client()),
			worker.StartModule(
				viper.GetString(flag.HttpBindAddressWorker),
				viper.GetDuration(flag.RetryCron),
				retrySchedule))
		require.NoError(t, workerApp.Start(context.Background()))

		healthCheckWorker(t)

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

		require.NoError(t, workerApp.Stop(context.Background()))
	})
}

func webhooksSuccessHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("formance-webhook-id")
	ts := r.Header.Get("formance-webhook-timestamp")
	signatures := r.Header.Get("formance-webhook-signature")
	timeInt, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ok, err := security.Verify(signatures, id, timeInt, secret, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "security.Verify NOK", http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintf(w, "WEBHOOK RECEIVED: MOCK OK RESPONSE\n")
	return
}

func webhooksFailHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "WEBHOOKS RECEIVED: MOCK ERROR RESPONSE", http.StatusNotFound)
	return
}
