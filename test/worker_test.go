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

	"github.com/numary/webhooks/constants"
	webhooks "github.com/numary/webhooks/pkg"
	"github.com/numary/webhooks/pkg/kafka"
	"github.com/numary/webhooks/pkg/security"
	"github.com/numary/webhooks/pkg/server"
	"github.com/numary/webhooks/pkg/worker"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx/fxtest"
)

func TestWorker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx,
		options.Client().ApplyURI(
			viper.GetString(constants.StorageMongoConnStringFlag)))
	require.NoError(t, err)

	// Cleanup collections
	require.NoError(t, mongoClient.Database(
		viper.GetString(constants.StorageMongoDatabaseNameFlag)).
		Collection(constants.MongoCollectionConfigs).Drop(context.Background()))

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
		server.StartModule(
			viper.GetString(constants.HttpBindAddressServerFlag)))
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
			viper.GetString(constants.StorageMongoDatabaseNameFlag)).
			Collection(constants.MongoCollectionRequests).Drop(context.Background()))

		workerApp := fxtest.New(t,
			worker.StartModule(
				viper.GetString(constants.HttpBindAddressWorkerFlag), httpServerSuccess.Client()))
		require.NoError(t, workerApp.Start(context.Background()))

		t.Run("health check", func(t *testing.T) {
			requestWorker(t, http.MethodGet, server.PathHealthCheck, http.StatusOK)
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
					viper.GetString(constants.StorageMongoDatabaseNameFlag)).
					Collection(constants.MongoCollectionRequests).
					Find(context.Background(), bson.M{}, nil)
				require.NoError(t, err)
				var results []webhooks.Request
				require.NoError(t, cur.All(context.Background(), &results))
				msgs = len(results)
				if msgs != expectedSentWebhooks {
					time.Sleep(time.Second)
				} else {
					for _, res := range results {
						require.Equal(t, true, res.Success)
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
			viper.GetString(constants.StorageMongoDatabaseNameFlag)).
			Collection(constants.MongoCollectionRequests).Drop(context.Background()))

		workerApp := fxtest.New(t,
			worker.StartModule(
				viper.GetString(constants.HttpBindAddressWorkerFlag), httpServerFail.Client()))
		require.NoError(t, workerApp.Start(context.Background()))

		t.Run("health check", func(t *testing.T) {
			requestWorker(t, http.MethodGet, server.PathHealthCheck, http.StatusOK)
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
					viper.GetString(constants.StorageMongoDatabaseNameFlag)).
					Collection(constants.MongoCollectionRequests).
					Find(context.Background(), bson.M{}, nil)
				require.NoError(t, err)
				var results []webhooks.Request
				require.NoError(t, cur.All(context.Background(), &results))
				msgs = len(results)
				if msgs != expectedSentWebhooks {
					time.Sleep(time.Second)
				} else {
					for _, res := range results {
						require.Equal(t, false, res.Success)
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

func webhooksSuccessHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("formance-webhook-id")
	ts := r.Header.Get("formance-webhook-timestamp")
	signatures := r.Header.Get("formance-webhook-signature")
	timeInt, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	timestamp := time.Unix(timeInt, 0)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ok, err := security.Verify(signatures, id, timestamp, secret, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintf(w, "WEBHOOK RECEIVED: FAKE OK RESPONSE\n")
	return
}

func webhooksFailHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "WEBHOOKS RECEIVED: FAKE NOT FOUND ERROR", http.StatusNotFound)
	return
}
