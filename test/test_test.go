package test_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/numary/go-libs/sharedapi"
	"github.com/numary/webhooks/cmd/constants"
	"github.com/numary/webhooks/internal/env"
	"github.com/numary/webhooks/pkg/model"
	"github.com/numary/webhooks/pkg/server"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

var testFlagSet *pflag.FlagSet

func TestMain(m *testing.M) {
	testFlagSet = pflag.NewFlagSet("test", pflag.ContinueOnError)
	if err := env.Init(testFlagSet); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestServer(t *testing.T) {
	serverApp := fxtest.New(t, server.StartModule())
	serverApp.RequireStart()

	baseURL := fmt.Sprintf("http://localhost%s", constants.DefaultBindAddress)
	c := http.DefaultClient

	t.Run("health check", func(t *testing.T) {
		resp, err := http.Get(baseURL + server.HealthCheckPath)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("clean existing configs", func(t *testing.T) {
		resp, err := http.Get(baseURL + server.ConfigsPath)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur := decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		for _, cfg := range cur.Data {
			req, err := http.NewRequest(http.MethodDelete, baseURL+server.ConfigsPath+"/"+cfg.ID, nil)
			require.NoError(t, err)
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}

		resp, err = http.Get(baseURL + server.ConfigsPath)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur = decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		assert.Equal(t, 0, len(cur.Data))
	})

	bind, err := testFlagSet.GetString(constants.ServerHttpBindAddressFlag)
	require.NoError(t, err)

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
			Endpoint:   "http://localhost" + bind + server.HealthCheckPath,
			EventTypes: []string{"TYPE1"},
		},
	}

	var insertedIds = make([]string, len(validConfigs))

	t.Run("POST "+server.ConfigsPath, func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			for i, cfg := range validConfigs {
				req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath, buffer(t, cfg))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				resp, err := c.Do(req)
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&insertedIds[i]))
			}
		})

		t.Run("invalid Content-Type", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
				buffer(t, validConfigs[0]))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "invalid")
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
		})

		t.Run("invalid nil body", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
				nil)
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("invalid body not json", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
				bytes.NewBuffer([]byte("{")))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("invalid body double json", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
				bytes.NewBuffer([]byte("{\"active\":false}{\"active\":false}")))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("invalid body unknown field", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
				bytes.NewBuffer([]byte("{\"field\":false}")))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("invalid body json syntax", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
				bytes.NewBuffer([]byte("{\"endpoint\":\"example.com\",}")))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("GET "+server.ConfigsPath, func(t *testing.T) {
		resp, err := http.Get(baseURL + server.ConfigsPath)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur := decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		assert.Equal(t, len(validConfigs), len(cur.Data))
		for i, cfg := range validConfigs {
			assert.Equal(t, cfg, cur.Data[len(validConfigs)-i-1].Config)
		}
	})

	/*
		t.Run("WORKER", func(t *testing.T) {
			workerApp := fxtest.New(t, worker.StartModule())
			workerApp.RequireStart()

			brokers, err := testFlagSet.GetStringSlice(constants.KafkaBrokersFlag)
			require.NoError(t, err)
			topics, err := testFlagSet.GetStringSlice(constants.KafkaTopicsFlag)
			require.NoError(t, err)

			conn, err := kafkago.DialLeader(context.Background(),
				"tcp", brokers[0], topics[0], 0)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, conn.Close())
			}()

			eventType := validConfigs[2].EventTypes[0]
			nbBytes, err := conn.WriteMessages(newEventMessage(t, eventType))
			require.NoError(t, err)
			require.NotEqual(t, 0, nbBytes)

			res, err := ts.Client().Get(ts.URL)
			require.NoError(t, err)
			message, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			require.NoError(t, res.Body.Close())
			fmt.Printf("MESSAGE: %s", message)

			workerApp.RequireStop()
		})
	*/

	t.Run("DELETE "+server.ConfigsPath, func(t *testing.T) {
		for _, id := range insertedIds {
			req, err := http.NewRequest(http.MethodDelete, baseURL+server.ConfigsPath+"/"+id, nil)
			require.NoError(t, err)
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("GET "+server.ConfigsPath, func(t *testing.T) {
		resp, err := http.Get(baseURL + server.ConfigsPath)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur := decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		assert.Equal(t, 0, len(cur.Data))
	})

	serverApp.RequireStop()
}

func buffer(t *testing.T, v any) *bytes.Buffer {
	data, err := json.Marshal(v)
	assert.NoError(t, err)
	return bytes.NewBuffer(data)
}

func decodeCursorResponse[T any](t *testing.T, reader io.Reader) *sharedapi.Cursor[T] {
	res := sharedapi.BaseResponse[T]{}
	err := json.NewDecoder(reader).Decode(&res)
	assert.NoError(t, err)
	return res.Cursor
}

/*
func newEventMessage(t *testing.T, eventType string) kafkago.Message {
	ev := kafka.Event{
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
*/
