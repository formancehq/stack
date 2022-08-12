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
	"github.com/numary/webhooks/constants"
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

	baseURL := fmt.Sprintf("http://localhost%s", constants.DefaultBindAddressServer)
	c := http.DefaultClient

	t.Run("health check", func(t *testing.T) {
		resp, err := http.Get(baseURL + server.PathHealthCheck)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("clean existing configs", func(t *testing.T) {
		resp, err := http.Get(baseURL + server.PathConfigs)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur := decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		for _, cfg := range cur.Data {
			req, err := http.NewRequest(http.MethodDelete, baseURL+server.PathConfigs+"/"+cfg.ID, nil)
			require.NoError(t, err)
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}

		resp, err = http.Get(baseURL + server.PathConfigs)
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
			EventTypes: []string{"TYPE1"},
		},
	}

	var insertedIds = make([]string, len(validConfigs))

	t.Run("POST "+server.PathConfigs, func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			for i, cfg := range validConfigs {
				req, err := http.NewRequest(http.MethodPost, baseURL+server.PathConfigs, buffer(t, cfg))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				resp, err := c.Do(req)
				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&insertedIds[i]))
			}
		})

		t.Run("invalid Content-Type", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, baseURL+server.PathConfigs,
				buffer(t, validConfigs[0]))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "invalid")
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
		})

		t.Run("invalid nil body", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, baseURL+server.PathConfigs,
				nil)
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("invalid body not json", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, baseURL+server.PathConfigs,
				bytes.NewBuffer([]byte("{")))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("invalid body double json", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, baseURL+server.PathConfigs,
				bytes.NewBuffer([]byte("{\"active\":false}{\"active\":false}")))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("invalid body unknown field", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, baseURL+server.PathConfigs,
				bytes.NewBuffer([]byte("{\"field\":false}")))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("invalid body json syntax", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, baseURL+server.PathConfigs,
				bytes.NewBuffer([]byte("{\"endpoint\":\"example.com\",}")))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})

	t.Run("GET "+server.PathConfigs+" all", func(t *testing.T) {
		resp, err := http.Get(baseURL + server.PathConfigs)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur := decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		assert.Equal(t, len(validConfigs), len(cur.Data))
		for i, cfg := range validConfigs {
			assert.Equal(t, cfg, cur.Data[len(validConfigs)-i-1].Config)
		}
	})

	t.Run("PUT "+server.PathConfigs, func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, baseURL+server.PathConfigs+"/"+insertedIds[0]+server.PathDeactivate, nil)
		require.NoError(t, err)
		resp, err := c.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp, err = http.Get(baseURL + server.PathConfigs)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur := decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		assert.Equal(t, len(validConfigs), len(cur.Data))
		assert.Equal(t, false, cur.Data[len(cur.Data)-1].Active)

		req, err = http.NewRequest(http.MethodPut, baseURL+server.PathConfigs+"/"+insertedIds[0]+server.PathActivate, nil)
		require.NoError(t, err)
		resp, err = c.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp, err = http.Get(baseURL + server.PathConfigs)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur = decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		assert.Equal(t, len(validConfigs), len(cur.Data))
		assert.Equal(t, true, cur.Data[len(cur.Data)-1].Active)
	})

	t.Run("DELETE "+server.PathConfigs, func(t *testing.T) {
		for _, id := range insertedIds {
			req, err := http.NewRequest(http.MethodDelete, baseURL+server.PathConfigs+"/"+id, nil)
			require.NoError(t, err)
			resp, err := c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			req, err = http.NewRequest(http.MethodDelete, baseURL+server.PathConfigs+"/"+id, nil)
			require.NoError(t, err)
			resp, err = c.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		}
	})

	t.Run("GET "+server.PathConfigs+" after delete", func(t *testing.T) {
		resp, err := http.Get(baseURL + server.PathConfigs)
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
