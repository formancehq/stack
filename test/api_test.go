package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/numary/go-libs/sharedapi"
	"github.com/numary/webhooks-cloud/api/server"
	"github.com/numary/webhooks-cloud/cmd/constants"
	"github.com/numary/webhooks-cloud/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestAPI(t *testing.T) {
	app := fxtest.New(t, server.StartModule())
	app.RequireStart()

	baseURL := fmt.Sprintf("http://localhost%s", constants.DefaultBindAddress)
	c := http.DefaultClient

	t.Run("health check", func(t *testing.T) {
		resp, err := http.Get(baseURL + server.HealthCheckPath)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("no configs", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, baseURL+server.ConfigsPath, nil)
		require.NoError(t, err)
		resp, err := c.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp, err = http.Get(baseURL + server.ConfigsPath)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur := decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		assert.Equal(t, 0, len(cur.Data))
	})

	req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
		buffer(t, model.Config{
			Active:     true,
			EventTypes: []string{"TYPE1", "TYPE2"},
			Endpoints:  []string{"https://www.site1.com", "https://www.site2.com"},
		}))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
		buffer(t, model.Config{
			Active:     true,
			EventTypes: []string{"TYPE3"},
			Endpoints:  []string{"https://www.site3.com"},
		}))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err = c.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
		buffer(t, model.Config{
			Active:     false,
			EventTypes: []string{},
			Endpoints:  []string{},
		}))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err = c.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	t.Run("get all configs", func(t *testing.T) {
		resp, err = http.Get(baseURL + server.ConfigsPath)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur := decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		assert.Equal(t, 3, len(cur.Data))
		assert.Equal(t, false, cur.Data[0].Active)
	})

	t.Run("delete all configs", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, baseURL+server.ConfigsPath, nil)
		require.NoError(t, err)
		resp, err := c.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp, err = http.Get(baseURL + server.ConfigsPath)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur := decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		assert.Equal(t, 0, len(cur.Data))
	})

	t.Run("invalid config", func(t *testing.T) {
		req, err = http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
			buffer(t, model.Config{
				Active:    false,
				Endpoints: []string{"https://www.site1.com"},
			}))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err = c.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		req, err = http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
			buffer(t, model.Config{
				Active:     false,
				EventTypes: []string{"TYPE"},
			}))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err = c.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid content type", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
			buffer(t, model.Config{}))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "invalid")
		resp, err = c.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)
	})

	t.Run("nil body", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
			nil)
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)
		resp, err = c.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid body", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
			bytes.NewBuffer([]byte("{")))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)
		resp, err = c.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid body double json", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
			bytes.NewBuffer([]byte("{\"active\":false}{\"active\":false}")))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err = c.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid body unknown field", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
			bytes.NewBuffer([]byte("{\"field\":false}")))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err = c.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid body invalid value", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
			bytes.NewBuffer([]byte("{\"active\":1}")))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err = c.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid body syntax", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
			bytes.NewBuffer([]byte("{\"active\":true,}")))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err = c.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	app.RequireStop()
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
