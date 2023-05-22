package grpc_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/formancehq/stack/components/stargate/internal/api"
	"github.com/formancehq/stack/components/stargate/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestStream(t *testing.T) {
	ctx := context.Background()

	client := NewClient()
	defer client.Close()

	organizationID := "test"
	stackID := "test"

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	responseChan := make(chan *api.StargateClientMessage)
	incomingMessageChan := client.RunStream(t, ctx, organizationID, stackID, responseChan)

	type testCase struct {
		name     string
		ev       *api.StargateServerMessage
		response *api.StargateClientMessage
	}

	testCases := []*testCase{
		{
			name: "success",
			ev: &api.StargateServerMessage{
				Event: &api.StargateServerMessage_ApiCall{
					ApiCall: &api.StargateServerMessage_APICall{
						Method:  http.MethodGet,
						Path:    "/api/ledger",
						Query:   nil,
						Body:    nil,
						Headers: nil,
					},
				},
			},
			response: &api.StargateClientMessage{
				Event: &api.StargateClientMessage_ApiCallResponse{
					ApiCallResponse: &api.StargateClientMessage_APICallResponse{
						StatusCode: 204,
						Body:       nil,
						Headers:    nil,
					},
				},
			},
		},
		{
			name: "success with query, body and headers",
			ev: &api.StargateServerMessage{
				Event: &api.StargateServerMessage_ApiCall{
					ApiCall: &api.StargateServerMessage_APICall{
						Method: http.MethodGet,
						Path:   "/api/ledger",
						Query: map[string]*api.Values{
							"foo": {
								Values: []string{"bar"},
							},
						},
						Body: []byte(`{"foo":"bar"}`),
						Headers: map[string]*api.Values{
							"fake-auth": {
								Values: []string{"fake-token"},
							},
						},
					},
				},
			},
			response: &api.StargateClientMessage{
				Event: &api.StargateClientMessage_ApiCallResponse{
					ApiCallResponse: &api.StargateClientMessage_APICallResponse{
						StatusCode: 200,
						Body:       []byte(`{"bar":"baz"}`),
						Headers: map[string]*api.Values{
							"fake-headers": {
								Values: []string{"fake-value"},
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		data, err := proto.Marshal(testCase.ev)
		require.NoError(t, err)

		goCtx, cancel := context.WithCancel(ctx)
		go func() {
			select {
			case <-goCtx.Done():
				return
			case ev := <-incomingMessageChan:
				testCaseMsg := testCase.ev.Event.(*api.StargateServerMessage_ApiCall)
				serverMsg := ev.Event.(*api.StargateServerMessage_ApiCall)
				assert.Equal(t, testCaseMsg.ApiCall.Method, serverMsg.ApiCall.Method)
				assert.Equal(t, testCaseMsg.ApiCall.Path, serverMsg.ApiCall.Path)
				for k, v := range testCaseMsg.ApiCall.Query {
					assert.Equal(t, v.Values, serverMsg.ApiCall.Query[k].Values)
				}
				for k, v := range testCaseMsg.ApiCall.Headers {
					assert.Equal(t, v.Values, serverMsg.ApiCall.Headers[k].Values)
				}
				assert.Equal(t, testCaseMsg.ApiCall.Body, serverMsg.ApiCall.Body)

				testCase.response.CorrelationId = ev.CorrelationId
				select {
				case <-goCtx.Done():
					return
				case responseChan <- testCase.response:
				}
			}
		}()

		resp, err := nc.Request(utils.GetNatsSubject(organizationID, stackID), data, 30*time.Second)
		cancel()
		require.NoError(t, err)

		var response api.StargateClientMessage
		err = proto.Unmarshal(resp.Data, &response)
		require.NoError(t, err)

		testCaseMsg := testCase.response.Event.(*api.StargateClientMessage_ApiCallResponse)
		responseMsg := response.Event.(*api.StargateClientMessage_ApiCallResponse)
		assert.Equal(t, testCase.response.CorrelationId, response.CorrelationId)
		assert.Equal(t, testCaseMsg.ApiCallResponse.StatusCode, responseMsg.ApiCallResponse.StatusCode)
		for k, v := range testCaseMsg.ApiCallResponse.Headers {
			assert.Equal(t, v.Values, responseMsg.ApiCallResponse.Headers[k].Values)
		}
		assert.Equal(t, testCaseMsg.ApiCallResponse.Body, responseMsg.ApiCallResponse.Body)
	}
}
