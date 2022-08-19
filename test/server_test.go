package test_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/numary/webhooks/internal/model"
	"github.com/numary/webhooks/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestServer(t *testing.T) {
	serverApp := fxtest.New(t, server.StartModule(context.Background(), httpClient))

	t.Run("start", func(t *testing.T) {
		serverApp.RequireStart()
	})

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

	t.Run("stop", func(t *testing.T) {
		serverApp.RequireStop()
	})
}
