package controllers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	service "github.com/formancehq/stack/components/stargate/internal/api"
	"github.com/formancehq/stack/components/stargate/internal/server/http/controllers"
	"github.com/formancehq/stack/components/stargate/internal/server/http/opentelemetry"
	"github.com/formancehq/stack/components/stargate/internal/server/http/routes"
	"github.com/formancehq/stack/components/stargate/internal/utils"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestStargateController(t *testing.T) {
	t.Parallel()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}

	organizationID := "testorgatest"
	stackID := "test"

	natsSubject := utils.GetNatsSubject(organizationID, stackID)
	sc := controllers.NewStargateController(
		nc,
		logging.Testing(),
		opentelemetry.NewNoOpMetricsRegistry(),
		controllers.NewStargateControllerConfig("test", 1*time.Second),
	)
	router := routes.NewRouter(nil, nil, sc)

	type testCase struct {
		name               string
		method             string
		queryParams        url.Values
		url                string
		response           service.StargateClientMessage
		expectedStatusCode int
		expectedHeaders    http.Header
		expectedBody       []byte
	}

	testCases := []*testCase{
		{
			name:   "success",
			method: http.MethodGet,
			queryParams: url.Values{
				"metadata[roles]": []string{"admin"},
			},
			url: "http://test.staging.formance.cloud/" + organizationID + "/" + stackID + "/api/ledger",
			response: service.StargateClientMessage{
				Event: &service.StargateClientMessage_ApiCallResponse{
					ApiCallResponse: &service.StargateClientMessage_APICallResponse{
						StatusCode: 204,
						Body:       []byte{},
						Headers:    map[string]*service.Values{},
					},
				},
			},
			expectedStatusCode: 204,
			expectedHeaders: http.Header{
				"Vary": []string{"Origin"},
			},
			expectedBody: []byte{},
		},
		{
			name:   "wrong organization id and stack id",
			method: http.MethodGet,
			queryParams: url.Values{
				"metadata[roles]": []string{"admin"},
			},
			url: "http://test.staging.formance.cloud/notfoundtest/test/api/ledger",
			response: service.StargateClientMessage{
				Event: &service.StargateClientMessage_ApiCallResponse{
					ApiCallResponse: &service.StargateClientMessage_APICallResponse{
						StatusCode: 204,
						Body:       []byte{},
						Headers:    map[string]*service.Values{},
					},
				},
			},
			expectedStatusCode: 524,
			expectedHeaders: http.Header{
				"Vary":         []string{"Origin"},
				"Content-Type": []string{"application/json"},
			},
			expectedBody: []byte{},
		},
		{
			name:   "wrong organization id and stack id bis",
			method: http.MethodGet,
			queryParams: url.Values{
				"metadata[roles]": []string{"admin"},
			},
			url: "http://test.staging.formance.cloud/api/ledger",
			response: service.StargateClientMessage{
				Event: &service.StargateClientMessage_ApiCallResponse{
					ApiCallResponse: &service.StargateClientMessage_APICallResponse{
						StatusCode: 204,
						Body:       []byte{},
						Headers:    map[string]*service.Values{},
					},
				},
			},
			expectedStatusCode: 400,
			expectedHeaders: http.Header{
				"Vary":         []string{"Origin"},
				"Content-Type": []string{"application/json"},
			},
			expectedBody: []byte("{\"errorCode\":\"VALIDATION\",\"errorMessage\":\"validation error\"}\n"),
		},
		{
			name:   "wrong organization id and stack id with numbers",
			method: http.MethodGet,
			queryParams: url.Values{
				"metadata[roles]": []string{"admin"},
			},
			url: "http://test.staging.formance.cloud/" + organizationID + "/1234/api/ledger",
			response: service.StargateClientMessage{
				Event: &service.StargateClientMessage_ApiCallResponse{
					ApiCallResponse: &service.StargateClientMessage_APICallResponse{
						StatusCode: 204,
						Body:       []byte{},
						Headers:    map[string]*service.Values{},
					},
				},
			},
			expectedStatusCode: 400,
			expectedHeaders: http.Header{
				"Vary":         []string{"Origin"},
				"Content-Type": []string{"application/json"},
			},
			expectedBody: []byte("{\"errorCode\":\"VALIDATION\",\"errorMessage\":\"validation error\"}\n"),
		},
		{
			name:   "failure, wrong url without orga and stack ids",
			method: http.MethodPost,
			queryParams: url.Values{
				"metadata[roles]": []string{"admin"},
			},
			url: "http://test.staging.formance.cloud/",
			response: service.StargateClientMessage{
				Event: &service.StargateClientMessage_ApiCallResponse{
					ApiCallResponse: &service.StargateClientMessage_APICallResponse{
						StatusCode: 204,
						Body:       []byte{},
						Headers:    map[string]*service.Values{},
					},
				},
			},
			expectedStatusCode: 400,
			expectedHeaders: http.Header{
				"Vary":         []string{"Origin"},
				"Content-Type": []string{"application/json"},
			},
			expectedBody: []byte("{\"errorCode\":\"VALIDATION\",\"errorMessage\":\"validation error\"}\n"),
		},
	}

	for _, test := range testCases {
		testCase := test
		sub, err := nc.QueueSubscribe(natsSubject, natsSubject, func(msg *nats.Msg) {
			var request service.StargateServerMessage
			if err := proto.Unmarshal(msg.Data, &request); err != nil {
				t.Fatal(err)
			}

			testCase.response.CorrelationId = request.CorrelationId
			data, err := proto.Marshal(&testCase.response)
			if err != nil {
				t.Fatal(err)
			}

			err = msg.Respond(data)
			if err != nil {
				t.Fatal(err)
			}
		})
		require.NoError(t, err)

		req := httptest.NewRequest(testCase.method, testCase.url, nil)
		rec := httptest.NewRecorder()
		req.URL.RawQuery = testCase.queryParams.Encode()

		router.ServeHTTP(rec, req)

		recBody, err := io.ReadAll(rec.Body)
		require.NoError(t, err)

		require.Equal(t, testCase.expectedStatusCode, rec.Result().StatusCode, "test '%s' failed", testCase.name)
		require.Equal(t, testCase.expectedBody, recBody, "test '%s' failed", testCase.name)
		require.Equal(t, testCase.expectedHeaders, rec.Result().Header, "test '%s' failed", testCase.name)

		err = sub.Unsubscribe()
		require.NoError(t, err)
	}
}
