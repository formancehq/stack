package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/numary/go-libs/sharedapi"
	"github.com/stretchr/testify/require"
)

func createJSONBuffer(t *testing.T, v any) io.Reader {
	data, err := json.Marshal(v)
	require.NoError(t, err)

	return bytes.NewBuffer(data)
}

func readObject[T any](t *testing.T, recorder *httptest.ResponseRecorder) T {
	body := sharedapi.BaseResponse[T]{}
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	return *body.Data
}
