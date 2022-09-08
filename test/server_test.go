package test_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/numary/webhooks/constants"
	webhooks "github.com/numary/webhooks/pkg"
	"github.com/numary/webhooks/pkg/server"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestServer(t *testing.T) {
	serverApp := fxtest.New(t,
		server.StartModule(
			viper.GetString(constants.HttpBindAddressServerFlag)))

	t.Run("start", func(t *testing.T) {
		serverApp.RequireStart()
	})

	t.Run("health check", func(t *testing.T) {
		requestServer(t, http.MethodGet, server.PathHealthCheck, http.StatusOK)
	})

	t.Run("clean existing configs", func(t *testing.T) {
		resBody := requestServer(t, http.MethodGet, server.PathConfigs, http.StatusOK)
		cur := decodeCursorResponse[webhooks.Config](t, resBody)
		for _, cfg := range cur.Data {
			requestServer(t, http.MethodDelete, server.PathConfigs+"/"+cfg.ID, http.StatusOK)
		}
		require.NoError(t, resBody.Close())

		resBody = requestServer(t, http.MethodGet, server.PathConfigs, http.StatusOK)
		cur = decodeCursorResponse[webhooks.Config](t, resBody)
		assert.Equal(t, 0, len(cur.Data))
		require.NoError(t, resBody.Close())
	})

	validConfigs := []webhooks.ConfigUser{
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
			Secret:     webhooks.NewSecret(),
			EventTypes: []string{"TYPE1"},
		},
	}

	insertedIds := make([]string, len(validConfigs))

	t.Run("POST "+server.PathConfigs, func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			for i, cfg := range validConfigs {
				resBody := requestServer(t, http.MethodPost, server.PathConfigs, http.StatusOK, cfg)
				assert.NoError(t, json.NewDecoder(resBody).Decode(&insertedIds[i]))
				require.NoError(t, resBody.Close())
			}
		})

		t.Run("invalid Content-Type", func(t *testing.T) {
			req, err := http.NewRequestWithContext(context.Background(),
				http.MethodPost, serverBaseURL+server.PathConfigs,
				buffer(t, validConfigs[0]))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "invalid")
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
			require.NoError(t, resp.Body.Close())
		})

		t.Run("invalid nil body", func(t *testing.T) {
			requestServer(t, http.MethodPost, server.PathConfigs, http.StatusBadRequest)
		})

		t.Run("invalid body not json", func(t *testing.T) {
			requestServer(t, http.MethodPost, server.PathConfigs, http.StatusBadRequest, []byte("{"))
		})

		t.Run("invalid body unknown field", func(t *testing.T) {
			requestServer(t, http.MethodPost, server.PathConfigs, http.StatusBadRequest, []byte("{\"field\":false}"))
		})

		t.Run("invalid body json syntax", func(t *testing.T) {
			requestServer(t, http.MethodPost, server.PathConfigs, http.StatusBadRequest, []byte("{\"endpoint\":\"example.com\",}"))
		})
	})

	t.Run("GET "+server.PathConfigs, func(t *testing.T) {
		resBody := requestServer(t, http.MethodGet, server.PathConfigs, http.StatusOK)
		cur := decodeCursorResponse[webhooks.Config](t, resBody)
		assert.Equal(t, len(validConfigs), len(cur.Data))
		for i, cfg := range validConfigs {
			assert.Equal(t, cfg.Endpoint, cur.Data[len(validConfigs)-i-1].Endpoint)
			assert.Equal(t, len(cfg.EventTypes), len(cur.Data[len(validConfigs)-i-1].EventTypes))
			for j, typ := range cfg.EventTypes {
				assert.Equal(t,
					strings.ToLower(typ),
					strings.ToLower(cur.Data[len(validConfigs)-i-1].EventTypes[j]))
			}
		}
		require.NoError(t, resBody.Close())

		cfg := validConfigs[0]
		ep := url.QueryEscape(cfg.Endpoint)
		resBody = requestServer(t, http.MethodGet, server.PathConfigs+"?endpoint="+ep, http.StatusOK)
		cur = decodeCursorResponse[webhooks.Config](t, resBody)
		assert.Equal(t, 1, len(cur.Data))
		assert.Equal(t, cfg.Endpoint, cur.Data[0].Endpoint)
		require.NoError(t, resBody.Close())

		resBody = requestServer(t, http.MethodGet, server.PathConfigs+"?id="+insertedIds[0], http.StatusOK)
		cur = decodeCursorResponse[webhooks.Config](t, resBody)
		assert.Equal(t, 1, len(cur.Data))
		assert.Equal(t, cfg.Endpoint, cur.Data[0].Endpoint)
		require.NoError(t, resBody.Close())
	})

	t.Run("PUT "+server.PathConfigs, func(t *testing.T) {
		t.Run(server.PathDeactivate, func(t *testing.T) {
			requestServer(t, http.MethodPut, server.PathConfigs+"/"+insertedIds[0]+server.PathDeactivate, http.StatusOK)

			resBody := requestServer(t, http.MethodGet, server.PathConfigs, http.StatusOK)
			cur := decodeCursorResponse[webhooks.Config](t, resBody)
			assert.Equal(t, len(validConfigs), len(cur.Data))
			assert.Equal(t, false, cur.Data[0].Active)
			require.NoError(t, resBody.Close())

			requestServer(t, http.MethodPut, server.PathConfigs+"/"+insertedIds[0]+server.PathDeactivate, http.StatusNotModified)
		})

		t.Run(server.PathActivate, func(t *testing.T) {
			requestServer(t, http.MethodPut, server.PathConfigs+"/"+insertedIds[0]+server.PathActivate, http.StatusOK)

			resBody := requestServer(t, http.MethodGet, server.PathConfigs, http.StatusOK)
			cur := decodeCursorResponse[webhooks.Config](t, resBody)
			assert.Equal(t, len(validConfigs), len(cur.Data))
			assert.Equal(t, true, cur.Data[len(cur.Data)-1].Active)
			require.NoError(t, resBody.Close())

			requestServer(t, http.MethodPut, server.PathConfigs+"/"+insertedIds[0]+server.PathActivate, http.StatusNotModified)
		})

		t.Run(server.PathRotateSecret, func(t *testing.T) {
			requestServer(t, http.MethodPut, server.PathConfigs+"/"+insertedIds[0]+server.PathRotateSecret, http.StatusOK)

			validSecret := webhooks.Secret{Secret: webhooks.NewSecret()}
			requestServer(t, http.MethodPut, server.PathConfigs+"/"+insertedIds[0]+server.PathRotateSecret, http.StatusOK, validSecret)

			invalidSecret := webhooks.Secret{Secret: "invalid"}
			requestServer(t, http.MethodPut, server.PathConfigs+"/"+insertedIds[0]+server.PathRotateSecret, http.StatusBadRequest, invalidSecret)

			invalidSecret2 := validConfigs[0]
			requestServer(t, http.MethodPut, server.PathConfigs+"/"+insertedIds[0]+server.PathRotateSecret, http.StatusBadRequest, invalidSecret2)
		})
	})

	t.Run("DELETE "+server.PathConfigs, func(t *testing.T) {
		for _, id := range insertedIds {
			requestServer(t, http.MethodDelete, server.PathConfigs+"/"+id, http.StatusOK)
			requestServer(t, http.MethodDelete, server.PathConfigs+"/"+id, http.StatusNotFound)
		}
	})

	t.Run("GET "+server.PathConfigs+" after delete", func(t *testing.T) {
		resBody := requestServer(t, http.MethodGet, server.PathConfigs, http.StatusOK)
		cur := decodeCursorResponse[webhooks.Config](t, resBody)
		assert.Equal(t, 0, len(cur.Data))
		require.NoError(t, resBody.Close())
	})

	t.Run("stop", func(t *testing.T) {
		serverApp.RequireStop()
	})
}
