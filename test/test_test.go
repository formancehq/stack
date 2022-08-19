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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	httpClient    *http.Client
	serverBaseURL string
	mongoClient   *mongo.Client
)

func TestMain(m *testing.M) {
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
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

// Interceptor intercepts every http request from httpClient to store webhooks sent
type Interceptor struct {
	core http.RoundTripper
}

type message struct {
	Url string `json:"url" bson:"url"`
}

// Intercept the message requests to the Svix API and store them in a 'messages' Mongo collection
func (i Interceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.Contains(req.URL.String(), "/msg/") {
		sharedlogging.Debugf("request intercepted: %s", req.URL.String())
		_, err := mongoClient.Database(
			viper.GetString(constants.StorageMongoDatabaseNameFlag)).
			Collection("messages").InsertOne(context.Background(), message{req.URL.String()})
		if err != nil {
			return nil, fmt.Errorf("Interceptor.RoundTrip: %w", err)
		}
	}

	// send the request using the DefaultTransport
	return i.core.RoundTrip(req)
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
