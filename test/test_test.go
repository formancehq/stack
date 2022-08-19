package test_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/numary/go-libs/sharedapi"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/constants"
	"github.com/numary/webhooks/internal/env"
	"github.com/numary/webhooks/internal/kafka"
	"github.com/numary/webhooks/pkg/model"
	"github.com/numary/webhooks/pkg/server"
	"github.com/numary/webhooks/pkg/worker"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx/fxtest"
)

var (
	flagSet       *pflag.FlagSet
	httpClient    *http.Client
	serverBaseURL string
	mongoClient   *mongo.Client
)

func TestMain(m *testing.M) {
	flagSet = pflag.NewFlagSet("test", pflag.ContinueOnError)
	if err := env.Init(flagSet); err != nil {
		panic(err)
	}

	httpClient = &http.Client{
		Transport: Interceptor{http.DefaultTransport},
	}
	serverBaseURL = fmt.Sprintf("http://localhost%s",
		viper.GetString(constants.HttpBindAddressServerFlag))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	mongoDBUri := viper.GetString(constants.StorageMongoConnStringFlag)
	if mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoDBUri)); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestServer(t *testing.T) {
	serverApp := fxtest.New(t, server.StartModule(httpClient))
	serverApp.RequireStart()

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
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}

		resp, err = http.Get(serverBaseURL + server.PathConfigs)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur = decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		assert.Equal(t, 0, len(cur.Data))
	})

	validConfigs := []model.Config{
		{
			Endpoint:   "https://www.site1.com",
			EventTypes: []string{"TYPE1", "TYPE2"},
		},
		{
			Endpoint:   "https://www.site2.com",
			EventTypes: []string{"TYPE3"},
		},
		{
			Endpoint:   "https://www.site3.com",
			Secret:     model.NewSecret(),
			EventTypes: []string{"TYPE1"},
		},
	}

	var insertedIds = make([]string, len(validConfigs))

	t.Run("POST "+server.PathConfigs, func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			for i, cfg := range validConfigs {
				req, err := http.NewRequest(http.MethodPost, serverBaseURL+server.PathConfigs, buffer(t, cfg))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				resp, err := httpClient.Do(req)
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&insertedIds[i]))
			}
		})

		t.Run("invalid Content-Type", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, serverBaseURL+server.PathConfigs,
				buffer(t, validConfigs[0]))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "invalid")
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
		})

		t.Run("invalid nil body", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, serverBaseURL+server.PathConfigs,
				nil)
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("invalid body not json", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, serverBaseURL+server.PathConfigs,
				bytes.NewBuffer([]byte("{")))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("invalid body unknown field", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, serverBaseURL+server.PathConfigs,
				bytes.NewBuffer([]byte("{\"field\":false}")))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("invalid body json syntax", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, serverBaseURL+server.PathConfigs,
				bytes.NewBuffer([]byte("{\"endpoint\":\"example.com\",}")))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("GET "+server.PathConfigs+" all", func(t *testing.T) {
		resp, err := http.Get(serverBaseURL + server.PathConfigs)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur := decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		assert.Equal(t, len(validConfigs), len(cur.Data))
		for i, cfg := range validConfigs {
			assert.Equal(t, cfg.Endpoint, cur.Data[len(validConfigs)-i-1].Config.Endpoint)
			assert.Equal(t, cfg.EventTypes, cur.Data[len(validConfigs)-i-1].Config.EventTypes)
		}
	})

	t.Run("PUT "+server.PathConfigs, func(t *testing.T) {
		t.Run(server.PathDeactivate, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, serverBaseURL+server.PathConfigs+
				"/"+insertedIds[0]+server.PathDeactivate, nil)
			require.NoError(t, err)
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			resp, err = http.Get(serverBaseURL + server.PathConfigs)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			cur := decodeCursorResponse[model.ConfigInserted](t, resp.Body)
			assert.Equal(t, len(validConfigs), len(cur.Data))
			assert.Equal(t, false, cur.Data[len(cur.Data)-1].Active)
		})

		t.Run(server.PathActivate, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, serverBaseURL+server.PathConfigs+
				"/"+insertedIds[0]+server.PathActivate, nil)
			require.NoError(t, err)
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			resp, err = http.Get(serverBaseURL + server.PathConfigs)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			cur := decodeCursorResponse[model.ConfigInserted](t, resp.Body)
			assert.Equal(t, len(validConfigs), len(cur.Data))
			assert.Equal(t, true, cur.Data[len(cur.Data)-1].Active)
		})

		t.Run(server.PathRotateSecret, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, serverBaseURL+server.PathConfigs+
				"/"+insertedIds[0]+server.PathRotateSecret, nil)
			require.NoError(t, err)
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			sec := model.Secret{Secret: model.NewSecret()}
			req, err = http.NewRequest(http.MethodPut, serverBaseURL+server.PathConfigs+
				"/"+insertedIds[0]+server.PathRotateSecret, buffer(t, sec))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			resp, err = httpClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			sec = model.Secret{Secret: "invalid"}
			req, err = http.NewRequest(http.MethodPut, serverBaseURL+server.PathConfigs+
				"/"+insertedIds[0]+server.PathRotateSecret, buffer(t, sec))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			resp, err = httpClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			req, err = http.NewRequest(http.MethodPut, serverBaseURL+server.PathConfigs+
				"/"+insertedIds[0]+server.PathRotateSecret, buffer(t, validConfigs[0]))
			require.NoError(t, err)
			resp, err = httpClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("DELETE "+server.PathConfigs, func(t *testing.T) {
		for _, id := range insertedIds {
			req, err := http.NewRequest(http.MethodDelete, serverBaseURL+server.PathConfigs+"/"+id, nil)
			require.NoError(t, err)
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			req, err = http.NewRequest(http.MethodDelete, serverBaseURL+server.PathConfigs+"/"+id, nil)
			require.NoError(t, err)
			resp, err = httpClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		}
	})

	t.Run("GET "+server.PathConfigs+" after delete", func(t *testing.T) {
		resp, err := http.Get(serverBaseURL + server.PathConfigs)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur := decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		assert.Equal(t, 0, len(cur.Data))
	})

	serverApp.RequireStop()
}

func TestWorker(t *testing.T) {
	serverApp := fxtest.New(t, server.StartModule(httpClient))
	serverApp.RequireStart()
	workerApp := fxtest.New(t, worker.StartModule(context.Background(), httpClient))
	workerApp.RequireStart()

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

	workerApp.RequireStop()
	serverApp.RequireStop()
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

// Intercept every http request from httpClient to store webhooks sent
type Interceptor struct {
	core http.RoundTripper
}

type message struct {
	Url string `json:"url" bson:"url"`
}

func (i Interceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	if strings.Contains(r.URL.String(), "msg") {
		_, err := mongoClient.Database(
			viper.GetString(constants.StorageMongoDatabaseNameFlag)).
			Collection("messages").InsertOne(context.Background(), message{r.URL.String()})
		if err != nil {
			return nil, fmt.Errorf("Interceptor.RoundTrip: %w", err)
		}
	}

	// send the request using the DefaultTransport
	return i.core.RoundTrip(r)
}

func buffer(t *testing.T, v any) *bytes.Buffer {
	data, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewBuffer(data)
}

func decodeCursorResponse[T any](t *testing.T, reader io.Reader) *sharedapi.Cursor[T] {
	res := sharedapi.BaseResponse[T]{}
	err := json.NewDecoder(reader).Decode(&res)
	require.NoError(t, err)
	return res.Cursor
}
