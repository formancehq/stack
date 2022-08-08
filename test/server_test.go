package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/numary/go-libs/sharedapi"
	"github.com/numary/webhooks-cloud/cmd/constants"
	"github.com/numary/webhooks-cloud/internal/env"
	"github.com/numary/webhooks-cloud/pkg/model"
	"github.com/numary/webhooks-cloud/pkg/server"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func TestServer(t *testing.T) {
	flagSet := pflag.NewFlagSet("TestServer", pflag.ContinueOnError)
	require.NoError(t, env.Init(flagSet))

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

	req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
		buffer(t, model.Config{
			Endpoint:   "https://www.site1.com",
			EventTypes: []string{"COMMITTED_TRANSACTIONS", "SAVED_METADATA"},
		}))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	insertedId1 := ""
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&insertedId1))

	req, err = http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath+server.TogglePath+"/"+insertedId1, nil)
	require.NoError(t, err)
	resp, err = c.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
		buffer(t, model.Config{
			Endpoint:   "https://www.site3.com",
			EventTypes: []string{"COMMITTED_TRANSACTIONS"},
		}))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err = c.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	insertedId2 := ""
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&insertedId2))

	t.Run("get all configs", func(t *testing.T) {
		resp, err = http.Get(baseURL + server.ConfigsPath)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		cur := decodeCursorResponse[model.ConfigInserted](t, resp.Body)
		assert.Equal(t, 2, len(cur.Data))
		assert.Equal(t, 1, len(cur.Data[0].EventTypes))
		assert.Equal(t, "COMMITTED_TRANSACTIONS", cur.Data[0].EventTypes[0])
	})

	t.Run("invalid Content-Type", func(t *testing.T) {
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

	t.Run("invalid body syntax", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, baseURL+server.ConfigsPath,
			bytes.NewBuffer([]byte("{\"endpoint\":\"example.com\",}")))
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
