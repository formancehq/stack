package test_test

//
//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"io"
//	"net/http"
//	"os"
//	"strings"
//	"testing"
//	"time"
//
//	"github.com/formancehq/stack/libs/go-libs/api"
//	"github.com/formancehq/stack/libs/go-libs/logging"
//	"github.com/formancehq/stack/libs/go-libs/pgtesting"
//	"github.com/formancehq/stack/libs/go-libs/publish"
//	"github.com/formancehq/webhooks/cmd/flag"
//	webhooks "github.com/formancehq/webhooks/pkg"
//	"github.com/formancehq/webhooks/pkg/server"
//	"github.com/spf13/pflag"
//	"github.com/spf13/viper"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//)
//
//var (
//	httpClient = http.DefaultClient
//
//	serverBaseURL string
//	workerBaseURL string
//
//	secret = webhooks.NewSecret()
//
//	topic = strings.ReplaceAll(
//		time.Now().UTC().Format(time.RFC3339Nano), ":", "-")
//
//	retrySchedule []time.Duration
//)
//
//func TestMain(m *testing.M) {
//	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
//	var errInit error
//	if _, errInit = flag.Init(flagSet); errInit != nil {
//		panic(errInit)
//	}
//
//	if err := pgtesting.CreatePostgresServer(); err != nil {
//		logging.Error(err)
//		os.Exit(1)
//	}
//
//	viper.Set(flag.StoragePostgresConnString, pgtesting.Server().GetDSN())
//	viper.Set(flag.KafkaTopics, []string{topic})
//	viper.Set(flag.RetriesCron, time.Second)
//	viper.Set(publish.PublisherKafkaEnabledFlag, true)
//	viper.Set(publish.PublisherKafkaBrokerFlag, []string{"localhost:9092"})
//
//	serverBaseURL = "http://localhost:8080"
//	workerBaseURL = "http://localhost:8080"
//
//	code := m.Run()
//	if err := pgtesting.DestroyPostgresServer(); err != nil {
//		logging.Error(err)
//	}
//
//	os.Exit(code)
//}
//
//func requestServer(t *testing.T, method, url string, expectedCode int, body ...any) io.ReadCloser {
//	return request(t, method, serverBaseURL+url, body, expectedCode)
//}
//
//func healthCheckWorker(t *testing.T) {
//	request(t, http.MethodGet, workerBaseURL+server.PathHealthCheck, nil, http.StatusOK)
//}
//
//func request(t *testing.T, method, url string, body []any, expectedCode int) io.ReadCloser {
//	var err error
//	var req *http.Request
//	if len(body) > 0 {
//		req, err = http.NewRequestWithContext(context.Background(), method, url, buffer(t, body[0]))
//		req.Header.Set("Content-Type", "application/json")
//	} else {
//		req, err = http.NewRequestWithContext(context.Background(), method, url, nil)
//	}
//	require.NoError(t, err)
//	resp, err := httpClient.Do(req)
//	require.NoError(t, err)
//	require.Equal(t, expectedCode, resp.StatusCode)
//	return resp.Body
//}
//
//func buffer(t *testing.T, v any) *bytes.Buffer {
//	data, err := json.Marshal(v)
//	require.NoError(t, err)
//	return bytes.NewBuffer(data)
//}
//
//func decodeCursorResponse[T any](t *testing.T, reader io.Reader) *api.Cursor[T] {
//	res := api.BaseResponse[T]{}
//	err := json.NewDecoder(reader).Decode(&res)
//	require.NoError(t, err)
//	return res.Cursor
//}
//
//func decodeSingleResponse[T any](t *testing.T, reader io.Reader) (T, bool) {
//	res := api.BaseResponse[T]{}
//	if !decode(t, reader, &res) {
//		var zero T
//		return zero, false
//	}
//	return *res.Data, true
//}
//
//func decode(t *testing.T, reader io.Reader, v interface{}) bool {
//	err := json.NewDecoder(reader).Decode(v)
//	return assert.NoError(t, err)
//}
