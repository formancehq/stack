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
	"github.com/numary/webhooks/constants"
	webhooks "github.com/numary/webhooks/pkg"
	"github.com/numary/webhooks/pkg/env"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

var (
	httpClient = http.DefaultClient

	serverBaseURL string
	workerBaseURL string

	secret = webhooks.NewSecret()

	topic = strings.ReplaceAll(
		time.Now().UTC().Format(time.RFC3339Nano), ":", "-")
)

func TestMain(m *testing.M) {
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	if err := env.Init(flagSet); err != nil {
		panic(err)
	}

	viper.Set(constants.KafkaTopicsFlag, []string{topic})

	serverBaseURL = fmt.Sprintf("http://localhost%s",
		viper.GetString(constants.HttpBindAddressServerFlag))
	workerBaseURL = fmt.Sprintf("http://localhost%s",
		viper.GetString(constants.HttpBindAddressWorkerFlag))

	os.Exit(m.Run())
}

func requestServer(t *testing.T, method, url string, expectedCode int, body ...any) io.ReadCloser {
	return request(t, method, serverBaseURL+url, body, expectedCode)
}

func requestWorker(t *testing.T, method, url string, expectedCode int, body ...any) {
	request(t, method, workerBaseURL+url, body, expectedCode)
}

func request(t *testing.T, method, url string, body []any, expectedCode int) io.ReadCloser {
	var err error
	var req *http.Request
	if len(body) > 0 {
		req, err = http.NewRequestWithContext(context.Background(), method, url, buffer(t, body[0]))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequestWithContext(context.Background(), method, url, nil)
	}
	require.NoError(t, err)
	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, expectedCode, resp.StatusCode)
	return resp.Body
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
